package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kunalsinghdadhwal/fib_notes/handlers"
	"github.com/kunalsinghdadhwal/fib_notes/middleware"
)

func SetupNotesRoutes(app *fiber.App) {
	notesHandler := handlers.NewNotesHandler()

	notes := app.Group("/notes", middleware.JWTMiddleware())

	notes.Post("/", notesHandler.CreateNote)
	notes.Get("/", notesHandler.GetNotes)
	notes.Get("/:id", notesHandler.GetNote)
	notes.Put("/:id", notesHandler.UpdateNote)
	notes.Delete("/:id", notesHandler.DeleteNote)
}
