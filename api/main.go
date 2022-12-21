package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/mesirendon/gredis/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading envvars: %s", err.Error())
	}

	app := fiber.New()

	router.SetRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
