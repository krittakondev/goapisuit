package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/krittakondev/goapisuit"
	"github.com/krittakondev/goapisuit/internal/models"
	"github.com/krittakondev/goapisuit/pkg/utils"
	// "gorm.io/gorm"
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

func (r *Route) Index_get(c fiber.Ctx) error{
	
	resp := Response{
		Message: "hello goapisuit",
	
		Data: []string{},
	};
	
	return c.JSON(resp);
}

func (r *Route) Login_post(c fiber.Ctx) error{
	type Body struct {
		Username string 
		Password string
	}
	var body Body
	if err := c.Bind().Body(body); err != nil {
		return c.SendStatus(401)
	}
	
	var user models.Users
	tx := r.Suit.DB.Where("username=?", body.Username).Find(&user)
	if tx.Error != nil{
		return c.JSON(Response{
			Message: "DB error",
		});
	
	}
	if !utils.CheckPassword(user.PasswordEnc, body.Password){
		return c.JSON(Response{
			Message: "wrong password",
		});
		
	}
	
	token, err := utils.SignJwt(utils.JwtClaims{
		Sub: int(user.ID),
		Permissions: []string{"user"},
	})
	if err != nil {
		return c.JSON(Response{
			Message: err.Error(),
		})
	}
	resp := Response{
		Message: "Success Login",
	
		Data: map[string]string{
			"token": token,
		},
	};
	
	return c.JSON(resp);
}

