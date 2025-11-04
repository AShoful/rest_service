# rest_service

> RESTful API service implemented in Go using Gin framework

## Overview
This repository contains a REST API service implemented in Go.  
It includes:
- User authentication endpoints (`/auth/sign-up`, `/auth/sign-in`)  
- Book management endpoints (`/books`)  
- PostgreSQL integration  
- Docker setup for database  

## Repository Structure

├── server/  REST API server implementation \
├── pkg/  Packages: handler, models, repository, etc. \
├── go.mod \
├── go.sum \
└── .gitignore


## Requirements
- Go 1.20 or higher  
- PostgreSQL database (can run via Docker)  
- Docker (optional, recommended for local development)

## Docker Setup for PostgreSQL
Run the PostgreSQL database using Docker:

docker run --name pg-container \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=sword \
  -e POSTGRES_DB=go_crud_db \
  -p 5433:5432 \
  -d postgres

This will: \
Create a PostgreSQL container named pg-container \
Use the user postgres with password sword \
Create a database go_crud_db \
Expose port 5433 on your host \

Running the Service \
From the root of the repository: \

cd server \
go run main.go \
The server will start and listen for HTTP requests (default port: 8080). \

API Endpoints \

Authentication \
POST	/auth/sign-up	Register a new user \
POST	/auth/sign-in	Authenticate a user \

Books \

POST	/books/	Create a new book \
GET	/books/	Retrieve all books \
GET	/books/:id	Retrieve a book by ID \
PUT	/books/:id	Update a book by ID \
DELETE	/books/:id	Delete a book by ID \

Configuration \
The service uses environment variables for configuration: \
DB_HOST — database host (default: localhost) \
DB_PORT — database port (default: 5433) \
DB_USER — database username (default: postgres) \
DB_PASSWORD — database password (default: sword) \

DB_NAME — database name (default: go_crud_db)

SERVER_PORT — HTTP server port (default: 8080)
