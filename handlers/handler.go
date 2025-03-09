package handlers

import(
	"context"
	"go-url-shortener/utils"
	"database/sql"
	"go-url-shortener/database"
	"github.com/gofiber/fiber/v2"	
)

func ShortenURL(c *fiber.Ctx) error{
	type Request struct{
		OriginalUrl string `json:"original_url"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err !=nil{
		return c.Status(400).JSON(fiber.Map{"error" : "invalid request"})
	}

	shortCode := utils.GenerateShortCode()

	_, err := database.DB.Exec(context.Background(), "INSERT INTO urls (original_url, short_url) values ($1,$2)", req.OriginalUrl, shortCode)
	if err!= nil{
		return c.Status(500).JSON(fiber.Map{"error" : "failed to convert url"})
	}

	return c.JSON(fiber.Map{"short_url" : "http://localhost:8080/" + shortCode})
}

func RedirectURL(c *fiber.Ctx) error{
	shortUrl := c.Params("shortURL")
	
	var originalUrl string

	err := database.DB.QueryRow(context.Background(), "SELECT original_url FROM urls WHERE short_url = $1", shortUrl).Scan(&originalUrl)
	if err == sql.ErrNoRows{
		return c.Status(404).JSON(fiber.Map{"error":"url not found"})
	} else if err != nil{
		return c.Status(500).JSON(fiber.Map{"error":"DataBase Error"})
	}

	return c.Redirect(originalUrl,301)
}

func CustomAliasGenerator(c *fiber.Ctx) error{
	type Request struct{
		OriginalURL string `json:"original_url"`
		CustomALias string `json:"custom_alias"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil{
		return c.Status(400).JSON(fiber.Map{"error" : " invalid request body"})
	}

	_, err := database.DB.Exec(context.Background(),"INSERT INTO urls (original_url, short_url) values ($1, $2)", req.OriginalURL, req.CustomALias)
	if err != nil{
		return c.Status(500).JSON(fiber.Map{"error" : "failed to convert url"})
	}

	return c.JSON(fiber.Map{"shortURL" : "http://localhost:8080/" + req.CustomALias})
}

func DeleteAll(c *fiber.Ctx) error{
	_, err := database.DB.Exec(context.Background(),"DELETE from urls")
	if err != nil{
		return c.Status(500).JSON(fiber.Map{"error":"couldn't delete all" + err.Error()},)
	}

	return c.JSON(fiber.Map{"status":"deleted"})
}