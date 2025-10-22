package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"

    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
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
    if err := godotenv.Load("../.env"); err != nil {
        log.Println("No .env file found")
    }
}

func main() {
   host := os.Getenv("DB_HOST")
   port := os.Getenv("DB_PORT")
   user := os.Getenv("DB_USER")
   pass := os.Getenv("DB_PASSWORD")
   name := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", user, pass, host, port, name)
    var err error
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("error open db: %v", err)
    }
    if err = db.Ping(); err != nil {
        log.Fatalf("error ping db: %v", err)
    }
    log.Println("Conectado a DB:", name)

    r := mux.NewRouter()
    r.Use(corsMiddleware)

    r.HandleFunc("/users", listUsers).Methods("GET")
    r.HandleFunc("/users/{id:[0-9]+}", getUser).Methods("GET")
    r.HandleFunc("/users", createUser).Methods("POST")
    r.HandleFunc("/users/{id:[0-9]+}", updateUser).Methods("PUT")
    r.HandleFunc("/users/{id:[0-9]+}", deleteUser).Methods("DELETE")

    r.HandleFunc("/solis", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"fullname": "Carlos Solis"})
    }).Methods("GET")

    addr := ":8000"
    log.Println("API backend corriendo en", addr)
    log.Fatal(http.ListenAndServe(addr, r))
}

func getenv(k, d string) string {
    v := os.Getenv(k)
    if v == "" {
        return d
    }
    return v
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
    rows, err := db.Query("SELECT id, first_name, last_name, email FROM users ORDER BY id DESC")
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    defer rows.Close()
    users := []User{}
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email); err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        users = append(users, u)
    }
    json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
    idS := mux.Vars(r)["id"]
    id, _ := strconv.Atoi(idS)
    var u User
    err := db.QueryRow("SELECT id, first_name, last_name, email FROM users WHERE id = ?", id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Not found", 404)
            return
        }
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(u)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var u User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        http.Error(w, "invalid body", 400)
        return
    }
    res, err := db.Exec("INSERT INTO users (first_name, last_name, email) VALUES (?, ?, ?)", u.FirstName, u.LastName, u.Email)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    id, _ := res.LastInsertId()
    u.ID = int(id)
    w.WriteHeader(201)
    json.NewEncoder(w).Encode(u)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
    idS := mux.Vars(r)["id"]
    id, _ := strconv.Atoi(idS)
    var u User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        http.Error(w, "invalid body", 400)
        return
    }
    _, err := db.Exec("UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?", u.FirstName, u.LastName, u.Email, id)
    if err != nil {
        http.Error(w, err.Error(), 500)
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
        http.Error(w, err.Error(), 500)
        return
    }
    w.WriteHeader(204)
}