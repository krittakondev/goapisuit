package routes

import "github.com/gofiber/fiber/v2"

type Response struct{
	Code int
	Message string
	Data interface{}
}

type Route struct{
}

func (r *Route) MainRoute(c *fiber.Ctx) error{
	
	resp := Response{
		Message: "hello goapisuit",
		Data: []string{},
	};
	
	return c.JSON(resp);
}
