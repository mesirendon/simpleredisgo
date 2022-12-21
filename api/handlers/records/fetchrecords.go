package records

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mesirendon/gredis/datasource"
	"github.com/mesirendon/gredis/services/cache/mdkredis"
)

func FetchRecords(c *fiber.Ctx) error {
	regexp := c.Params("regexp")

	recordsDB := datasource.CreateClient(datasource.DefaultDatabase)
	recordsCache := mdkredis.New(recordsDB)

	defer func() {
		if err := recordsCache.Close(); err != nil {
			log.Fatalf("Error closing connection: %s", err.Error())
		}
	}()

	records := make([]string, 0)

	if err := recordsCache.GetRecordsByPattern(regexp, &records); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Can't get user records",
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Regexp": regexp,
		"Values": records,
	})
}
