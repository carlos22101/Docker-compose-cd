package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var db *sql.DB

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func init() {
	// Intentar cargar .env desde varias ubicaciones comunes
	paths := []string{".env", "../.env", "/app/.env"}
	loaded := false
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			if err := godotenv.Load(p); err == nil {
				log.Println("Loaded .env from", p)
				loaded = true
				break
			}
		}
	}
	if !loaded {
		// intenta cargar por defecto (godotenv buscaría .env en cwd)
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		} else {
			log.Println("Loaded .env from default location")
		}
	}
}

func main() {
	// Leer variables de entorno (si no existen, getenv permite default)
	host := getenv("DB_HOST", "db_carlos")
	port := getenv("DB_PORT", "3306")
	user := getenv("DB_USER", "carlos")
	pass := getenv("DB_PASSWORD", "1234")
	name := getenv("DB_NAME", "carlos_DB")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", user, pass, host, port, name)

	// Reintentos con backoff para esperar a que MySQL esté listo
	maxAttempts := 30
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("attempt %d/%d: error ping db: %v", attempt, maxAttempts, err)
		if db != nil {
			_ = db.Close()
		}
		sleepSec := attempt
		if sleepSec > 10 {
			sleepSec = 10
		}
		time.Sleep(time.Duration(sleepSec) * time.Second)
	}
	if err != nil {
		log.Fatalf("could not connect to db after %d attempts: %v", maxAttempts, err)
	}
	log.Println("Conectado a DB:", name)

	r := mux.NewRouter()
	// registrar middleware normal (por compatibilidad)
	r.Use(corsMiddleware)

	// endpoints CRUD
	r.HandleFunc("/users", listUsers).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", deleteUser).Methods("DELETE")

	// endpoint con tu apellido
	r.HandleFunc("/solis", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"fullname": "Carlos Solis"})
	}).Methods("GET")

	// Wrapper global que asegura CORS en todas las respuestas y responde OPTIONS
	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Cabeceras CORS globales
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")
		// Responder preflight
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Dejar que el router procese la petición
		r.ServeHTTP(w, req)
	})

	addr := ":8000"
	log.Println("API backend corriendo en", addr)
	log.Fatal(http.ListenAndServe(addr, corsHandler))
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func corsMiddleware(next http.Handler) http.Handler {
	// middleware adicional (ya tenemos wrapper global pero mantenemos esto por compatibilidad)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cabeceras para CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT id, first_name, last_name, email FROM users ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idS := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idS)
	var u User
	err := db.QueryRow("SELECT id, first_name, last_name, email FROM users WHERE id = ?", id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(u)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	res, err := db.Exec("INSERT INTO users (first_name, last_name, email) VALUES (?, ?, ?)", u.FirstName, u.LastName, u.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	u.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idS := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idS)
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	_, err := db.Exec("UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?", u.FirstName, u.LastName, u.Email, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.ID = id
	json.NewEncoder(w).Encode(u)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	idS := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idS)
	_, err := db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
