package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserSuite struct {
	IntegrationTestSuite
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (s *UserSuite) SetupTest() {
	cleanCollection(s.T(), s.MongoDB.Collection("users"))
}

type registerUserResponse struct {
	ID string `json:"id"`
}

type loginUserResponse struct {
	Token string `json:"token"`
}

type registerUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`
}

type loginUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *UserSuite) registerUser(body registerUserBody) *httptest.ResponseRecorder {
	s.T().Helper()

	payload, err := json.Marshal(body)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *UserSuite) loginUser(body loginUserBody) *httptest.ResponseRecorder {
	s.T().Helper()

	payload, err := json.Marshal(body)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *UserSuite) mustRegisterUser(body registerUserBody) string {
	s.T().Helper()

	w := s.registerUser(body)
	if w.Code != http.StatusCreated {
		s.T().Logf("mustRegisterUser failed – status: %d, body: %s", w.Code, w.Body.String())
	}
	s.Require().Equal(http.StatusCreated, w.Code)

	var resp registerUserResponse
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp), "failed to decode RegisterUser response")
	s.Require().NotEmpty(resp.ID, "RegisterUser response must include an ID")

	return resp.ID
}

func (s *UserSuite) mustLoginUser(body loginUserBody) string {
	s.T().Helper()

	w := s.loginUser(body)
	if w.Code != http.StatusOK {
		s.T().Logf("mustLoginUser failed – status: %d, body: %s", w.Code, w.Body.String())
	}
	s.Require().Equal(http.StatusOK, w.Code)

	var resp loginUserResponse
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp), "failed to decode LoginUser response")
	s.Require().NotEmpty(resp.Token, "LoginUser response must include a token")

	return resp.Token
}

func (s *UserSuite) decodeRegisterUserResponse(w *httptest.ResponseRecorder) registerUserResponse {
	s.T().Helper()

	var resp registerUserResponse
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp), "failed to decode POST /v1/users/register response")
	return resp
}

func (s *UserSuite) decodeLoginUserResponse(w *httptest.ResponseRecorder) loginUserResponse {
	s.T().Helper()

	var resp loginUserResponse
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp), "failed to decode POST /v1/users/login response")
	return resp
}

func (s *UserSuite) dbUserCount() int64 {
	s.T().Helper()

	count, err := s.MongoDB.Collection("users").CountDocuments(context.Background(), bson.D{})
	s.Require().NoError(err)
	return count
}

func validUser() registerUserBody {
	return registerUserBody{
		Username: "testuser",
		Password: "securepass123",
		Role:     "user",
	}
}

func (s *UserSuite) TestRegisterUser_ReturnsHTTP201() {
	w := s.registerUser(validUser())

	s.Equal(http.StatusCreated, w.Code)
}

func (s *UserSuite) TestRegisterUser_ReturnsCreatedID() {
	w := s.registerUser(validUser())
	s.Require().Equal(http.StatusCreated, w.Code)

	resp := s.decodeRegisterUserResponse(w)
	s.NotEmpty(resp.ID, "response body must include the created user ID")
	s.Regexp(`^[a-fA-F0-9]{24}$`, resp.ID, "ID must be a valid BSON ObjectID")
}

func (s *UserSuite) TestRegisterUser_PersistsToMongoDB() {
	body := validUser()
	s.mustRegisterUser(body)

	s.Equal(int64(1), s.dbUserCount())

	var doc bson.M
	err := s.MongoDB.Collection("users").FindOne(context.Background(), bson.D{}).Decode(&doc)
	s.Require().NoError(err)

	s.Equal(body.Username, doc["username"])
	s.NotEmpty(doc["_id"])
	s.NotEmpty(doc["createdAt"])
}

func (s *UserSuite) TestRegisterUser_PasswordIsHashed() {
	body := validUser()
	s.mustRegisterUser(body)

	var doc bson.M
	err := s.MongoDB.Collection("users").FindOne(context.Background(), bson.D{}).Decode(&doc)
	s.Require().NoError(err)

	storedPassword, ok := doc["password"].(string)
	s.Require().True(ok, "password field should be a string")
	s.NotEqual(body.Password, storedPassword, "password must be stored as a hash, not plain text")
	s.NotEmpty(storedPassword)
}

func (s *UserSuite) TestRegisterUser_WithRoleAdmin() {
	body := registerUserBody{
		Username: "adminuser",
		Password: "adminpass123",
		Role:     "admin",
	}

	s.mustRegisterUser(body)

	var doc bson.M
	err := s.MongoDB.Collection("users").FindOne(context.Background(), bson.D{}).Decode(&doc)
	s.Require().NoError(err)
	s.Equal("admin", doc["role"])
}

func (s *UserSuite) TestRegisterUser_WithoutRole_UsesDefault() {
	body := registerUserBody{
		Username: "noroleuser",
		Password: "somepass123",
		// Role omitted
	}

	w := s.registerUser(body)
	s.Equal(http.StatusCreated, w.Code)
	s.Equal(int64(1), s.dbUserCount())
}

func (s *UserSuite) TestRegisterUser_DuplicateUsername_Returns409() {
	body := validUser()
	s.mustRegisterUser(body)

	// Register same username again
	w := s.registerUser(body)

	s.Equal(http.StatusConflict, w.Code)
	s.Equal(int64(1), s.dbUserCount(), "duplicate registration must not create a second user")
}

func (s *UserSuite) TestRegisterUser_ValidationErrors() {
	cases := []struct {
		name       string
		body       registerUserBody
		wantStatus int
	}{
		{
			name:       "Username too short (2 chars)",
			body:       registerUserBody{Username: "ab", Password: "validpass123", Role: "user"},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "Username too long (51 chars)",
			body:       registerUserBody{Username: strings.Repeat("a", 51), Password: "validpass123", Role: "user"},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "Password too short (7 chars)",
			body:       registerUserBody{Username: "validuser", Password: "short12", Role: "user"},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "Password too long (101 chars)",
			body:       registerUserBody{Username: "validuser", Password: strings.Repeat("a", 101), Role: "user"},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "Invalid role",
			body:       registerUserBody{Username: "validuser", Password: "validpass123", Role: "superuser"},
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			cleanCollection(s.T(), s.MongoDB.Collection("users"))

			w := s.registerUser(tc.body)

			s.Equal(tc.wantStatus, w.Code, "case: %q", tc.name)
			s.Equal(int64(0), s.dbUserCount(), "no user should be persisted for invalid input in case %q", tc.name)
		})
	}
}

func (s *UserSuite) TestRegisterUser_BoundaryValues() {
	exactly3Chars := strings.Repeat("a", 3)
	exactly50Chars := strings.Repeat("a", 50)
	exactly8Chars := strings.Repeat("a", 8)
	exactly100Chars := strings.Repeat("a", 100)

	cases := []struct {
		name string
		body registerUserBody
	}{
		{
			name: "Username at minLength (3 chars)",
			body: registerUserBody{Username: exactly3Chars, Password: "validpass123", Role: "user"},
		},
		{
			name: "Username at maxLength (50 chars)",
			body: registerUserBody{Username: exactly50Chars, Password: "validpass123", Role: "user"},
		},
		{
			name: "Password at minLength (8 chars)",
			body: registerUserBody{Username: "validuser1", Password: exactly8Chars, Role: "user"},
		},
		{
			name: "Password at maxLength (100 chars)",
			body: registerUserBody{Username: "validuser2", Password: exactly100Chars, Role: "user"},
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			cleanCollection(s.T(), s.MongoDB.Collection("users"))

			w := s.registerUser(tc.body)

			s.Equal(http.StatusCreated, w.Code, "boundary value should be accepted for case %q", tc.name)
			s.Equal(int64(1), s.dbUserCount(), "one document should be persisted for case %q", tc.name)
		})
	}
}

func (s *UserSuite) TestRegisterUser_MultipleUsers_AllPersisted() {
	users := []registerUserBody{
		{Username: "user_alpha", Password: "passalpha1", Role: "user"},
		{Username: "user_beta", Password: "passbeta123", Role: "user"},
		{Username: "user_gamma", Password: "passgamma12", Role: "admin"},
	}

	ids := make(map[string]bool)
	for _, u := range users {
		ids[s.mustRegisterUser(u)] = true
	}

	s.Equal(int64(len(users)), s.dbUserCount())
	s.Len(ids, len(users), "every registered user must have a unique ID")
}

func (s *UserSuite) TestLoginUser_ReturnsHTTP200() {
	reg := validUser()
	s.mustRegisterUser(reg)

	w := s.loginUser(loginUserBody{Username: reg.Username, Password: reg.Password})

	s.Equal(http.StatusOK, w.Code)
}

func (s *UserSuite) TestLoginUser_ReturnsToken() {
	reg := validUser()
	s.mustRegisterUser(reg)

	w := s.loginUser(loginUserBody{Username: reg.Username, Password: reg.Password})
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeLoginUserResponse(w)
	s.NotEmpty(resp.Token, "response body must include a JWT token")
	s.Regexp(`^[\w-]+\.[\w-]+\.[\w-]+$`, resp.Token, "token must be a valid JWT format")
}

func (s *UserSuite) TestLoginUser_WrongPassword_Returns401() {
	reg := validUser()
	s.mustRegisterUser(reg)

	w := s.loginUser(loginUserBody{Username: reg.Username, Password: "wrongpassword"})

	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *UserSuite) TestLoginUser_NonExistentUser_Returns401() {
	w := s.loginUser(loginUserBody{Username: "ghost_user", Password: "somepassword"})

	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *UserSuite) TestLoginUser_ValidationErrors() {
	cases := []struct {
		name       string
		body       loginUserBody
		wantStatus int
	}{
		{
			name:       "Username too short",
			body:       loginUserBody{Username: "ab", Password: "validpass123"},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "Password too short",
			body:       loginUserBody{Username: "validuser", Password: "short"},
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			w := s.loginUser(tc.body)
			s.Equal(tc.wantStatus, w.Code, "case: %q", tc.name)
		})
	}
}

func (s *UserSuite) TestRegisterThenLogin_RoundTrip() {
	reg := validUser()
	s.mustRegisterUser(reg)

	token := s.mustLoginUser(loginUserBody{Username: reg.Username, Password: reg.Password})
	s.NotEmpty(token)
}

func (s *UserSuite) TestRegisterThenLogin_TokenIsUsableOnProtectedRoute() {
	reg := validUser()
	s.mustRegisterUser(reg)
	token := s.mustLoginUser(loginUserBody{Username: reg.Username, Password: reg.Password})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/books", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	s.Router.ServeHTTP(w, req)

	s.NotEqual(http.StatusUnauthorized, w.Code, "a valid token should be accepted by protected routes")
	s.NotEqual(http.StatusForbidden, w.Code, "a valid token should not be rejected as forbidden")
}

func (s *UserSuite) TestProtectedRoute_WithoutToken_Returns401() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test-auth/test", nil)
	s.Router.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *UserSuite) TestProtectedRoute_WithInvalidToken_Returns401() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test-auth/test", nil)
	req.Header.Set("Authorization", "Bearer this.is.not.a.valid.token")
	s.Router.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *UserSuite) TestLoginTwice_ReturnsDifferentTokens() {
	reg := validUser()
	s.mustRegisterUser(reg)
	creds := loginUserBody{Username: reg.Username, Password: reg.Password}

	token1 := s.mustLoginUser(creds)
	time.Sleep(time.Second)
	token2 := s.mustLoginUser(creds)

	// JWTs issued at different times will have different iat claims.
	s.NotEqual(token1, token2, "two login calls should produce distinct tokens")
}
