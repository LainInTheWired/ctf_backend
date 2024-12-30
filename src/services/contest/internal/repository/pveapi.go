package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/LainInTheWired/ctf_backend/contest/model"
	"github.com/cockroachdb/errors"
)

type PVEAPIRepository interface {
	GetIPByVMID(vmid int) (*model.ResponseIPs, error)
	GetClusterResource() ([]model.ClusterResources, error)
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

func (r *pveapiRepository) GetIPByVMID(vmid int) (*model.ResponseIPs, error) {
	endpoint := fmt.Sprintf("%s/vm/%d/ips", r.URL, vmid)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't create http request: %w")
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "fail http request:")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading clone response body: %w")
	}

	var ifs model.ResponseIPs
	if err := json.Unmarshal(body, &ifs); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body: %w")
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return nil, errors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}
	return &ifs, nil
}

func (r *pveapiRepository) GetClusterResource() ([]model.ClusterResources, error) {
	endpoint := fmt.Sprintf("%s/cluster", r.URL)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't create http request: %w")
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "fail http request:")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading clone response body: %w")
	}

	var cluster []model.ClusterResources
	if err := json.Unmarshal(body, &cluster); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body: %w")
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return nil, errors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}
	return cluster, nil
}
