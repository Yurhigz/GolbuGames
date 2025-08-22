package sudoku

//  intégrer les types spécifiques au jeu sudoku plutôt que de tout stocket dans pkg

// type Leaderboard struct {
// 	UserID   int `json:"userid"`
// 	EloScore int `json:"elo_score"`
// 	Rank     int `json:"rank"`
// }

type Leaderboard struct {
    Username string `json:"username"`
    EloScore int    `json:"elo_score"`
    Rank     int    `json:"rank"`
}


type GameScore struct {
	GameID         int `json:"gameid"`
	UserID         int `json:"userid"`
	Opponent_id    int `json:"opponent_id"`
	GameMode       int `json:"game_mode"`
	Result         int `json:"results"`
	CompletionTime int `json:"completion_time"`
}
