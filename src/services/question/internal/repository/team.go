package repository

import (
	"net/http"
)

type TeamRepository interface {
}

type teamRepository struct {
	HTTPClient *http.Client
	URL        string
}

func NewTeamRepository(h *http.Client, url string) TeamRepository {
	return &teamRepository{
		HTTPClient: h,
		URL:        url,
	}
}
