// Package middleware provides JWT authentication middleware for Fiber handlers.
//
// The middleware validates JWT tokens from the Authorization header and adds
// user claims to the request context. It includes both required and optional
// authentication middleware functions.
//
// Usage with Fiber:
//
//	app.Use("/protected", middleware.JWTMiddleware())
//	app.Use("/optional", middleware.OptionalJWTMiddleware())
//
// Token format expected in Authorization header:
//
//	Authorization: Bearer <jwt-token>
//
// The middleware automatically checks:
// - Token presence and format
// - Token signature validity
// - Token expiration (exp claim)
// - Token issued at time (iat claim)
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kunalsinghdadhwal/fib_notes/utils"
)

// UserContextKey is the key used to store user data in Fiber's locals
const UserContextKey = "user"

// UserInfo represents simplified user information extracted from JWT
type UserInfo struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

// JWTMiddleware validates JWT tokens and adds user claims to Fiber's locals
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		// Extract the token (remove "Bearer " prefix)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token required",
			})
		}

		// Validate the JWT token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Store user claims in Fiber's locals
		c.Locals(UserContextKey, claims)

		// Continue to the next handler
		return c.Next()
	}
}

// GetUserFromContext extracts user claims from Fiber's locals
func GetUserFromContext(c *fiber.Ctx) (*utils.JWTClaims, bool) {
	user, ok := c.Locals(UserContextKey).(*utils.JWTClaims)
	return user, ok
}

// GetUserInfoFromContext extracts simplified user info from Fiber's locals
func GetUserInfoFromContext(c *fiber.Ctx) (*UserInfo, bool) {
	claims, ok := GetUserFromContext(c)
	if !ok {
		return nil, false
	}

	userInfo := &UserInfo{
		ID:    claims.UserID,
		Name:  claims.Name,
		Email: claims.Email,
	}

	return userInfo, true
}

// OptionalJWTMiddleware validates JWT tokens if present but doesn't require them
// Useful for endpoints that can work with or without authentication
func OptionalJWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")

		// If no header, continue without authentication
		if authHeader == "" {
			return c.Next()
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Next()
		}

		// Extract the token (remove "Bearer " prefix)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return c.Next()
		}

		// Validate the JWT token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			// If token is invalid, continue without authentication
			return c.Next()
		}

		// Store user claims in Fiber's locals
		c.Locals(UserContextKey, claims)

		// Continue to the next handler
		return c.Next()
	}
}
