package docs

import (
	"goploy/src/apis/dtos"
	"goploy/src/core/apidoc"

	"github.com/danielgtaylor/huma/v2"
)

var userTags = []string{"User"}

// UserOpenApi registers OpenAPI 3.1 specifications for user endpoints.
func UserOpenApi(api huma.API) {
	r := api.OpenAPI()

	r.Paths["/api/user/me"] = &huma.PathItem{
		Get: &huma.Operation{
			Tags:        userTags,
			OperationID: "get-me",
			Summary:     "Get Current User",
			Description: "Get the profile details of the currently authenticated user.",
			Responses: apidoc.Response{
				"200": apidoc.JsonContent(
					api,
					dtos.UserResDto{},
					"User profile details",
				),
				"401": apidoc.ErrContent("User not authenticated"),
				"500": apidoc.ErrContent("Internal server error"),
			},
		},
	}
}
