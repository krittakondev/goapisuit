package maketemplate

import (
	"fmt"

	"github.com/krittakondev/goapisuit/v2/pkg/utils"
)



var templateRouter = `package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/krittakondev/goapisuit/v2"
)

type Response struct {
	Code    int         `+"`"+`json:"code"`+"`"+`
	Message string      `+"`"+`json:"message"`+"`"+`
	Data    interface{} `+"`"+`json:"data"`+"`"+`
}

type Route struct{
	Suit *goapisuit.Suit
}

func (r *Route) Middleware(c fiber.Ctx) error{
	
	// middleware for group
	
	return c.Next()
}

func (r *Route) Index_get(c fiber.Ctx) error{
	
	resp := Response{
		Message: "hello goapisuit",
	
		Data: []string{},
	};
	
	return c.JSON(resp);
}
`

var templateServer = `package main
import (
	"log"

	"github.com/krittakondev/goapisuit/v2"
	"{{.ProjectName}}/internal/routes"
	"{{.ProjectName}}/internal/setup"
)




func main(){
	suit, err := goapisuit.New("{{.ProjectName}}")
	if err != nil{
		log.Fatal(err)
	}
	if suit.Config.DbConnection == "mysql"{
	
		sqlDB, err := suit.DB.DB()
		if err != nil {
			log.Println("failed to get database object:", err)
		}
		defer sqlDB.Close()
	}
	suit.SetupRoutes(&routes.Route{})
	setup.GroupsSetup(suit)

	suit.Run()
}
`


type EnvStruct struct{
	AppName string
	AppPort string
	JwtSecret string
	DbUsername string
	DbPassword string
	DbDatabase string
	DbPort string
	DbHost string
	
}

var templateEnv = `
# APP

APP_NAME={{.AppName}}
APP_HOST=0.0.0.0
APP_PORT=3002

# API
API_PREFIX=/api
API_LIMIT_PAGE=20


# jwt

JWT_SECRET={{.JwtSecret}}
JWT_EXPIRE=1d


# database

DB_CONNECTION=mysql #now support just mysql
DB_HOST={{.DbHost}}
DB_PORT={{.DbPort}}
DB_DATABASE={{.DbDatabase}}
DB_USERNAME={{.DbUsername}}
DB_PASSWORD={{.DbPassword}}


#BCRYPT_ROUNDS=12
`

var templatePublicIndex = `
<h1>public goAPIsuit</h1>

<h2>fetch from /api</h2>
<div id="result"></div>

<script>
	async function getData(){
		const resp = await fetch("/api")
		const data = await resp.text()
		document.getElementById("result").innerHTML = data
	}
	getData()

</script>
`


var templateMakeRouter = `package routes

import (
	"github.com/gofiber/fiber/v3"
	"{{.PathProject}}/internal/models"
	"gorm.io/gorm"
)

type Response{{.Name}} struct {
	Code    int         `+"`"+`json:"code"`+"`"+`
	Message string      `+"`"+`json:"message"`+"`"+`
	Data    interface{} `+"`"+`json:"data,omitempty"`+"`"+`
}


func (r *Route) {{.Name}}_get(c fiber.Ctx) error {
	// uncomment below for protect route with jwt  header["Token"] 
	// if err :=  r.Suit.RequireJwtAuth(c); err!=nil{
	// 	return err
	// }
	// fmt.Printf("%+v", c.Locals("user"))
	limit := fiber.Query[int](c, "limit", r.Suit.LimitPage)
	page := fiber.Query[int](c, "page", 1)
	code := 0
	message := "response from {{.Name}}"
	id := c.Params("id", "")


	var result []models.{{.ModelName}}
	var tx *gorm.DB
	if id != ""{
		tx = r.Suit.DB.Where("id=?", id).Find(&result)
		message = "response from {{.Name}} id "+ id
	}else{
		tx = r.Suit.DB.Limit(limit).Offset((page-1)*limit).Order("created_at DESC").Find(&result)
	}

	if tx.Error != nil {
		code = -1
		message = "database error"
	}
	resp := Response{{.Name}}{
		Code:    code,
		Message: message,
		Data:    result,
	}

	return c.JSON(resp)
}

func (r *Route) {{.Name}}_post(c fiber.Ctx) error {
	// uncomment below for protect route with jwt  header["Token"] 
	// if err :=  r.Suit.RequireJwtAuth(c); err!=nil{
	// 	return err
	// }
	// fmt.Printf("%+v", c.Locals("user"))

	var body models.{{.ModelName}};
	
	if err := c.Bind().Body(&body); err != nil {
		return c.JSON(Response{{.Name}}{
			Code: -1,
			Message: err.Error(),
			Data: nil,
		})
	}
	body.ID = 0
	tx := r.Suit.DB.Create(&body)
	if tx.Error != nil {
		return c.JSON(Response{{.Name}}{
			Code:    -1,
			Message: "db error",
			Data:    nil,
		})
	}

	return c.JSON(Response{{.Name}}{
		Code:    0,
		Message: "Create {{.Name}} Success",
		Data:    body,
	})

}

func (r *Route) {{.Name}}_delete(c fiber.Ctx) error {
	// uncomment below for protect route with jwt  header["Token"] 
	// if err :=  r.Suit.RequireJwtAuth(c); err!=nil{
	// 	return err
	// }
	// fmt.Printf("%+v", c.Locals("user"))
	id := c.Params("id", "")
	if id == "" {
		return c.JSON(Response{{.Name}}{
			Code:    -1,
			Message: "{{.Name}} Id error",
			Data:    nil,
		})
	}

	tx := r.Suit.DB.Delete(&models.{{.ModelName}}{}, id)
	if tx.Error != nil {
		return c.JSON(Response{{.Name}}{
			Code:    -1,
			Message: "db error",
			Data:    nil,
		})
	}

	return c.JSON(Response{{.Name}}{
		Code:    0,
		Message: "Delete {{.Name}} id "+id+" Success",
		Data:    nil,
	})
}

func (r *Route) {{.Name}}_put(c fiber.Ctx) error {
	// uncomment below for protect route with jwt  header["Token"] 
	// if err :=  r.Suit.RequireJwtAuth(c); err!=nil{
	// 	return err
	// }
	// fmt.Printf("%+v", c.Locals("user"))
	id := c.Params("id", "")
	if id == "" {
		return c.JSON(Response{{.Name}}{
			Code:    -1,
			Message: "{{.Name}} Id error",
			Data:    nil,
		})
	}
	var body models.{{.ModelName}}
	err := c.Bind().Body(&body)
	if err != nil {
		return c.JSON(Response{{.Name}}{
			Code:    -1,
			Message: err.Error(),
			Data:    nil,
		})
	}
	tx := r.Suit.DB.Model(&models.{{.ModelName}}{}).Where("id=?", id).Updates(body)
	if tx.Error != nil {
		return c.JSON(Response{{.Name}}{
			Code:    -1,
			Message: "db error",
			Data:    nil,
		})
	}

	return c.JSON(Response{{.Name}}{
		Code:    0,
		Message: "Update {{.Name}} id "+id+" Success",
		Data:    nil,
	})
}
`

var templateMakeModel = `package models

import "gorm.io/gorm"


  // Edit this models
type {{.Name}} struct {
  gorm.Model
  // Name       string     `+"`"+`json:"name" gorm:"index"`+"`"+`
}
`

var templateDbMigrate = `
package main

import (
	"{{.PathProject}}/internal/models"
	"github.com/krittakondev/v2/goapisuit/database"
)

func main(){
	db, err := database.MysqlConnect()
	if err != nil{
		panic(err)
	}
	db.AutoMigrate(&models.{{.Name}}{})
}
`

var templateDockerfile = `FROM golang:alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app cmd/server.go

CMD ["app"]
`

var templateDockerCompose = `
services:

  app:
    build: .
    restart: unless-stopped
    network_mode: host
    volumes: 
      - .env:/usr/src/app/.env
    depends_on:
      db:
        condition: service_started
    

  db:
    image: "mysql:8.4.3"
    restart: always
    environment:
      MYSQL_USER: {{.DbUsername}}
      MYSQL_PASSWORD: {{.DbPassword}}
      MYSQL_ROOT_PASSWORD: {{.DbPassword}}
      MYSQL_DATABASE: {{.DbDatabase}}
    ports:
      - "{{.DbPort}}:3306"
    networks:
      - backend
    volumes:
      - ./db_suit:/var/lib/mysql

networks:
  backend:
`

type GroupsLoader struct{
	ImportRouteGroup string
	SetupGroups string
}

func CreateTemplateGroupsSetupCall(groups_path []string) (arr []string){
	for _, val := range groups_path{
		name := utils.PathToModelFormatName(val)
		if name != ""{
			arr = append(arr, fmt.Sprintf("\tsuit.SetupGroups(suit.Config.ApiPrefix+\"%s\", &route_%s.Route{})", val, name))
		}
	}
	return
}

func CreateTemplateGroupsSetupImport(project_name string, groups_path []string)(arr []string){
	for _, val := range groups_path{
		name := utils.PathToModelFormatName(val)
		if name != ""{
			arr = append(arr, fmt.Sprintf("\troute_%s \"%s/internal/routes%s\"", name, project_name, val))
		}
	}
	return
}
var templateGroupsSetup = `package setup

import (
{{.ImportRouteGroup}}
	"github.com/krittakondev/goapisuit/v2"
)

func GroupsSetup(suit *goapisuit.Suit){
{{.SetupGroups}}
}
`

