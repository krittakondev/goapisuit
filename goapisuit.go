package goapisuit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	// "strings"

	"github.com/caarlos0/env/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/krittakondev/goapisuit/database"
	"github.com/krittakondev/goapisuit/pkg/utils"
	"github.com/krittakondev/goapisuit/middlewares"
	"gorm.io/gorm"
	// routesAll "github.com/krittakondev/goapisuit/internal/api/routes"
)

const Version = "v1.0.2"

type Suit struct {
	ProjectName    string
	DB             *gorm.DB
	LimitPage      int
	RequireJwtAuth func(*fiber.Ctx) error
	Fiber          *fiber.App
	Routes         *interface{}
	Config         Config
}

type Config struct {
	AppName string `env:"APP_NAME"`
	AppHost string `env:"APP_HOST"`
	AppPort string `env:"APP_PORT"`

	ApiPrefix    string `env:"API_PREFIX"`
	ApiLimitPage int    `env:"API_LIMIT_PAGE"`
	DbConnection string `env:"DB_CONNECTION"`
}

func LoadEnv() Config {
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

func LoadTmpModel() (arr []string, err error) {
	read, err := os.ReadFile(".tmpmodels")
	if err != nil {
		return
	}
	arr = strings.Split(string(read), "\n")
	return
}

func New(project_name string, fiberConfig ...fiber.Config) (suit *Suit, err error) {
	cfg := LoadEnv()
	suit = &Suit{
		RequireJwtAuth: middlewares.RequireJwtAuth,
	}
	suit.ProjectName = project_name
	suit.Config = cfg
	if strings.ToLower(suit.Config.DbConnection) == "mysql" {
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
	if len(fiberConfig) > 0{
		app = fiber.New(fiberConfig...)
	}
	suit.Fiber = app

	limit := 20
	if cfg.ApiLimitPage != 0 {
		limit = cfg.ApiLimitPage
	}
	suit.LimitPage = limit

	return suit, nil
}

func (s *Suit) GroupScan() (groups []string, err error) {
	err = filepath.Walk("internal/routes", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), "init_suit.go") {
			groupPath := strings.TrimPrefix(filepath.Dir(path), "internal/routes")
			groupPath = strings.ReplaceAll(groupPath, "\\", "/")
			if groupPath == "" {
				groupPath = "/"
			} else {
				groupPath = "/" + strings.Trim(groupPath, "/")
			}
			if groupPath != "/" {
				fmt.Println(groupPath)
				groups = append(groups, groupPath)
			}
		}
		return nil
	})
	return groups, err
}
func handlerReflect(reflect_val reflect.Value, namemethod string) func(c *fiber.Ctx) error {

	method := reflect.Value.MethodByName(reflect_val, namemethod)
	return func(c *fiber.Ctx) error {
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
}

func setupDynamicRoutes(api fiber.Router, r interface{}) (routes []map[string]string, err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	t := reflect.TypeOf(r)
	reflect_val := reflect.ValueOf(r)
	for i := 0; i < t.NumMethod(); i++ {

		namemethod := t.Method(i).Name

		handler := handlerReflect(reflect_val, namemethod)

		path_split := strings.Split(utils.CamelToKebab(namemethod), "_")
		apipath := path_split[0]
		route_method := ""
		if len(path_split) > 1 {
			route_method = path_split[1]
		}
		if apipath == "index" {
			api.Get("/", handler)
			apipath = ""
		}

		if route_method == "get" {
			// api.Get("/"+apipath, handler)
			// log.Print(route_method + " " + api_prefix + "/" + apipath)
			routes = append(routes, map[string]string{"method": route_method, "path": apipath})
			api.Get("/"+apipath+"/:id?", handler)
			if apipath != "" {
				// log.Print(route_method + " " + api_prefix + "/" + apipath + "/:id")
				routes = append(routes, map[string]string{"method": route_method, "path": apipath, "param": ":id"})
			}
		}
		if route_method == "delete" {
			api.Delete("/"+apipath+"/:id", handler)
			// log.Print(route_method + " " + api_prefix + "/" + apipath + "/:id")
			routes = append(routes, map[string]string{"method": route_method, "path": apipath, "param": ":id"})
		}
		if route_method == "put" {
			api.Put("/"+apipath+"/:id", handler)
			routes = append(routes, map[string]string{"method": route_method, "path": apipath, "param": ":id"})
			// log.Print(route_method + " " + api_prefix + "/" + apipath + "/:id")
		}
		if route_method == "post" {
			api.Post("/"+apipath, handler)
			routes = append(routes, map[string]string{"method": route_method, "path": apipath})
			// log.Print(route_method + " " + api_prefix + "/" + apipath)
		}
	}
	return

}

func (s *Suit) SetupGroups(api_prefix string, r interface{}, middleware ...fiber.Handler) (err error) {
	api_prefix = strings.ReplaceAll(api_prefix, "//", "/")
	api := s.Fiber.Group(api_prefix)

	reflect_val := reflect.ValueOf(r)

	if reflect_val.Kind() == reflect.Ptr {
		val := reflect_val.Elem()

		field := val.FieldByName("Suit")
		newVal := reflect.ValueOf(s)
		field.Set(newVal)
	}
	first_middle := reflect_val.MethodByName("Middleware")
	if first_middle.Kind() != reflect.Invalid {

		handler_middle := handlerReflect(reflect_val, "Middleware")
		api.Use(handler_middle)
	}
	for _, mid := range middleware {
		api.Use(mid)
	}
	routes, err := setupDynamicRoutes(api, r)
	if err != nil {
		return
	}
	re := regexp.MustCompile(`/+`)
	for _, val := range routes {
		join_path := re.ReplaceAllString(fmt.Sprintf("/%s/%s/%s", api_prefix, val["path"], val["param"]), "/")
		fmt.Printf("%s\n%s\t%s\n",strings.Repeat("-", 30), strings.ToUpper(val["method"]), join_path)

	}
	return
}

func (s *Suit) SetupRoutes(r interface{}) {
	api_prefix := "/api"
	if prefix := os.Getenv("API_PREFIX"); prefix != "" {
		api_prefix = prefix
	}

	s.SetupGroups(api_prefix, r)
}

func (s *Suit) Run() {

	HOST := os.Getenv("APP_HOST")
	PORT := os.Getenv("APP_PORT")
	if PORT == "" {
		PORT = "3000"
	}

	s.Fiber.Static("/", "./public")

	if err := s.Fiber.Listen(HOST + ":" + PORT); err != nil {
		log.Fatal(err)
	}
}
