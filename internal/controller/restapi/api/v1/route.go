package v1

import (
	"github.com/danielgtaylor/huma/v2"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
)

func RegisterGroup(public huma.API, protected huma.API, h *handler.Handlers) {
	publicGroup := huma.NewGroup(public, "/v1")
	protectedGroup := huma.NewGroup(protected, "/v1")

	modifier := func(op *huma.Operation) {
		op.OperationID = op.OperationID + "-v1"
		op.Summary = op.Summary + " (V1)"
	}
	publicGroup.UseSimpleModifier(modifier)
	protectedGroup.UseSimpleModifier(modifier)

	RegisterRoutes(publicGroup, protectedGroup, h)
}

func RegisterRoutes(public huma.API, protected huma.API, h *handler.Handlers) {
	RegisterFileGroup(public, protected, h)
	RegisterBookGroup(public, protected, h)
	RegisterBookPageGroup(public, protected, h)
	RegisterPetGroup(public, protected, h)
	RegisterReviewGroup(public, protected, h)
	RegisterUserGroup(public, protected, h)
}
