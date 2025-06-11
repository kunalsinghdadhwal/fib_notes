package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv"
	"github.com/kunalsinghdadhwal/fib_notes/db"
	"github.com/kunalsinghdadhwal/fib_notes/routes"
	"github.com/watchakorn-18k/scalar-go"
)

func main() {

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	if err := db.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	routes.SetupAuthRoutes(app)
	routes.SetupNotesRoutes(app)
	routes.SetupSeederRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/api/reference", fiber.StatusMovedPermanently)
	})

	app.Use("/api/reference", func(c *fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL:  "./docs/swagger.yaml",
			Theme:    scalar.ThemeId("saturn"),
			DarkMode: true,
		})

		if err != nil {
			return err
		}
		c.Type("html")
		return c.SendString(htmlContent)
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	port := os.Getenv("PORT")

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
