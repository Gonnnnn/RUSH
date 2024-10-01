package aggregation

import "time"

type Aggregation struct {
	Id string
	// The aggregated information for each user.
	UserInfos []UserInfo
	// The IDs of the sessions that have been used to get the aggregation.
	SessionIds []string
	// The time when the document was created.
	CreatedAt time.Time
}

type UserInfo struct {
	UserId     string
	UserName   string
	Generation float64
	Score      int
}
