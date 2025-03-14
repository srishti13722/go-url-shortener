package main

import (
	"log"

	"go-url-shortener/authhandler"
	"go-url-shortener/authmiddleware"
	"go-url-shortener/database"
	"go-url-shortener/delete"
	"go-url-shortener/handlers"
	"go-url-shortener/middleware"
	"go-url-shortener/tasks"

	"github.com/gofiber/fiber/v2"
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
	limiter := middleware.RateLimiter(10, "M")

	// Login API (Public)
	app.Post("/login", authhandler.Login)

	app.Post("/url/shorten", authmiddleware.AuthMiddlerware, limiter, handlers.ShortenURL)
	app.Get("/:shortURL", limiter, handlers.RedirectURL)
	app.Post("/url/custom", authmiddleware.AuthMiddlerware, limiter, handlers.CustomAliasGenerator)
	app.Get("/url/stats/:shortURL", authmiddleware.AuthMiddlerware, limiter, handlers.GetAnalytics)
	
	//only for personal db cleanup
	app.Delete("/delete", authmiddleware.AuthMiddlerware, delete.DeleteAll) 

	// Start server
	log.Fatal(app.Listen(":8080"))

}
