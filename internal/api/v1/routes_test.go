package v1_test

// import (
// 	"testing"

// 	"github.com/danielgtaylor/huma/v2/humatest"

// 	. "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/api/v1"
// )

// func TestPostReviewInvalid(t *testing.T) {
// 	_, testAPI := humatest.New(t)
// 	RegisterRoutes(testAPI)

// 	// rating 10 violates your model constraints (1–5)
// 	resp := testAPI.Post("/reviews", map[string]any{
// 		"rating": 10,
// 	})
// 	if resp.Code != 422 {
// 		t.Fatalf("expected 422 Unprocessable Entity, got %d", resp.Code)
// 	}
// }

// TODO: Mock MongoDB to be able to test Read/Write
// func TestPostReviewValid(t *testing.T) {
// 	_, testAPI := humatest.New(t)
// 	api.RegisterRoutes(testAPI)

// 	resp := testAPI.Post("/reviews", map[string]any{
// 		"author":  "daniel",
// 		"rating":  5,
// 		"message": "Great!",
// 	})
// 	if resp.Code != 201 {
// 		t.Fatalf("expected 201 Created, got %d, body: %s", resp.Code, resp.Body.String())
// 	}
// }
