package records

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/mesirendon/gredis/datasource"
	"github.com/mesirendon/gredis/services/cache"
	"github.com/mesirendon/gredis/services/cache/mdkredis"
)

type Record struct {
	Success bool
	Key     string
}

const ITER int = 100

var (
	recordsDB    redis.Cmdable
	recordsCache cache.ICache
)

func setRecord(i, j int) Record {
	key := fmt.Sprintf("user-%d-%d", i, j)
	record := Record{
		Success: true,
		Key:     key,
	}

	if err := recordsCache.Set(key, key, time.Hour); err != nil {
		record.Success = false
	}

	return record
}

func ReloadRecords(c *fiber.Ctx) error {
	recordsDB = datasource.CreateClient(datasource.DefaultDatabase)
	recordsCache = mdkredis.New(recordsDB)

	defer func() {
		if err := recordsCache.Close(); err != nil {
			log.Fatalf("Error closing db: %s", err.Error())
		}
	}()

	if err := recordsCache.FlushDB(); err != nil {
		log.Fatalf("Error flushing db: %s", err.Error())
	}

	recordsChan := make(chan Record)

	go func() {
		wg := sync.WaitGroup{}
		for i := 0; i < ITER; i++ {
			for j := 0; j < ITER; j++ {
				wg.Add(1)

				go func(i, j int) {
					defer wg.Done()
					recordsChan <- setRecord(i, j)
				}(i, j)
			}
		}
		wg.Wait()
		close(recordsChan)
	}()

	successData := make([]string, 0)
	failData := make([]string, 0)

	for r := range recordsChan {
		if r.Success {
			successData = append(successData, r.Key)
		} else {
			failData = append(failData, r.Key)
		}
	}

	if len(failData) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":        "Error setting records",
			"failed_rules": failData,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"operation":  "Reload records",
		"successful": successData,
	})
}
