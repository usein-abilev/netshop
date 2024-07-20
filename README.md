# A tiny online shop backend
This is a simple online shop backend written in Go. It is a RESTful API that allows you to manage products, customers, orders and order items.

My personal goal with this project is to learn about Go and how to build a RESTful API with it.

## Installation
1. Clone the repository
2. Run `go mod download` to download all dependencies
3. Create a `.env` file in the root directory and add the following environment variables:
```properties
JWT_SECRET=secret
JWT_EXPIRATION=expiration_in_ms ; default is 3600000

ARGON2_SALT=argon2_salt
ARGON2_THREADS=argon2_threads ; default is 4

export DATABASE_URL=postgres://username:password@localhost:1900/dbname
export SERVER_URL=0.0.0.0
export SERVER_PORT=8080
```
4. Install `docker-compose` and run `docker-compose up -d` to install and start the PostgreSQL database
5. Run `go run .` to start the server

!!! Note: You can also run `make setup` to automatically install the dependencies, create the `.env` file and start the database.
