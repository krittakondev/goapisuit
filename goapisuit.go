package goapisuit

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	// "strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/krittakondev/goapisuit/database"
	"github.com/krittakondev/goapisuit/middlewares"
	"gorm.io/gorm"
	// routesAll "github.com/krittakondev/goapisuit/internal/api/routes"
)


type Suit struct{
	ProjectName string
	DB *gorm.DB
	LimitPage int
	RequireJwtAuth func(*fiber.Ctx) error
	Fiber *fiber.App
	Groups *[]interface{}
}


func LoadEnv(){
	if err := godotenv.Load(); err != nil {
		log.Fatal("'cp .env.example .env' and edit .env")
	}
}

func LoadTmpModel() (arr []string, err error){
	read, err := os.ReadFile(".tmpmodels")
	if err != nil {
		return
	}
	arr = strings.Split(string(read), "\n")
	return
}

func New(project_name string) (*Suit, error){
	LoadEnv()
	conn, err := database.MysqlConnect()
	if err != nil {
		return &Suit{}, err
	}
	app := fiber.New(fiber.Config{
		ServerHeader: "goapisuit",
		AppName:      os.Getenv("APP_NAME"),
	})

	limit := 20
	if os.Getenv("API_LIMIT_PAGE") != ""{
		limit, err = strconv.Atoi(os.Getenv("API_LIMIT_PAGE"))
		if err != nil {
			panic(err)
		}
	}

	return &Suit{
		RequireJwtAuth: middlewares.RequireJwtAuth,
		DB: conn,
		LimitPage: limit,
		ProjectName: project_name,
		Fiber: app,
	}, nil
}

func (s *Suit) SetupRoutes(r interface{}){
	api_prefix := "/api"
	if prefix := os.Getenv("API_PREFIX"); prefix != "" {
		api_prefix = prefix
	}
	api := s.Fiber.Group(api_prefix)


	t := reflect.TypeOf(r)
	reflect_val := reflect.ValueOf(r)

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
			api.Get("/", handler)
			apipath = ""
		}
		
		if route_method == "get"{
			// api.Get("/"+apipath, handler)
			log.Print(route_method+" "+api_prefix+"/"+apipath)
			api.Get("/"+apipath+"/:id?", handler)
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

}

func (s *Suit) Run(r interface{}){

	HOST := os.Getenv("APP_HOST")
	PORT := os.Getenv("APP_PORT")
	if PORT == "" {
		PORT = "3000"
	}

	s.Fiber.Static("/", "./public")

	s.SetupRoutes(r)

	if err := s.Fiber.Listen(HOST + ":" + PORT); err != nil{
		log.Fatal(err)
	}
}
