package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/LainInTheWired/ctf-backend/pveapi/model"
)

var app *App

func TestMain(m *testing.M) {
	var err error

	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
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

	code := m.Run()

	// テスト終了後のクリーンアップ
	app.DB.Close()
	app.Redis.Close()

	os.Exit(code)
}
