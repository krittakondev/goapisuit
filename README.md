[![Project Logo](logo.png)](https://github.com/krittakondev/goapisuit)

# goapisuit

[![Go Reference](https://pkg.go.dev/badge/github.com/krittakondev/goapisuit.svg)](https://pkg.go.dev/github.com/krittakondev/goapisuit)

**goapisuit** is a lightweight and easy-to-use Golang framework that simplifies API creation and database management using GoFiber. It aims to accelerate the development of APIs by providing built-in commands for generating models, routers, and handling database migrations.

## Features

- **Fast API Creation**: Automatically generate routes and models.
- **Database Management**: Seamlessly handle migrations with built-in commands.
- **GoFiber Integration**: Utilize the powerful and efficient GoFiber web framework.

## Getting Started

### Prerequisites

To use goapisuit, you need the following installed:

- [Golang](https://golang.org/dl/)
- A relational database (e.g., MySQL etc.)
- Git for cloning the repository

### Installation

1. Install CLI `heykrit` for the `goapisuit`:

   ```bash
   go install -v github.com/krittakondev/goapisuit/v2/cmd/heykrit@latest
   ```
2. Make your project:

   ```bash
   mkdir nameproject
   cd nameproject
   go mod init you/projectpath
   heykrit init
   ```
3. Config your project in `.env`

4. Run your project:

   ```bash
   go run cmd/server.go
   ```

### Directory Structure

Here's a brief overview of the generated directory structure:

```
nameproject/
├── cmd/
│   └── server.go        # Entry point of your project
├── internal/
│   ├── models/          # Contains your data models (e.g., user.go)
│   └── routes/             # Contains API route handlers
├── public/              # Directory for static files (e.g., CSS, JavaScript, images)
├── go.mod
└── go.sum
└── .env                 # config project

```

## Usage

Once the setup is complete, you can build upon this project by adding new models and routes using the provided commands. For example, to create a new module for handling products:


1. Generate a new `product` Route and Model:

   ```bash
   heykrit make product
   ```
   Generate 2 files `internal/routes/Product.go` and `internal/models/Product.go`

2. Modify the generated `internal/models/Product.go` model as needed, and then apply the migration (gorm model):

   ```bash
   heykrit db:migrate product
   ```
3. Run server:

   ```bash
   go run cmd/server.go
   ```

## Contributing

If you'd like to contribute to `goapisuit`, feel free to open a pull request or issue on the [GitHub repository](https://github.com/krittakondev/goapisuit).

## License

This project is licensed under the MIT License.

## Credits

[![GoFiber](https://img.shields.io/badge/GoFiber-API_Framework-blue)](https://gofiber.io/)A web framework that brings lightning-fast performance to your Golang applications.
[![GORM](https://img.shields.io/badge/GORM-ORM_Library-lightgrey)](https://gorm.io/) 
A powerful ORM library for Golang, simplifying database handling and migrations.
