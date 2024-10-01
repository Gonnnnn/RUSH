package aggregation

import "time"

type Aggregation struct {
	Id string
	// The scores of the users. It is not ordered by the scores.
	UserScores []UserScorePair
	// The IDs of the sessions that have been used to get the aggregation.
	SessionIds []string
	// The time when the document was created.
	CreatedAt time.Time
}

type UserScorePair struct {
	UserId     string
	UserName   string
	Generation float64
	Score      int
}
