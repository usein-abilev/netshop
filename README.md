# A tiny online shop backend
This is a simple online shop backend written in Go. It is a RESTful API that allows you to manage products, customers, orders and order items.

My personal goal with this project is to learn about Go and how to build a RESTful API with it.

- [Requirements](#requirements)
- [Installation](#installation)
- [Development Environment](#development-environment)
- [Project Structure](#project-structure)
- [API Endpoints](#api-endpoints)

# Requirements
- Go 1.21+
- PostgreSQL 16+
- [libvips](https://www.libvips.org/install) 8.15.2

This project uses [libvips](https://www.libvips.org) to compress and convert images to WebP format. Make sure to install it on your system.

# Installation
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

# Development Environment
1. [go-task](https://github.com/go-task/task) used to simplify the development process by providing a **hot-reload** feature.
2. [dbmate](https://github.com/amacneil/dbmate) is used to manage the database migrations.

# Project Structure
The project is organized into the following packages for clarity and maintainability:

#### `/api`
This package contains all the API handlers and routes.
### `/config`
This package contains the configuration for the application. It reads the environment variables and convert into an internal configuration struct.
### `/db`
All database-related code is contained in this package, including the initialization of the database connection and the database models.
### `/migrations`
This package contains all the database migrations. The [dbmate](https://github.com/amacneil/dbmate) tool is used to organize and execute these migrations.
### `/tools` 
This package contains all the tools and utilities used in the project. For example, the `image` package contains the image compression and conversion logic.

# API Endpoints
API versioning is used in this project, with the current version being v1. To see the full list of available endpoints, run the server and navigate to /api/v1/docs in your browser. Below is a summary of the main endpoints:
### Auth
- `POST /api/v1/auth/login` - Authenticate a customer or employee and receive a JWT token.
- `POST /api/v1/auth/customer/signup` - Register a new customer account
- `POST /api/v1/auth/employee/signup` - Register a new employee account
- `POST /api/v1/auth/verify` - Verify a JWT token

### Products
- `GET /api/v1/products` - Get all products (public access)
- `POST /api/v1/products` - Create a new product (admin users only)
- `PUT /api/v1/products/{id}` - Update a product (admin users only)
- `DELETE /api/v1/products/{id}` - Delete a product (admin users only)

### Orders
- `GET /api/v1/orders` - Get all orders (admin users or customers only)
- `POST /api/v1/orders` - Create a new order (customers only)
- `PUT /api/v1/orders/{id}` - Update an order (admin users or customers only)

### Files
- `POST /api/v1/file/upload` - Upload a new file. Files stored as a compressed WEBP file. Supported formats are PNG, JPEG, JPG, and WEBP (authenticated users only).
- `GET /static/files/{filename}` - Receive a file by its filename