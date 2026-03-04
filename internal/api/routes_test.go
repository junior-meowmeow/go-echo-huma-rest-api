package api_test

// import (
// 	"strings"
// 	"testing"

// 	"github.com/danielgtaylor/huma/v2/humatest"

// 	. "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/api"
// )

// func TestGetGreeting(t *testing.T) {
// 	_, testAPI := humatest.New(t)
// 	RegisterRoutes(testAPI)

// 	resp := testAPI.Get("/greeting/world")
// 	if resp.Code != 200 {
// 		t.Fatalf("expected status 200 got %d", resp.Code)
// 	}
// 	if !strings.Contains(resp.Body.String(), "Hello, world!") {
// 		t.Fatalf("expected greeting in body got: %s", resp.Body.String())
// 	}
// }
