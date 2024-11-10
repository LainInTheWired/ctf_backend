package repository

import "net/http"

type pveapiRepository struct {
	HTTPClient *http.Client
}
type PVEAPIRepository interface {
}

func NewPVEAPIRepository(h *http.Client) PVEAPIRepository {
	return &pveapiRepository{
		HTTPClient: h,
	}
}

func (r *pveapiRepository) CreateVM() {

}
