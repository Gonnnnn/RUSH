package season

import "time"

type Season struct {
	// The unique identifier for the season. E.g. "1"
	Id string
	// The name of the season. E.g. "Rush 10기"
	Name string
	// The start date of the season.
	StartDate time.Time
	// The end date of the season.
	EndDate time.Time
	// The order of the season. It's used to determine the season of the user.
	// It is either a positive integer or a positive decimal number with 0.5 step.
	// E.g. 10, 10.5, 11, 11.5, 12
	Order float64
	// The attendance score that the members have to get
	// to extend their membership to the next season.
	RequiredAttendanceScore int
	// The time when the document was created.
	CreatedAt time.Time
}
