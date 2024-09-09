[![Project Logo](logo.png)](https://github.com/krittakondev/goapisuit)

# goapisuit

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

1. Clone the `goapisuit` repository:

   ```bash
   git clone https://github.com/krittakondev/goapisuit.git nameproject
   ```

2. Change into your project directory:

   ```bash
   cd nameproject
   ```

3. Run your project and generate your first user module:

   ```bash
   go run cmd/heykrit/main.go make user
   ```

   This will automatically create the `user` router and the `user` model under the `internal/models` directory.

### Database Migration

After modifying the model file (e.g., `internal/models/user.go`), you can apply the migration using the following command:

1. Run the migration for the `user` model:

   ```bash
   go run cmd/heykrit/main.go db:migrate user
   ```

2. Confirm the migration by typing `y` when prompted.

### Directory Structure

Here's a brief overview of the generated directory structure:

```
nameproject/
├── cmd/
│   └── server/
│       └── main.go      # Entry point of your project
├── internal/
│   ├── models/          # Contains your data models (e.g., user.go)
│   └── routers/         # Contains API route handlers
├── public/              # Directory for static files (e.g., CSS, JavaScript, images)
├── go.mod
└── go.sum
└── .env                 # config project

```

## Usage

Once the setup is complete, you can build upon this project by adding new models and routes using the provided commands. For example, to create a new module for handling products:

1. Generate a new `product` module:

   ```bash
   go run cmd/heykrit/main.go make product
   ```

2. Modify the generated `product` model as needed, and then apply the migration:

   ```bash
   go run cmd/heykrit/main.go db:migrate product
   ```
3. Run server:

   ```bash
   go run cmd/serever/main.go
   ```

## Contributing

If you'd like to contribute to `goapisuit`, feel free to open a pull request or issue on the [GitHub repository](https://github.com/krittakondev/goapisuit).

## License

This project is licensed under the MIT License.
