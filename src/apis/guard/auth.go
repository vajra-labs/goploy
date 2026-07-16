package guard

import (
	"goploy/src/core/throw"
	"goploy/src/db/repos"
	"goploy/src/service"
	"goploy/src/types"

	"github.com/gofiber/fiber/v3"
)

type Guard struct {
	token   *service.TokenService
	queries *repos.Queries
}

func NewGuard(
	token *service.TokenService,
	queries *repos.Queries,
) *Guard {
	return &Guard{token, queries}
}

const userIDKey = "userID"

func (g *Guard) Auth() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		accessToken := ctx.Cookies(string(types.ACC_TOKEN))
		if accessToken == "" {
			return throw.UnauthorizedError(
				"Access token is required",
				"TOKEN_REQUIRED",
			)
		}
		payload, err := g.token.Verify(accessToken, types.ACC_TOKEN)
		if err != nil {
			return err
		}
		ctx.Locals(userIDKey, payload.Subject)
		return ctx.Next()
	}
}
