package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/mesirendon/gredis/handlers/records"
)

func loadRecordsRoutes(r fiber.Router) {
	rl := r.Group("/records")

	rl.Get("/", timeout.New(records.ReloadRecords, 30*time.Second))
	rl.Get("/:regexp", timeout.New(records.FetchRecords, 30*time.Second))

}
