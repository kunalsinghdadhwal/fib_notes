package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kunalsinghdadhwal/fib_notes/utils"
)

const UserContextKey = "user"

type UserInfo struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func JWTMiddleware() fiber.Handler {

	return func(c *fiber.Ctx) error {

		accessToken := c.Cookies("access_token")
		if accessToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Access token is missing",
			})
		}

		claims, err := utils.ValidateJWT(accessToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid access token",
			})
		}
		c.Locals(UserContextKey, claims)
		return c.Next()

	}

}

func GetUserFromContext(c *fiber.Ctx) (*utils.JWTClaims, bool) {
	user, ok := c.Locals(UserContextKey).(*utils.JWTClaims)
	return user, ok
}

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
