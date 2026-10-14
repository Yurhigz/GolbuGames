package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserHandlersRejectInvalidRequests(t *testing.T) {
	testCases := []struct {
		name    string
		handler http.HandlerFunc
		method  string
		target  string
		body    string
		userID  string
	}{
		{name: "create user malformed JSON", handler: CreateUser, method: http.MethodPost, target: "/users", body: `{"username":`},
		{name: "create user missing required fields", handler: CreateUser, method: http.MethodPost, target: "/users", body: `{}`},
		{name: "delete user nonnumeric ID", handler: DeleteUser, method: http.MethodDelete, target: "/users/abc", userID: "abc"},
		{name: "delete user nonpositive ID", handler: DeleteUser, method: http.MethodDelete, target: "/users/0", userID: "0"},
		{name: "get user nonnumeric ID", handler: GetUser, method: http.MethodGet, target: "/users/abc", userID: "abc"},
		{name: "get user ID missing username", handler: GetUserId, method: http.MethodGet, target: "/users/id?accountname=account"},
		{name: "get user ID missing account name", handler: GetUserId, method: http.MethodGet, target: "/users/id?username=user"},
		{name: "update password malformed JSON", handler: UpdateUserPassword, method: http.MethodPut, target: "/users/password", body: `{"id":`},
		{name: "update password missing new password", handler: UpdateUserPassword, method: http.MethodPut, target: "/users/password", body: `{"id":1}`},
		{name: "signin malformed JSON", handler: UserSignin, method: http.MethodPost, target: "/signin", body: `{"username":`},
		{name: "signin missing required fields", handler: UserSignin, method: http.MethodPost, target: "/signin", body: `{}`},
		{name: "login malformed JSON", handler: UserLogin, method: http.MethodPost, target: "/login", body: `{"username":`},
		{name: "login missing required fields", handler: UserLogin, method: http.MethodPost, target: "/login", body: `{}`},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.body))
			if tt.userID != "" {
				req.SetPathValue("id", tt.userID)
			}
			rec := httptest.NewRecorder()

			tt.handler(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d; body: %s", http.StatusBadRequest, rec.Code, rec.Body.String())
			}
		})
	}
}
