package docs

import (
	"net/http"

	"goploy/src/apis/dtos"
	"goploy/src/core/apidoc"

	"github.com/danielgtaylor/huma/v2"
)

func AuthOpenApi(api huma.API) {
	r := api.OpenAPI()
	tags := []string{"Authentication"}

	r.Paths["/api/auth/setup"] = &huma.PathItem{
		Get: &huma.Operation{
			Tags:        tags,
			OperationID: "auth-setup",
			Summary:     "Check Setup Status",
			Description: "Returns whether the owner has already been registered.",
			Responses: apidoc.Responses(
				apidoc.JsonContent(api, http.StatusOK, struct {
					IsOwnerPresent bool `json:"isOwnerPresent" doc:"True if owner is already registered"`
				}{}, "Setup status"),
			),
		},
	}

	r.Paths["/api/auth/register"] = &huma.PathItem{
		Post: &huma.Operation{
			Tags:        tags,
			OperationID: "register-user",
			Summary:     "Register Owner",
			Description: "Registers the first user as the application owner. Fails if owner already exists.",
			RequestBody: apidoc.Body(
				api,
				dtos.RegisterDto{},
				true,
				"Owner registration details",
			),
			Responses: apidoc.Responses(
				apidoc.TextContent(
					http.StatusCreated,
					"Owner registered successfully",
				),
				apidoc.ErrContent(
					http.StatusBadRequest,
					"Invalid request body",
				),
				apidoc.ErrContent(http.StatusConflict, "Owner already exists"),
				apidoc.ErrContent(
					http.StatusInternalServerError,
					"Internal server error",
				),
			),
		},
	}

	r.Paths["/api/auth/login"] = &huma.PathItem{
		Post: &huma.Operation{
			Tags:        tags,
			OperationID: "login-user",
			Summary:     "Login",
			Description: "Authenticates a user using email and password and returns access + refresh tokens.",
			RequestBody: apidoc.Body(
				api,
				dtos.LoginDto{},
				true,
				"Login credentials",
			),
			Responses: apidoc.Responses(
				apidoc.TextContent(
					http.StatusOK,
					"Access and refresh tokens set in HTTP-only cookies",
				),
				apidoc.ErrContent(
					http.StatusBadRequest,
					"Invalid request body",
				),
				apidoc.ErrContent(
					http.StatusUnauthorized,
					"Invalid email or password",
				),
				apidoc.ErrContent(
					http.StatusInternalServerError,
					"Internal server error",
				),
			),
		},
	}

	r.Paths["/api/auth/refresh"] = &huma.PathItem{
		Post: &huma.Operation{
			Tags:        tags,
			OperationID: "refresh-token",
			Summary:     "Refresh Access Token",
			Description: "Regenerates a new HTTP-only access token using the valid HTTP-only refresh token cookie.",
			Responses: apidoc.Responses(
				apidoc.TextContent(
					http.StatusOK,
					"Access token refreshed successfully in cookies",
				),
				apidoc.ErrContent(
					http.StatusBadRequest,
					"Refresh token is required",
				),
				apidoc.ErrContent(
					http.StatusUnauthorized,
					"Refresh token is invalid or expired",
				),
				apidoc.ErrContent(
					http.StatusInternalServerError,
					"Internal server error",
				),
			),
		},
	}

	r.Paths["/api/auth/logout"] = &huma.PathItem{
		Post: &huma.Operation{
			Tags:        tags,
			OperationID: "logout-user",
			Summary:     "Logout",
			Description: "Blacklists the current refresh token and clears both authentication cookies.",
			Parameters:  apidoc.QueryParams(dtos.LogoutDto{}),
			Responses: apidoc.Responses(
				apidoc.TextContent(http.StatusOK, "Logged out successfully"),
				apidoc.ErrContent(
					http.StatusBadRequest,
					"Refresh token is required",
				),
				apidoc.ErrContent(
					http.StatusInternalServerError,
					"Internal server error",
				),
			),
		},
	}
}
