package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	routes "github.com/krittakondev/goapisuit/internal/api"
)

func main() {
	if err := godotenv.Load(); err != nil{
		log.Fatal("'cp .env.example .env' and edit .env")
	}
	app := fiber.New(fiber.Config{
	    ServerHeader:  "goapisuit",
	    AppName: os.Getenv("APP_NAME"),
	})
	HOST := os.Getenv("APP_HOST")
	PORT := os.Getenv("APP_PORT")
	if PORT == ""{
		PORT = "3000"
	}

	app.Static("/", "./public")

	api_prefix := "/api"
	if prefix:=os.Getenv("API_PREFIX"); prefix!=""{
		prefix = prefix
	}
	api := app.Group(api_prefix)

	api.Get("/", routes.MainRoute)

	app.Listen(HOST+":"+PORT)
}
