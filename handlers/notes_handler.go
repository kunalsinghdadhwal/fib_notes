package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kunalsinghdadhwal/fib_notes/db"
	"github.com/kunalsinghdadhwal/fib_notes/middleware"
	"github.com/kunalsinghdadhwal/fib_notes/models"
	"gorm.io/gorm"
)

type NotesHandler struct{}

func NewNotesHandler() *NotesHandler {
	return &NotesHandler{}
}

type CreateNoteRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=255" example:"My First Note"`
	Content string `json:"content" validate:"required,min=1" example:"This is the content of my first note"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=255" example:"Updated Note Title"`
	Content string `json:"content" validate:"required,min=1" example:"This is the updated content of my note"`
}

type NoteResponse struct {
	ID        uint   `json:"id" example:"1"`
	Title     string `json:"title" example:"My First Note"`
	Content   string `json:"content" example:"This is the content of my first note"`
	CreatedAt uint   `json:"created_at" example:"1672531200"`
	UpdatedAt uint   `json:"updated_at" example:"1672531200"`
}

type NotesListResponse struct {
	Notes      []NoteResponse `json:"notes"`
	Count      int            `json:"count" example:"5"`
	Total      int64          `json:"total" example:"25"`
	Page       int            `json:"page" example:"1"`
	Limit      int            `json:"limit" example:"10"`
	TotalPages int64          `json:"total_pages" example:"3"`
}

// CreateNote handles creating a new note
// @Summary Create a new note
// @Description Create a new note for the authenticated user
// @Tags Notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateNoteRequest true "Note creation data"
// @Success 201 {object} NoteResponse "Note created successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "User not authenticated"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notes [post]
func (h *NotesHandler) CreateNote(c *fiber.Ctx) error {
	var req CreateNoteRequest

	claims, ok := middleware.GetUserFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Title == "" || req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title and content are required",
		})
	}

	if len(req.Title) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title cannot exceed 255 characters",
		})
	}

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
// @Summary Get all notes
// @Description Get all notes for the authenticated user with pagination and search
// @Tags Notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of notes per page" default(10)
// @Param search query string false "Search term for notes title or content"
// @Success 200 {object} NotesListResponse "Notes retrieved successfully"
// @Failure 401 {object} map[string]string "User not authenticated"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notes [get]
func (h *NotesHandler) GetNotes(c *fiber.Ctx) error {
	claims, ok := middleware.GetUserFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

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

	query := db.DB.Where("user_id = ?", claims.UserID)

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
// @Summary Get a specific note
// @Description Get a specific note by ID for the authenticated user
// @Tags Notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Note ID"
// @Success 200 {object} NoteResponse "Note retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid note ID"
// @Failure 401 {object} map[string]string "User not authenticated"
// @Failure 404 {object} map[string]string "Note not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notes/{id} [get]
func (h *NotesHandler) GetNote(c *fiber.Ctx) error {

	claims, ok := middleware.GetUserFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	noteIDStr := c.Params("id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid note ID",
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
// @Summary Update a specific note
// @Description Update a specific note by ID for the authenticated user
// @Tags Notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Note ID"
// @Param request body UpdateNoteRequest true "Note update data"
// @Success 200 {object} NoteResponse "Note updated successfully"
// @Failure 400 {object} map[string]string "Invalid request data or note ID"
// @Failure 401 {object} map[string]string "User not authenticated"
// @Failure 404 {object} map[string]string "Note not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notes/{id} [put]
func (h *NotesHandler) UpdateNote(c *fiber.Ctx) error {
	claims, ok := middleware.GetUserFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	noteIDStr := c.Params("id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid note ID",
		})
	}

	var req UpdateNoteRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Title == "" || req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title and content are required",
		})
	}

	if len(req.Title) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title cannot exceed 255 characters",
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
// @Summary Delete a specific note
// @Description Delete a specific note by ID for the authenticated user
// @Tags Notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Note ID"
// @Success 200 {object} map[string]string "Note deleted successfully"
// @Failure 400 {object} map[string]string "Invalid note ID"
// @Failure 401 {object} map[string]string "User not authenticated"
// @Failure 404 {object} map[string]string "Note not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notes/{id} [delete]
func (h *NotesHandler) DeleteNote(c *fiber.Ctx) error {
	claims, ok := middleware.GetUserFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	noteIDStr := c.Params("id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid note ID",
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
