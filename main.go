package main

import (
	"log"

	"go-url-shortener/database"
	"go-url-shortener/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"go-url-shortener/delete"
)

//global database connection

var db *pgx.Conn

func main() {
	//connect to PostgreSQL
	
	database.ConnectDB()

	//create a fibre app

	app := fiber.New()

	//test endpoint

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Server is running...🚀")
	})

	app.Post("/url/shorten", handlers.ShortenURL)
	app.Get("/:shortURL", handlers.RedirectURL)
	app.Post("/url/custom", handlers.CustomAliasGenerator)
	app.Delete("/delete", delete.DeleteAll) //only for personal db cleanup

	// Start server
	log.Fatal(app.Listen(":8080"))

}
