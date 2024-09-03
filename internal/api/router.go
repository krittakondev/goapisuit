package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Route struct{
	DB *gorm.DB
	PageLimit int
}

func (r *Route) Index_get(c *fiber.Ctx) error{
	
	resp := Response{
		Message: "hello goapisuit",
		Data: []string{},
	};
	
	return c.JSON(resp);
}
