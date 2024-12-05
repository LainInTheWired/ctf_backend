package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/LainInTheWired/ctf_backend/contest/model"
	"github.com/cockroachdb/errors"
	"golang.org/x/xerrors"
)

type QuestionRepository interface {
	GetListQuestionsByContest(cid int) ([]model.Question, error)
	CloneQuestion(conf model.QuesionRequest) error
	GetListQuestionsByQuestionID(qid int) (model.Question, error)
}

type questionRepository struct {
	HTTPClient *http.Client
	URL        string
}

func NewQuestionRepository(hc *http.Client, url string) QuestionRepository {
	return &questionRepository{
		HTTPClient: hc,
		URL:        url,
	}
}

func (r *questionRepository) GetListQuestionsByContest(cid int) ([]model.Question, error) {
	endpoint := fmt.Sprintf("%s/question/%d", r.URL, cid)
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
	var questions []model.Question
	if err := json.Unmarshal(body, &questions); err != nil {
		return nil, xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return nil, xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}

	// クローン作成のUPIDを表示
	log.Printf("VM クローンの作成が開始されました。UPID: %v\n", questions)
	return questions, nil
}

func (r *questionRepository) GetListQuestionsByQuestionID(qid int) (model.Question, error) {
	endpoint := fmt.Sprintf("%s/question?id=%d", r.URL, qid)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return model.Question{}, xerrors.Errorf("can't create http request: %w", err)
	}
	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// リクエストの送信
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return model.Question{}, xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	// json.Unmarshalでデコード
	var question model.Question
	if err := json.Unmarshal(body, &question); err != nil {
		return model.Question{}, xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return model.Question{}, xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}

	// クローン作成のUPIDを表示
	log.Printf("VM クローンの作成が開始されました。UPID: %v\n", question)

	return question, nil
}

func (r *questionRepository) CloneQuestion(conf model.QuesionRequest) error {
	// フォームデータの作成
	endpoint := fmt.Sprintf("%s/question/clone", r.URL)
	// フォームデータの作成

	jsend, err := json.Marshal(conf)
	if err != nil {
		return errors.Wrap(err, "can't change json")
	}

	// 新しいPOSTリクエストの作成
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsend))
	if err != nil {
		return xerrors.Errorf("can't create http request: %w", err)
	}

	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/json")

	// リクエストの送信
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスの読み取り
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalf("Error reading clone response body: %v", err)
	// }

	// // json.Unmarshalでデコード
	// var pveresp model.[string]
	// if err := json.Unmarshal(body, &pveresp); err != nil {
	// 	return xerrors.Errorf("can't unmarshal response body: %w", err)
	// }

	// エラーチェック
	if resp.StatusCode >= 400 {
		return xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}
	return nil
}
