package main

import (
	"log"

	"go-url-shortener/database"
	"go-url-shortener/handlers"
	"github.com/gofiber/fiber/v2"
	"go-url-shortener/delete"
)

//global database connection

func main() {
	//connect to PostgreSQL and redis
	
	database.ConnectDB()
    database.ConnectRedis()
	//create a fibre app

	app := fiber.New()

	//test endpoint

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Server is running...🚀")
	})

	app.Post("/url/shorten", handlers.ShortenURL)
	app.Get("/:shortURL", handlers.RedirectURL)
	app.Post("/url/custom", handlers.CustomAliasGenerator)
	app.Get("/url/stats/:shortURL", handlers.GetAnalytics)
	app.Delete("/delete", delete.DeleteAll) //only for personal db cleanup

	// Start server
	log.Fatal(app.Listen(":8080"))

}
