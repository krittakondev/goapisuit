package goapisuit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	// "strings"

	"github.com/gofiber/fiber/v2"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/krittakondev/goapisuit/database"
	"github.com/krittakondev/goapisuit/middlewares"
	"gorm.io/gorm"
	// routesAll "github.com/krittakondev/goapisuit/internal/api/routes"
)

type RouteGroup struct{
	Parent string
	Children *[]interface{}
}


type Suit struct{
	ProjectName string
	DB *gorm.DB
	LimitPage int
	RequireJwtAuth func(*fiber.Ctx) error
	Fiber *fiber.App
	Groups *[]RouteGroup
	Routes *interface{}
	Config Config
}

type Config struct {
	AppName string `env:"APP_NAME"`
	AppHost string `env:"APP_HOST"`
	AppPort string `env:"APP_PORT"`

	ApiPrefix string `env:"API_PREFIX"`
	ApiLimitPage int `env:"API_LIMIT_PAGE"`
	DbConnection string `env:"DB_CONNECTION"`
	
}

func LoadEnv() Config{
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error parsing env variables: %s", err)
	}
	return cfg
}

func LoadTmpModel() (arr []string, err error){
	read, err := os.ReadFile(".tmpmodels")
	if err != nil {
		return
	}
	arr = strings.Split(string(read), "\n")
	return
}

func New(project_name string) (suit *Suit, err error){
	cfg := LoadEnv()
	suit = &Suit{
		RequireJwtAuth: middlewares.RequireJwtAuth,
	}
	suit.ProjectName = project_name
	suit.Config = cfg
	if strings.ToLower(suit.Config.DbConnection) == "mysql"{
		conn, err := database.MysqlConnect()
		if err != nil {
			return suit, err
		}
		suit.DB = conn
	}
	app := fiber.New(fiber.Config{
		ServerHeader: "goapisuit",
		AppName:      suit.Config.AppName,
	})
	suit.Fiber = app

	limit := 20
	if cfg.ApiLimitPage != 0 {
		limit = cfg.ApiLimitPage
	}
	suit.LimitPage = limit

	return suit, nil
}

func (s *Suit) groupScan()error{
	err := filepath.Walk("internal/routes", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ตรวจสอบว่าเป็นไฟล์ Go และไม่ใช่ main package
		if strings.HasSuffix(info.Name(), ".go") && info.Name() != "main.go" {
			// นำ path ของโฟลเดอร์มาสร้าง API group
			groupPath := strings.TrimPrefix(filepath.Dir(path), "internal/routes")
			groupPath = strings.ReplaceAll(groupPath, "\\", "/") // สำหรับ Windows
			groupPath = "/" + strings.Trim(groupPath, "/")

			fmt.Println(groupPath)
		}
		return nil
	})
	return err
}

func (s *Suit) SetupGroups(api_prefix string, r interface{}, middleware ...fiber.Handler){
	api := s.Fiber.Group(api_prefix)
	
	t := reflect.TypeOf(r)
	reflect_val := reflect.ValueOf(r)

	if reflect_val.Kind() == reflect.Ptr {
		val := reflect_val.Elem()

		field := val.FieldByName("Suit")
		newVal := reflect.ValueOf(s)
		field.Set(newVal)
	} 
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
		for _, mid := range middleware{
			api.Use(mid)
		}
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

func (s *Suit) SetupRoutes(r interface{}){
	api_prefix := "/api"
	if prefix := os.Getenv("API_PREFIX"); prefix != "" {
		api_prefix = prefix
	}
	
	s.SetupGroups(api_prefix, r)
}

func (s *Suit) Run(){

	HOST := os.Getenv("APP_HOST")
	PORT := os.Getenv("APP_PORT")
	if PORT == "" {
		PORT = "3000"
	}

	s.Fiber.Static("/", "./public")


	if err := s.Fiber.Listen(HOST + ":" + PORT); err != nil{
		log.Fatal(err)
	}
}
