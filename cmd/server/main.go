package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
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

	prefix := "/"
	if prefix_api:=os.Getenv("PREFIX_API"); prefix_api!=""{
		prefix = prefix_api
	}
	api := app.Group(prefix)

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello goAPIsuit!")
	})

	app.Listen(HOST+":"+PORT)
}
