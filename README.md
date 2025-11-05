# rest_service

> RESTful API service implemented in Go using the Gin framework

---

## Overview

This repository contains a **REST API service** implemented in **Go** with the **Gin** web framework.  

The project includes:
- User authentication endpoints (`/auth/sign-up`, `/auth/sign-in`)
- Book management endpoints (`/books`)
- PostgreSQL integration
- Docker setup for running the database

---

## Repository Structure

```bash
├── server.go
├── cmd/            # main.go (entry point)
├── models/         # models.go 
├── configs/        # config.yml
├── pkg/            # Packages: handler, repository, service.
├── go.mod
├── go.sum
├── .env
└── .gitignore
````

---

## Requirements

* Go **1.20+**
* PostgreSQL database (local or Docker)
* Docker *(recommended for local development)*

---

## Docker Setup for PostgreSQL

Run the PostgreSQL database using Docker:

```bash
docker run --name pg-container \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=sword \
  -e POSTGRES_DB=go_crud_db \
  -p 5433:5432 \
  -d postgres
```

This command will:

* Create a PostgreSQL container named **pg-container**
* Set credentials → user: `postgres`, password: `sword`
* Create a database named **go_crud_db**
* Expose port **5433** on your host

---

## Running the Service

From the root of the repository:

```bash
go run ./cmd/main.go
```

The server will start and listen for HTTP requests on **port 8080** (by default).

---

## API Endpoints

### Authentication

| Method | Endpoint        | Description                   |
| ------ | --------------- | ----------------------------- |
| `POST` | `/auth/sign-up` | Register a new user           |
| `POST` | `/auth/sign-in` | Authenticate an existing user |

### Books

| Method   | Endpoint     | Description           |
| -------- | ------------ | --------------------- |
| `POST`   | `/books/`    | Create a new book     |
| `GET`    | `/books/`    | Retrieve all books    |
| `GET`    | `/books/:id` | Retrieve a book by ID |
| `PUT`    | `/books/:id` | Update a book by ID   |
| `DELETE` | `/books/:id` | Delete a book by ID   |

---

## Configuration

The service uses **environment variables** for configuration:

| Variable      | Description       | Default      |
| ------------- | ----------------- | ------------ |
| `DB_HOST`     | Database host     | `localhost`  |
| `DB_PORT`     | Database port     | `5433`       |
| `DB_USER`     | Database username | `postgres`   |
| `DB_PASSWORD` | Database password | `sword`      |
| `DB_NAME`     | Database name     | `go_crud_db` |
| `SERVER_PORT` | HTTP server port  | `8080`       |

