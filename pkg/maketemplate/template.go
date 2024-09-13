package maketemplate



var templateRouter = `package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krittakondev/goapisuit"
)

type Response struct {
	Code    int         `+"`"+`json:"code"`+"`"+`
	Message string      `+"`"+`json:"message"`+"`"+`
	Data    interface{} `+"`"+`json:"data"`+"`"+`
}

type Route struct{
	Suit *goapisuit.Suit
}

func (r *Route) Index_get(c *fiber.Ctx) error{
	
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

	"github.com/krittakondev/goapisuit"
	"{{.ProjectName}}/internal/routes"
)




func main(){
	suit, err := goapisuit.New("{{.ProjectName}}/goapisuit")
	if err != nil{
		log.Fatal(err)
	}
	suit.Run(&routes.Route{Suit: suit})
}
`


type envStruct struct{
	AppName string
	JwtSecret string
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
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=goapisuit
DB_USERNAME=root
DB_PASSWORD=


BCRYPT_ROUNDS=12
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

	"github.com/gofiber/fiber/v2"
	"{{.PathProject}}/internal/models"
	"gorm.io/gorm"
)

type Response{{.Name}} struct {
	Code    int         `+"`"+`json:"code"`+"`"+`
	Message string      `+"`"+`json:"message"`+"`"+`
	Data    interface{} `+"`"+`json:"data,omitempty"`+"`"+`
}

func (r *Route) {{.Name}}_get(c *fiber.Ctx) error {
	// uncomment below for protect route with jwt  header["Token"] 
	// if err :=  r.Suit.RequireJwtAuth(c); err!=nil{
	// 	return err
	// }
	// fmt.Printf("%+v", c.Locals("user"))
	limit := c.QueryInt("limit", r.Suit.LimitPage)
	page := c.QueryInt("page", 1)
	code := 0
	message := "response from {{.Name}}"
	id := c.Params("id", "")


	var result []models.{{.Name}}
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

func (r *Route) {{.Name}}_post(c *fiber.Ctx) error {
	// uncomment below for protect route with jwt  header["Token"] 
	// if err :=  r.Suit.RequireJwtAuth(c); err!=nil{
	// 	return err
	// }
	// fmt.Printf("%+v", c.Locals("user"))

	var body models.{{.Name}};
	
	if err := c.BodyParser(&body); err != nil {
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

func (r *Route) {{.Name}}_delete(c *fiber.Ctx) error {
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

	tx := r.Suit.DB.Delete(&models.{{.Name}}{}, id)
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

func (r *Route) {{.Name}}_put(c *fiber.Ctx) error {
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
	var body models.{{.Name}}
	err := c.BodyParser(&body)
	if err != nil {
		return c.JSON(Response{{.Name}}{
			Code:    -1,
			Message: err.Error(),
			Data:    nil,
		})
	}
	tx := r.Suit.DB.Model(&models.{{.Name}}{}).Where("id=?", id).Updates(body)
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
