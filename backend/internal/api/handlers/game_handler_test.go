package handlers

import (
	"bytes"
	"encoding/json"
	"golbugames/pkg/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubmitSoloGameWrong(t *testing.T) {
	test_cases := []struct {
		name           string
		userid         string
		gameMode       string
		opponentID     string
		completionTime int
		results        *int
		difficulty     string
		expectedCode   int
	}{
		{
			name:           "Invalid Game Mode",
			userid:         "1",
			gameMode:       "1v1",
			opponentID:     "",
			completionTime: 100,
			results:        nil,
			difficulty:     "easy",
			expectedCode:   http.StatusBadRequest,
		},
		{
			name:           "Missing Opponent ID",
			gameMode:       "solo",
			opponentID:     "2",
			completionTime: 100,
			results:        nil,
			difficulty:     "easy",
			expectedCode:   http.StatusBadRequest,
		},
		{
			name:           "Missing Completion Time",
			gameMode:       "solo",
			opponentID:     "",
			completionTime: 0,
			results:        nil,
			difficulty:     "easy",
			expectedCode:   http.StatusBadRequest,
		},
		{
			name:           "Invalid Difficulty",
			gameMode:       "solo",
			opponentID:     "",
			completionTime: 100,
			results:        nil,
			difficulty:     "",
			expectedCode:   http.StatusBadRequest,
		},
		{
			name:           "Wring Results for Solo Game",
			userid:         "1",
			gameMode:       "solo",
			opponentID:     "",
			completionTime: 100,
			results:        &[]int{1}[0],
			difficulty:     "easy",
			expectedCode:   http.StatusBadRequest,
		},
	}

	for _, tt := range test_cases {
		t.Run(tt.name, func(t *testing.T) {
			games := types.Game{UserID: tt.userid, GameMode: tt.gameMode, OpponentID: tt.opponentID, Completion_time: tt.completionTime, Results: tt.results, Difficulty: tt.difficulty}
			body, err := json.Marshal(games)
			if err != nil {
				t.Fatalf("failed to marshal game: %v", err)
			}
			req := httptest.NewRequest(http.MethodPost, "/games/solo", bytes.NewReader(body))
			rec := httptest.NewRecorder()
			SubmitSoloGame(rec, req)

			if rec.Code != tt.expectedCode {
				t.Fatalf("expected %d, got %d", tt.expectedCode, rec.Code)
			}
		})
	}

}

func TestSubmitMultiGameWrongGameMode(t *testing.T) {
	test_cases := []struct {
		name           string
		userid         string
		gameMode       string
		opponentID     string
		completionTime int
		results        *int
		difficulty     string
		expectedCode   int
	}{
		{
			name:           "Invalid Game Mode",
			userid:         "1",
			gameMode:       "solo",
			opponentID:     "2",
			completionTime: 100,
			results:        &[]int{1}[0],
			difficulty:     "easy",
			expectedCode:   http.StatusBadRequest,
		},
		{
			name:           "Missing Opponent ID",
			userid:         "1",
			gameMode:       "1v1",
			opponentID:     "",
			completionTime: 100,
			results:        &[]int{1}[0],
			difficulty:     "easy",
			expectedCode:   http.StatusBadRequest,
		},
		{
			name:           "Missing Completion Time",
			userid:         "1",
			gameMode:       "1v1",
			opponentID:     "2",
			completionTime: 0,
			results:        &[]int{1}[0],
			difficulty:     "easy",
			expectedCode:   http.StatusBadRequest,
		},
		{
			name:           "Missing Results",
			userid:         "1",
			gameMode:       "1v1",
			opponentID:     "2",
			completionTime: 100,
			results:        nil,
			difficulty:     "easy",
			expectedCode:   http.StatusBadRequest,
		},
		{
			name:           "Invalid Difficulty",
			userid:         "1",
			gameMode:       "1v1",
			opponentID:     "2",
			completionTime: 100,
			results:        &[]int{1}[0],
			difficulty:     "",
			expectedCode:   http.StatusBadRequest,
		},
	}

	for _, tt := range test_cases {
		t.Run(tt.name, func(t *testing.T) {
			games := types.Game{UserID: tt.userid, GameMode: tt.gameMode, OpponentID: tt.opponentID, Completion_time: tt.completionTime, Results: tt.results, Difficulty: tt.difficulty}
			body, err := json.Marshal(games)
			if err != nil {
				t.Fatalf("failed to marshal game: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/games/multi", bytes.NewReader(body))
			rec := httptest.NewRecorder()
			SubmitMultiGame(rec, req)

			if rec.Code != tt.expectedCode {
				t.Fatalf("expected %d, got %d", tt.expectedCode, rec.Code)
			}

		})

	}
}
