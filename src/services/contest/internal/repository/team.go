package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/LainInTheWired/ctf_backend/contest/model"
	"golang.org/x/xerrors"
)

type TeamRepository interface {
	ListTeamUsersByContest(cid int) ([]model.Team, error)
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

func (r *teamRepository) ListTeamUsersByContest(cid int) ([]model.Team, error) {
	endpoint := fmt.Sprintf("%s/teamusers?id=%d", r.URL, cid)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, xerrors.Errorf("can't create http request: %w", err)
	}
	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// リクエストの送信
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	// json.Unmarshalでデコード
	var teams []model.Team
	if err := json.Unmarshal(body, &teams); err != nil {
		return nil, xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return nil, xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}

	// クローン作成のUPIDを表示
	log.Printf("VM クローンの作成が開始されました。UPID: %s\n", teams)

	return teams, nil
}
