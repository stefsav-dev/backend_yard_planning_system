// main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"backend_yard_planning_system/controllers"
	"backend_yard_planning_system/database"
)

func main() {

	database.ConnectDB()

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))

	yardController := controllers.NewYardController()

	api := app.Group("/api")
	{
		api.Get("/yard-plans", yardController.GetYardPlans)
		api.Post("/suggestion", yardController.GetSuggestion)
		api.Post("/placement", yardController.PlaceContainer)
		api.Post("/pickup", yardController.PickupContainer)
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Yard Planning System API",
			"version": "1.0.0",
		})
	})

	log.Fatal(app.Listen(":8080"))
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return ctx.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
