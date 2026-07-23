package handler

import (
	"strconv"

	"goploy/src/apis/dtos"
	"goploy/src/core/throw"
	"goploy/src/service"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Me handles GET /api/user/me.
func (h *UserHandler) Me(ctx fiber.Ctx) error {
	userIDStr, ok := ctx.Locals("userID").(string)
	if !ok || userIDStr == "" {
		return throw.UnauthorizedError(
			"User not authenticated",
			"NOT_AUTHENTICATED",
		)
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return throw.UnauthorizedError(
			"Invalid user ID format",
			"INVALID_USER_ID",
		)
	}
	user, err := h.userService.GetByID(ctx.Context(), userID)
	if err != nil {
		return err
	}
	isOwner := false
	if user.IsOwner != nil && *user.IsOwner == 1 {
		isOwner = true
	}
	return ctx.JSON(dtos.UserResDto{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		IsOwner:   isOwner,
	})
}
