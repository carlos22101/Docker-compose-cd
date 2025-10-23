# Microservicios - Proyecto (Carlos Solis)

**Autor:** Carlos Solis  
**Repositorio:** reemplaza con la URL de tu repo: `https://github.com/carlos22101/Docker-compose-cd.git`

---

## Resumen
Proyecto de microservicios con Docker Compose que contiene:
- **Frontend**: `frontend_carlos` (Node/Express, puerto 80) — sirve una SPA HTML/JS/CSS con CRUD de usuarios.
- **Backend**: `backend_solis` (Go, puerto 8000) — API REST con CRUD y endpoint `/solis` que devuelve `{"fullname":"Carlos Solis"}`.
- **Base de datos**: `db_carlos` (MySQL 8.0, puerto 3306) — volumen explícito para persistencia.

Red interna Docker: `carlos_net`.  
Base de datos: **carlos_DB** (usuario `carlos`, contraseña ``).

---

## Estructura
microservices-project/
├─ docker-compose.yml
├─ README.md
├─ db/
│ └─ init.sql
├─ frontend/
│ ├─ Dockerfile
│ ├─ package.json
│ ├─ server.js
│ └─ public/
│ ├─ index.html
│ ├─ app.js
│ └─ style.css
└─ backend/
├─ Dockerfile
├─ go.mod
└─ main.go


---

## Requisitos implementados (checklist)
- [x] Tres servicios definidos en `docker-compose.yml`.
- [x] Redes internas (`carlos_net`) y comunicación por nombre de servicio (`db_carlos`, `backend_solis`, `frontend_carlos`).
- [x] Nombres de contenedores contienen `carlos` o `solis`.
- [x] Volumen explícito `db_data` para persistencia de MySQL.
- [x] `depends_on` configurado para orden de arranque.
- [x] Cada servicio (frontend, backend) tiene su propio `Dockerfile`.
- [x] No se usa `nginx`/imágenes preconstruidas para servir estático; frontend usa imagen construida con Node.
- [x] Backend en Go con variables de entorno para DB y endpoint `/solis`.
- [x] CORS manejado en backend para permitir llamadas desde frontend.
- [x] Procedimiento de pruebas de persistencia documentado.

---

## Cómo levantar el entorno (EC2 / local)

1. Con Docker y Docker Compose instalados (ver comandos de instalación en la repo original).
2. Clona el repo (o asegúrate de estar en la carpeta con `docker-compose.yml`):
```bash
cd ~/Docker-compose-cd
# si aun no esta git init/push, realiza git init/push antes de clonar en EC2
# (reconstruir imágenes si hiciste cambios)
docker compose build --no-cache

# levantar en background
docker compose up -d
Variables de entorno
Verificar servicios:

Frontend: http://<IP_ELASTICA>

Backend: http://<IP_ELASTICA>:8000/

Endpoint de prueba: http://<IP_ELASTICA>:8000/solis → {"fullname":"Carlos "}

El proyecto usa un archivo .env en la raíz con variables
