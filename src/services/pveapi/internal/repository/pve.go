package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/LainInTheWired/ctf-backend/pveapi/model"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v3"
)

type pveRepository struct {
	pveConf    *model.PVEConfig
	HTTPClient *http.Client
}

// ProxmoxConfig はProxmoxへの接続設定を保持します
type PVEConfig struct {
	APIURL        string
	Authorization string
}

type PVERepository interface {
	CloneVM(name string, id int, node string, cloneid int) error
	GetVM(string) (*model.VMConfig, error)
	EditVM(model.VMEdit) error
	GetNodeList() ([]model.NodeList, error)
	GetVMList(nodes *model.NodeList) ([]model.VMList, error)
	DeleteVM(vmdelete *model.VMDelete) error
	CloudinitGenerator(fname string, host string, users []model.User) error
	TransferFileViaSCP(fname string) error
	NextVMID() (string, error)
	GetClusterResourcesList() ([]model.ClusterResources, error)
}

func NewPVERepository(conf *model.PVEConfig, client *http.Client) PVERepository {
	return &pveRepository{
		pveConf:    conf,
		HTTPClient: client,
	}
}

func (r *pveRepository) GetNodeList() ([]model.NodeList, error) {
	endpoint := fmt.Sprintf("%s/nodes", r.pveConf.APIURL)

	formData := url.Values{}

	req, err := http.NewRequest("GET", endpoint, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, xerrors.Errorf("can't create http request: %w", err)
	}
	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
	var pveresp model.ResponsePVE[[]model.NodeList]
	if err := json.Unmarshal(body, &pveresp); err != nil {
		return nil, xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return nil, xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}

	return pveresp.Data, nil
}

func (r *pveRepository) GetVMList(nodes *model.NodeList) ([]model.VMList, error) {
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu", r.pveConf.APIURL, strings.Split(nodes.ID, "/")[1])
	formData := url.Values{}
	req, err := http.NewRequest("GET", endpoint, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, xerrors.Errorf("can't create http request: %w", err)
	}
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	var pveresp model.ResponsePVE[[]model.VMList]
	if err := json.Unmarshal(body, &pveresp); err != nil {
		return nil, xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return nil, xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}
	return pveresp.Data, nil
}

func (r *pveRepository) EditVM(vmedit model.VMEdit) error {
	// フォームデータの作成
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu/%d/config", r.pveConf.APIURL, vmedit.Node, vmedit.Vmid)
	// フォームデータの作成
	formData := url.Values{}
	formData.Set("cores", "2")
	for i, v := range vmedit.Ipconfig {
		formData.Set(fmt.Sprintf("ipconfig%d", i), v)
	}
	for i, v := range vmedit.Scsi {
		formData.Set(fmt.Sprintf("scsi%d", i), v)
	}
	if vmedit.Cicustom != "" {
		formData.Set("cicustom", fmt.Sprintf("user=cephfs:snippets/%s", vmedit.Cicustom))
	}

	fmt.Println("Request Body:", formData.Encode())

	// 新しいPOSTリクエストの作成
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return xerrors.Errorf("can't create http request: %w", err)
	}

	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// リクエストの送信
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	// json.Unmarshalでデコード
	var pveresp model.ResponsePVE[string]
	if err := json.Unmarshal(body, &pveresp); err != nil {
		return xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}

	// クローン作成のUPIDを表示
	log.Printf("VM クローンの作成が開始されました。UPID: %s\n", pveresp.Data)

	return nil
}

func (r *pveRepository) GetVM(id string) (*model.VMConfig, error) {
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu/%s/config", r.pveConf.APIURL, "pve02", id)
	formData := url.Values{}

	req, err := http.NewRequest("GET", endpoint, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, xerrors.Errorf("can't create http request: %w", err)
	}

	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return nil, xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}

	var getVMconfig model.ResponsePVE[model.VMConfig]

	err = json.Unmarshal([]byte(body), &getVMconfig)
	if err != nil {
		return nil, xerrors.Errorf("json Unmarshal error: %w", err)
	}
	fmt.Println(getVMconfig)

	return &getVMconfig.Data, nil
}

func (r *pveRepository) CloneVM(name string, id int, node string, cloneid int) error {

	// フォームデータの作成
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu/%d/clone", r.pveConf.APIURL, node, cloneid)
	fmt.Printf("%s/nodes/%s/qemu/%b/clone\n", r.pveConf.APIURL, node, cloneid)
	// フォームデータの作成
	formData := url.Values{}
	formData.Set("name", name)
	formData.Set("newid", strconv.Itoa(id))
	formData.Set("node", node)
	formData.Set("full", "1")

	// 新しいPOSTリクエストの作成
	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return xerrors.Errorf("can't create http request: %w", err)
	}

	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// リクエストの送信
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスの読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	// json.Unmarshalでデコード
	var pveresp model.ResponsePVE[string]
	if err := json.Unmarshal(body, &pveresp); err != nil {
		return xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}

	// クローン作成のUPIDを表示
	log.Printf("VM クローンの作成が開始されました。UPID: %s\n", pveresp.Data)

	return nil
}

func (r *pveRepository) DeleteVM(vmdelete *model.VMDelete) error {
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu/%d/", r.pveConf.APIURL, vmdelete.Node, vmdelete.Vmid)
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return xerrors.Errorf("can't create http request: %w", err)
	}
	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// リクエストの送信
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスの読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	// json.Unmarshalでデコード
	var pveresp model.ResponsePVE[string]
	if err := json.Unmarshal(body, &pveresp); err != nil {
		return xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}
	return nil
}

func (r *pveRepository) CloudinitGenerator(fname string, host string, users []model.User) error {
	config := model.CloudConfig{
		Hostname: host,
		Users:    users,
		Packages: []string{
			"git",
			"curl",
		},
	}
	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Printf("Error marshalling YAML: %v\n", err)
		return nil
	}
	yamlData = append([]byte("#cloud-config\n"), yamlData...)

	// ファイルに書き出し
	err = os.WriteFile(fname, yamlData, 0644)
	if err != nil {
		fmt.Printf("Error writing YAML to file: %v\n", err)
		return nil
	}
	return nil
}

func (r *pveRepository) TransferFileViaSCP(fname string) error {
	localPath := fname
	remoteUser := "root"
	// remoteHost := "10.0.10.30"
	remoteHost := os.Getenv("PROXMOX_API_URL")
	remotePath := "/mnt/pve/cephfs/snippets/"
	cmd := exec.Command("scp",
		"-i", "/ssh/id_ed25519",
		"-o", "StrictHostKeyChecking=no",
		localPath,
		fmt.Sprintf("%s@%s:%s", remoteUser, remoteHost, remotePath),
	)

	// コマンドの標準出力と標準エラーを取得
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// コマンドの実行
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("scp コマンドの実行に失敗しました: %w", err)
	}

	return nil

}

func (r *pveRepository) NextVMID() (string, error) {
	endpoint := fmt.Sprintf("%s/cluster/nextid", r.pveConf.APIURL)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", xerrors.Errorf("can't create http request: %w", err)
	}
	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// リクエストの送信
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return "", xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	// json.Unmarshalでデコード
	var pveresp model.ResponsePVE[string]
	if err := json.Unmarshal(body, &pveresp); err != nil {
		return "", xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return "", xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}

	// クローン作成のUPIDを表示
	log.Printf("VM クローンの作成が開始されました。UPID: %s\n", pveresp.Data)

	return pveresp.Data, nil
}

func (r *pveRepository) GetClusterResourcesList() ([]model.ClusterResources, error) {
	endpoint := fmt.Sprintf("%s/cluster/resources", r.pveConf.APIURL)
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
	var pveresp model.ResponsePVE[[]model.ClusterResources]
	if err := json.Unmarshal(body, &pveresp); err != nil {
		return nil, xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return nil, xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}

	// クローン作成のUPIDを表示
	log.Printf("VM クローンの作成が開始されました。UPID: %s\n", pveresp.Data)

	return pveresp.Data, nil
}
