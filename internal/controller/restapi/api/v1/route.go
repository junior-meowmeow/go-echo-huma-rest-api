package v1

import (
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterGroup(api huma.API, h *handler.Handlers) {
	v1Group := huma.NewGroup(api, "/v1")

	v1Group.UseSimpleModifier(func(op *huma.Operation) {
		op.OperationID = op.OperationID + "-v1"
		op.Summary = op.Summary + " (V1)"
	})

	RegisterRoutes(v1Group, h)
}

func RegisterRoutes(api huma.API, h *handler.Handlers) {
	RegisterReviewGroup(api, h)
	RegisterFileGroup(api, h)
	RegisterBookGroup(api, h)
	RegisterBookPageGroup(api, h)
	RegisterPetGroup(api, h)
}
