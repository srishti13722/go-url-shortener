package authhandler

import (
	"go-url-shortener/auth"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error{
	type Credentials struct{
		Username string `json:"username"`
		Password string `json:"password"`
	}

	req := new(Credentials)

	if err := c.BodyParser(req); err != nil{
		return c.Status(400).JSON(fiber.Map{"error":"invalid request"})
	}

	// Hardcoded login credentials (Replace with database check)
	if req.Username != "admin" || req.Password != "password" {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(req.Username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{"token": token})
}