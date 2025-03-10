package main

import (
	"log"

	"go-url-shortener/database"
	"go-url-shortener/handlers"
	"github.com/gofiber/fiber/v2"
	"go-url-shortener/delete"
	"go-url-shortener/middleware"
	"go-url-shortener/tasks"
)

func main() {
	//connect to PostgreSQL and redis	
	database.ConnectDB()
    database.ConnectRedis()

	// background cleanup job
	go tasks.DeleteExpiredUrls()
	
	//create a fibre app
	app := fiber.New()

	//create rate limiter
	limiter := middleware.RateLimiter(10,"M")

	app.Post("/url/shorten",limiter, handlers.ShortenURL)
	app.Get("/:shortURL",limiter, handlers.RedirectURL)
	app.Post("/url/custom",limiter, handlers.CustomAliasGenerator)
	app.Get("/url/stats/:shortURL",limiter, handlers.GetAnalytics)
	app.Delete("/delete", delete.DeleteAll) //only for personal db cleanup

	// Start server
	log.Fatal(app.Listen(":8080"))

}
