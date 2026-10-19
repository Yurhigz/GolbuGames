package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func TestGetUserFriendsInvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/abc/friends", nil)
	req.SetPathValue("id", "abc")

	rec := httptest.NewRecorder()

	var resp MessageResponse
	GetUserFriends(rec, req)

	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected JSON response, got %q", contentType)
	}

	if resp.Message != "invalid user id" {
		t.Fatalf("expected message %q, got %q", "invalid user id", resp.Message)
	}
}

func TestRemoveFriendInvalidIDS(t *testing.T) {
	test_cases := []struct {
		name            string
		userID          string
		friendID        string
		expectedCode    int
		expectedMessage string
	}{
		{
			name:            "Invalid User ID",
			userID:          "invalid",
			friendID:        "1",
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid user id",
		},
		{
			name:            "Invalid Friend ID",
			userID:          "1",
			friendID:        "invalid",
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid friend id",
		},
	}

	for _, tt := range test_cases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.userID+tt.friendID, nil)
			req.SetPathValue("id", tt.userID)
			req.SetPathValue("f_id", tt.friendID)

			rec := httptest.NewRecorder()

			RemoveFriend(rec, req)

			var resp MessageResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if rec.Code != tt.expectedCode {
				t.Fatalf("expected %d, got %d", tt.expectedCode, rec.Code)
			}

			if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
				t.Fatalf("expected JSON response, got %q", contentType)
			}

			if resp.Message != tt.expectedMessage {
				t.Fatalf("expected message %q, got %q", tt.expectedMessage, resp.Message)
			}
		})

	}
}

func TestAddFriendInvalidRequest(t *testing.T) {}
