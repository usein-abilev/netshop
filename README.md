# A tiny online shop backend
This is a simple online shop backend written in Go. It is a RESTful API that allows you to manage products, customers, orders and order items.

My personal goal with this project is to learn about Go and how to build a RESTful API with it.

## Requirements
- Go 1.21+
- PostgreSQL 16+
- [libvips](https://www.libvips.org/install) 8.15.2

This project uses [libvips](https://www.libvips.org) to compress and convert images to WebP format. Make sure to install it on your system.

## Installation
1. Clone the repository
2. Run `go mod download` to download all dependencies
3. Create a `.env` file in the root directory and add the following environment variables:
```properties
JWT_SECRET=secret
JWT_EXPIRATION=duration ; value for `time.ParseDuration`, default is 24h

export DATABASE_URL=postgres://username:password@localhost:1900/dbname
export SERVER_HOST=0.0.0.0
export SERVER_PORT=8080
```
4. Install `docker-compose` and run `docker-compose up -d` to install and start the PostgreSQL database
5. Run `go run .` to start the server

> Note: If you have `make` util installed, you can also run `make setup` to automatically install the dependencies, create the `.env` file and start the database.

## Development Environment
1. [go-task](https://github.com/go-task/task) used to simplify the development process by providing a **hot-reload** feature.
2. [dbmate](https://github.com/amacneil/dbmate) is used to manage the database migrations.