# Animal API

A RESTful API in Go (Gin, GORM, MySQL) for managing animal data, following Clean Architecture and Dockerized.

## Key Features

* **POST /v1/animal**: Create (client-provided id, fails if id exists).
* **PUT /v1/animal**: Upsert (client-provided id in payload; updates if id exists, creates if not).
* **DELETE /v1/animal/:id**: Delete by path id (404 if not found).
* **GET /v1/animal**: List all (404 if empty).
* **GET /v1/animal/:id**: Get by path id (404 if not found).

## Prerequisites

* Git
* Go (1.21+)
* Docker & Docker Compose

## Configuration

Configure via environment variables or a .env file (copy from .env.example). Key variables:
* GIN_MODE: debug (default), release, test.
* APP_PORT: Application port (default: 8080).
* Database: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, or full DB_DSN.

## Running the Application

*1. Clone Repository:*
```bash
git clone git@github.com:Chyoad/go-animal-api.git

cd go-animal-api
```
*2. Using Docker Compose (Recommended):*

Creates .env from .env.example in project root if you need to customize docker-compose.yml substitutions (e.g. APP_PORT_HOST).

To run storage (MySQL) and application:
```bash
docker-compose up --build -d
```
API Address: http://localhost:<APP_PORT_HOST>/v1/animal (default: http://localhost:8080/v1/animal)

To view the log:
```bash
docker-compose logs -f app       # Log Go App
docker-compose logs -f mysql_db  # Log database MySQL
```

Stops all services: 
```bash
docker-compose down
```

To delete a data volume (all database data will be lost!):
```bash
docker-compose down -v
```

*3. Running Locally (for Development):*

- Run Storage System: Start your MySQL instance (e.g., local install or docker run ... mysql:8.0). Using docker for MySQL only :
```bash
docker run --name local-animal-mysql -e MYSQL_ROOT_PASSWORD=rootpass \ -e MYSQL_DATABASE=animal_db_local -e MYSQL_USER=user_local -e MYSQL_PASSWORD=pass_local \ -p 3306:3306 -d mysql:8.0
```
- Configure Application: Create .env from .env.example and set your local DB connection details and APP_PORT.

Run Application:
```bash 
go mod tidy 
go run ./cmd/api/main.go
```

API Address: http://localhost:<APP_PORT_IN_DOTENV>/v1/animal (default: http://localhost:8080/v1/animal)