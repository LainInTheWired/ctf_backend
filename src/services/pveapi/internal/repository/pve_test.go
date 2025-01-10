package repository

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/LainInTheWired/ctf-backend/pveapi/model"
	"github.com/joho/godotenv"
	echoLog "github.com/labstack/gommon/log" // エイリアスを付ける

	"golang.org/x/xerrors"
)

// 認証ヘッダーを自動的に付与するカスタムクライアントを作成
type MiddlewareTransport struct {
	Transport http.RoundTripper
	Token     string
}

type httpClientRequestLog struct {
	method  string
	url     string
	headers map[string][]string
	body    string
}
type httpClinetResponseLog struct {
	status string
	header http.Header
	body   string
}

// リクエストのたびに Authorization ヘッダーを自動で追加
func (a *MiddlewareTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s", a.Token))
	var reqbody []byte
	if req.Body != nil {
		reqbody, err := io.ReadAll(req.Body)
		if err != nil {
			xerrors.Errorf("can't read reqest body : %w", err)
		}
		req.Body = io.NopCloser(bytes.NewBuffer(reqbody))

	}

	reqlog := httpClientRequestLog{
		method:  req.Method,
		url:     req.URL.String(),
		headers: req.Header,
		body:    string(reqbody),
	}

	echoLog.Debug(reqlog)
	res, err := a.Transport.RoundTrip(req)
	if err != nil {
		xerrors.Errorf("http clinet send error : %w", err)
	}

	if res != nil && res.Body != nil {

		resbody, err := io.ReadAll(res.Body)
		if err != nil {
			xerrors.Errorf("can't read response body : %w", err)
		}
		res.Body = io.NopCloser(bytes.NewBuffer(resbody))
		reslog := httpClinetResponseLog{
			status: res.Status,
			header: res.Header,
			body:   string(resbody),
		}
		echoLog.Debug(reslog)
	}
	return res, err
}
func TestGetNodeList(t *testing.T) {
	// .envファイルを読み込む
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// proxmoxapi初期化処理
	config := &model.PVEConfig{
		APIURL:        os.Getenv("PROXMOX_API_URL"),
		Authorization: os.Getenv("PROXMOX_API_TOKEN"),
	}
	// httpclinet auth middleware
	// カスタムトランスポートを作成
	// カスタム Transport を設定（InsecureSkipVerify は本番環境では false に）
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 本番では false にする
	}

	// カスタム AuthTransport を作成
	authTransport := &MiddlewareTransport{
		Transport: tr,
		Token:     config.Authorization,
	}
	// カスタム HTTP クライアントの作成
	client := &http.Client{
		Transport: authTransport,
		Timeout:   60 * time.Second,
	}

	// 設定のバリデーション
	if config.APIURL == "" || config.Authorization == "" {
		log.Fatal("必要な環境変数が設定されていません。PROXMOX_API_URL および PROXMOX_AUTHORIZATION を設定してください。")
	}

	r := NewPVERepository(config, client)
	a, err := r.GetNodeList()
	if err != nil {
		log.Fatal(err)
	}
	t.Log(a)
}
func TestEditVMACL(t *testing.T) {
	// .envファイルを読み込む
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// proxmoxapi初期化処理
	config := &model.PVEConfig{
		APIURL:        os.Getenv("PROXMOX_API_URL"),
		Authorization: os.Getenv("PROXMOX_API_TOKEN"),
	}
	// httpclinet auth middleware
	// カスタムトランスポートを作成
	// カスタム Transport を設定（InsecureSkipVerify は本番環境では false に）
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 本番では false にする
	}

	// カスタム AuthTransport を作成
	authTransport := &MiddlewareTransport{
		Transport: tr,
		Token:     config.Authorization,
	}
	// カスタム HTTP クライアントの作成
	client := &http.Client{
		Transport: authTransport,
		Timeout:   60 * time.Second,
	}

	// 設定のバリデーション
	if config.APIURL == "" || config.Authorization == "" {
		log.Fatal("必要な環境変数が設定されていません。PROXMOX_API_URL および PROXMOX_AUTHORIZATION を設定してください。")
	}

	r := NewPVERepository(config, client)
	err = r.EditVMACL(137)
	if err != nil {
		log.Fatal(err)
	}

}
