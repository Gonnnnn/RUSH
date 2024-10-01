package season

type repo struct {
	currentSeason Season
}

// TODO(#165): Implement the repo to handle the season data.
func NewRepo(currentSeason Season) *repo {
	return &repo{currentSeason: currentSeason}
}

func (r *repo) GetCurrentSeason() Season {
	return r.currentSeason
}
