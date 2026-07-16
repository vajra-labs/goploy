package guard

import (
	"strconv"

	"goploy/src/core/throw"
	"goploy/src/db/repos"

	"github.com/gofiber/fiber/v3"
)

const orgIDHeader = "X-Org-ID"

// Authorize checks if the user has the required permission.
// Owner (is_owner = 1) bypasses all permission checks.
// Usage:
//
//	guard.Authorize("project", "create")
func (g *Guard) Authorize(resource, action string) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		userIDStr, ok := ctx.Locals(userIDKey).(string)
		if !ok || userIDStr == "" {
			return throw.UnauthorizedError(
				"User not authenticated",
				"NOT_AUTHENTICATED",
			)
		}
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return throw.UnauthorizedError("Invalid user ID", "INVALID_USER_ID")
		}
		// Fetch user to check is_owner
		user, err := g.queries.GetUserByID(ctx.Context(), userID)
		if err != nil {
			return throw.UnauthorizedError("User not found", "USER_NOT_FOUND")
		}
		// Owner bypasses all permission checks
		if user.IsOwner != nil && *user.IsOwner == 1 {
			return ctx.Next()
		}
		// Non-owner — org ID required
		orgIDStr := ctx.Get(orgIDHeader)
		if orgIDStr == "" {
			return throw.UnauthorizedError(
				"Organization ID is required",
				"ORG_ID_REQUIRED",
			)
		}
		orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
		if err != nil {
			return throw.UnauthorizedError(
				"Invalid organization ID",
				"INVALID_ORG_ID",
			)
		}
		// Fetch final permissions (group + GRANT - DENY)
		key := resource + ":" + action
		permissions, err := g.queries.GetUserFinalPermissions(
			ctx.Context(),
			repos.GetUserFinalPermissionsParams{
				UserID: userID,
				OrgID:  orgID,
			},
		)
		if err != nil {
			return throw.InternalServerError(
				"Failed to check permissions",
				"PERMISSION_CHECK_ERROR",
				throw.WithCause(err),
			)
		}
		for _, p := range permissions {
			if p == key {
				return ctx.Next()
			}
		}
		return throw.ForbiddenError("Permission denied", "PERMISSION_DENIED")
	}
}
