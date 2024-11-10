package repository

import "net/http"

type PVEAPIRepository interface {
}

type pveapiRepository struct {
	HTTPClient *http.Client
}

func NewPVEAPIRepository(hc *http.Client) PVEAPIRepository {
	return &pveapiRepository{
		HTTPClient: hc,
	}
}
