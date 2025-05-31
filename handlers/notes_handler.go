package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kunalsinghdadhwal/fib_notes/db"
	"github.com/kunalsinghdadhwal/fib_notes/middleware"
	"github.com/kunalsinghdadhwal/fib_notes/models"
	"gorm.io/gorm"
)

// NotesHandler contains all notes related handlers
type NotesHandler struct{}

// NewNotesHandler creates a new notes handler instance
func NewNotesHandler() *NotesHandler {
	return &NotesHandler{}
}

// CreateNoteRequest represents the request body for creating a note
type CreateNoteRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=255"`
	Content string `json:"content" validate:"required,min=1"`
}

// UpdateNoteRequest represents the request body for updating a note
type UpdateNoteRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=255"`
	Content string `json:"content" validate:"required,min=1"`
}

// NoteResponse represents the response for note operations
type NoteResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt uint   `json:"created_at"`
	UpdatedAt uint   `json:"updated_at"`
}

// NotesListResponse represents the response for listing notes
type NotesListResponse struct {
	Notes []NoteResponse `json:"notes"`
	Count int            `json:"count"`
}

// CreateNote handles creating a new note
func (h *NotesHandler) CreateNote(c *fiber.Ctx) error {
	var req CreateNoteRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Title == "" || req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title and content are required",
		})
	}

	// Validate title length
	if len(req.Title) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title cannot exceed 255 characters",
		})
	}

	// Get user from JWT middleware
	claims, ok := middleware.GetUserFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Create new note
	note := models.Note{
		UserID:  claims.UserID,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := db.DB.Create(&note).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create note",
		})
	}

	// Return created note
	response := NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetNotes handles getting all notes for the authenticated user
func (h *NotesHandler) GetNotes(c *fiber.Ctx) error {
	// Get user from JWT middleware
	claims, ok := middleware.GetUserFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Parse query parameters for pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Build query
	query := db.DB.Where("user_id = ?", claims.UserID)

	// Add search functionality
	if search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var notes []models.Note
	if err := query.Order("created_at desc").Limit(limit).Offset(offset).Find(&notes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch notes",
		})
	}

	// Convert to response format
	var responses []NoteResponse
	for _, note := range notes {
		responses = append(responses, NoteResponse{
			ID:        note.ID,
			Title:     note.Title,
			Content:   note.Content,
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		})
	}

	// Get total count for pagination info
	var totalCount int64
	countQuery := db.DB.Model(&models.Note{}).Where("user_id = ?", claims.UserID)
	if search != "" {
		countQuery = countQuery.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	countQuery.Count(&totalCount)

	return c.JSON(fiber.Map{
		"notes":       responses,
		"count":       len(responses),
		"total":       totalCount,
		"page":        page,
		"limit":       limit,
		"total_pages": (totalCount + int64(limit) - 1) / int64(limit),
	})
}

// GetNote handles getting a specific note by ID (only for user's own notes)
func (h *NotesHandler) GetNote(c *fiber.Ctx) error {
	// Get note ID from URL params
	noteIDStr := c.Params("id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid note ID",
		})
	}

	// Get user from JWT middleware
	claims, ok := middleware.GetUserFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var note models.Note
	if err := db.DB.Where("id = ? AND user_id = ?", uint(noteID), claims.UserID).First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Note not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Return note
	response := NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}

	return c.JSON(response)
}

// UpdateNote handles updating a specific note (only for user's own notes)
func (h *NotesHandler) UpdateNote(c *fiber.Ctx) error {
	// Get note ID from URL params
	noteIDStr := c.Params("id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid note ID",
		})
	}

	var req UpdateNoteRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Title == "" || req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title and content are required",
		})
	}

	// Validate title length
	if len(req.Title) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title cannot exceed 255 characters",
		})
	}

	// Get user from JWT middleware
	claims, ok := middleware.GetUserFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Find and update note (only if it belongs to the user)
	var note models.Note
	if err := db.DB.Where("id = ? AND user_id = ?", uint(noteID), claims.UserID).First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Note not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Update note
	note.Title = req.Title
	note.Content = req.Content

	if err := db.DB.Save(&note).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update note",
		})
	}

	// Return updated note
	response := NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}

	return c.JSON(response)
}

// DeleteNote handles deleting a specific note (only for user's own notes)
func (h *NotesHandler) DeleteNote(c *fiber.Ctx) error {
	// Get note ID from URL params
	noteIDStr := c.Params("id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid note ID",
		})
	}

	// Get user from JWT middleware
	claims, ok := middleware.GetUserFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Find note (only if it belongs to the user)
	var note models.Note
	if err := db.DB.Where("id = ? AND user_id = ?", uint(noteID), claims.UserID).First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Note not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Delete note
	if err := db.DB.Delete(&note).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete note",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Note deleted successfully",
	})
}
