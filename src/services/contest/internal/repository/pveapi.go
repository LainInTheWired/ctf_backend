package repository

import "net/http"

type PVEAPIRepository interface {
}

type pveapiRepository struct {
	HTTPClient *http.Client
	URL        string
}

func NewPVEAPIRepository(hc *http.Client, url string) PVEAPIRepository {
	return &pveapiRepository{
		HTTPClient: hc,
		URL:        url,
	}
}
