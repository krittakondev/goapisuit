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
`

