package main

import (
	"log"
	"os"
	"reflect"
	"strings"

	// "strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	routes "github.com/krittakondev/goapisuit/internal/api"
	"github.com/krittakondev/goapisuit/internal/database"
	// routesAll "github.com/krittakondev/goapisuit/internal/api/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("'cp .env.example .env' and edit .env")
	}
	app := fiber.New(fiber.Config{
		ServerHeader: "goapisuit",
		AppName:      os.Getenv("APP_NAME"),
	})
	db, err := database.MysqlConnect()
	if err != nil {
		log.Fatal(err)
	}

	// auto migrate if you want
	db.AutoMigrate(
		// &models.Users{},
	)

	HOST := os.Getenv("APP_HOST")
	PORT := os.Getenv("APP_PORT")
	if PORT == "" {
		PORT = "3000"
	}

	app.Static("/", "./public")

	api_prefix := "/api"
	if prefix := os.Getenv("API_PREFIX"); prefix != "" {
		prefix = prefix
	}
	api := app.Group(api_prefix)

	
	mkroute := &routes.Route{
		DB: db,
		PageLimit: 20,
	}

	t := reflect.TypeOf(mkroute)
	reflect_val := reflect.ValueOf(mkroute)

	for i := 0; i < t.NumMethod(); i++ {
		namemethod := t.Method(i).Name
		method := reflect.Value.MethodByName(reflect_val, namemethod)
		handler := func(c *fiber.Ctx) error {
			// Prepare the arguments for the function
			args := []reflect.Value{reflect.ValueOf(c)}

			// Call the method
			results := method.Call(args)

			// Check for an error and return it
			if len(results) > 0 && results[0].Interface() != nil {
				return results[0].Interface().(error)
			}
			return nil
		}
		path_split := strings.Split(strings.ToLower(namemethod), "_")
		apipath :=  path_split[0]
		route_method := ""
		if len(path_split) > 1{
			route_method = path_split[1]
		}
		if apipath == "index"{
			apipath = ""
		}
		
		if route_method == "get"{
			api.Get("/"+apipath, handler)
			log.Print(route_method+" "+api_prefix+"/"+apipath)
			api.Get("/"+apipath+"/:id", handler)
			if apipath != ""{
				log.Print(route_method+" "+api_prefix+"/"+apipath+"/:id")
			}
		}
		if route_method == "delete"{
			api.Delete("/"+apipath+"/:id", handler)
			log.Print(route_method+" "+api_prefix+"/"+apipath+"/:id")
		}
		if route_method == "put"{
			api.Put("/"+apipath+"/:id", handler)
			log.Print(route_method+" "+api_prefix+"/"+apipath+"/:id")
		}
		if route_method == "post"{
			api.Post("/"+apipath, handler)
			log.Print(route_method+" "+api_prefix+"/"+apipath)
		}
		// api.All("/"+apipath, handler)
	}

	app.Listen(HOST + ":" + PORT)
}
