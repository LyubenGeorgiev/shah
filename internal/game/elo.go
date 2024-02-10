package game

import (
	"errors"
	"math"
	"sort"
)

// User is the user model (defined in models/user.go)
type User struct {
	Username    string  `json:"username"`
	Email       string  `json:"email" gorm:"type:varchar(100);unique_index"`
	Password    string  `json:"password"`
	Rating      float64 `json:"rating"`
	GamesPlayed int     `json:"gamesPlayed"`
	GamesWon    int     `json:"gamesWon"`
}

func UpdateRating(winner *User, loser *User) {

	KFactor := 32.0
	WinProbability := 1.0 / (1.0 + math.Pow(10, (loser.Rating-winner.Rating)/400.0))

	// Calculate the change in ratings
	winnerDelta  := (KFactor * (1 - WinProbability))
	loserDelta :=  (KFactor * (0 - WinProbability))

	// Update ratings and game statistics
	winner.Rating += winnerDelta
	loser.Rating += loserDelta

	winner.GamesPlayed++
	winner.GamesWon++

	loser.GamesPlayed++
}

type MatchmakingParams struct {
	MaxRatingDifference int
}

// Matchmake finds an opponent for a given player based on their ratings.
func Matchmake(player *User, pool []*User, params MatchmakingParams) (*User, error) {
	// Filter out the player from the pool
	var opponents []*User
	for _, p := range pool {
		if p.Username != player.Username {
			opponents = append(opponents, p)
		}
	}

	// Sort opponents by rating difference
	sort.Slice(opponents, func(i, j int) bool {
		return math.Abs(player.Rating-opponents[i].Rating) < math.Abs(player.Rating-opponents[j].Rating)
	})

	// Find a suitable opponent within the rating difference limit
	for _, opponent := range opponents {
		ratingDifference := math.Abs(player.Rating - opponent.Rating)
		if int(ratingDifference) <= params.MaxRatingDifference {
			return opponent, nil
		}
	}

	return nil, errors.New("no suitable opponent found")
}
