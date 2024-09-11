package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krittakondev/goapisuit"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Route struct{
	// DB *gorm.DB
	// LimitPage int
	// RequireJwtAuth func(*fiber.Ctx) error
	Suit *goapisuit.Suit
}

func (r *Route) Index_get(c *fiber.Ctx) error{
	
	resp := Response{
		Message: "hello goapisuit",
	
		Data: []string{},
	};
	
	return c.JSON(resp);
}
