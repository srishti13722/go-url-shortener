package authmiddleware

import(
	"go-url-shortener/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddlerware(c *fiber.Ctx) error{
	//get token from authorization header
	tokenString := c.Get("Authorization")

	_ , err := auth.ValidateJWT(tokenString)
	if err != nil{
		return c.Status(401).JSON(fiber.Map{
			"error" : "Unauthorized - Invalid Token",
		})
	}

	//continue if token is valid
	return c.Next()
}