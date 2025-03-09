package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-url-shortener/database"
	"go-url-shortener/utils"
	"time"

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

	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	//Insert click detail into click table

	_, err = database.DB.Exec(context.Background(), "INSERT INTO clicks (short_url, ip_address, user_agent) values ($1, $2, $3)", shortUrl, ipAddress,  userAgent)
	if err != nil{
		fmt.Println("Failed to log click analysis:", err)
	}

    //redirect user to original url
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

func GetAnalytics(c *fiber.Ctx) error {
	shortURL := c.Params("shortURL")

	//check for chahced data first in redis
	cachedData , err := database.RedisClient.Get(context.Background(), shortURL).Result()
	if err == nil{
		fmt.Println("returning from redis")
		var cachedResponse fiber.Map
		json.Unmarshal([]byte(cachedData), &cachedResponse)
		return c.JSON(cachedResponse)
	}

	//query the clicks db for the short url
	var totalClicks int
	err = database.DB.QueryRow(context.Background(), "SELECT COUNT(*) FROM clicks WHERE short_url = $1", shortURL).Scan(&totalClicks)
	if err != nil{
		return c.Status(500).JSON(fiber.Map{"error":"Couldn't fetch click analytics" + err.Error()})
	}

	//get last accessed time
	var lastAccessed sql.NullTime
	err = database.DB.QueryRow(context.Background(), "SELECT MAX(click_time) FROM clicks WHERE short_url = $1", shortURL).Scan(&lastAccessed)
	if err != nil{
		return c.Status(500).JSON(fiber.Map{"error":"Couldn't fetch click last accessed" + err.Error()})
	}

	//get most frequest user ip adress
	var modtFrequentIP string
	err = database.DB.QueryRow(context.Background(), "SELECT ip_address FROM clicks WHERE short_url = $1 GROUP BY ip_address ORDER BY COUNT(*) DESC LIMIT 1", shortURL).Scan(&modtFrequentIP)
	if err != nil{
		fmt.Println("Unable to fetch most frequest user", err)
	}

	response := fiber.Map{
		"short_url":      "http://localhost:8080/" + shortURL,
		"total_clicks":   totalClicks,
		"last_accessed":  lastAccessed.Time,
		"mostFrequentUser" : modtFrequentIP,
	}

	//Cache the analytics response in Redis (expire after 10 minutes)
	responseJSON, _ := json.Marshal(response)
	database.RedisClient.Set(context.Background(), shortURL, responseJSON, 10*time.Minute)

	// Return analytics response
	return c.JSON(response)
}