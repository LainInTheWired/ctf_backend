package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/LainInTheWired/ctf-backend/pveapi/model"
	"github.com/cockroachdb/errors"
	"golang.org/x/xerrors"
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
	CloneVM(name string, id int, cnode string, cloneid int, tnode string) error
	GetVM(string) (*model.VMConfig, error)
	EditVM(model.VMEdit) error
	GetNodeList() ([]model.NodeList, error)
	GetVMList(nodes *model.NodeList) ([]model.VMList, error)
	DeleteVM(vmdelete *model.VMDelete) error
	CloudinitGenerator(fname string, host string, fqdn string, sshPwauth int, users []model.User) error
	TransferFileViaSCP(fname string) error
	NextVMID() (string, error)
	GetClusterResourcesList() ([]model.ClusterResources, error)
	ResizeDisk(node string, disk string, size int, vmid int) error
	Boot(node string, vmid int) error
	Shutdown(node string, vmid int) error
	Template(node string, vmid int) error
	DeleteFile(fname string) error
	GetNetIntFormQumeAgent(node string, vmid int) ([]model.NetworkIntQumeAgent, error)
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
	if vmedit.Memory != 0 {
		formData.Set("memory", strconv.Itoa(vmedit.Memory))
	}
	if vmedit.Cores != 0 {
		formData.Set("cores", strconv.Itoa(vmedit.Cores))

	}

	for i, v := range vmedit.Ipconfig {
		formData.Set(fmt.Sprintf("ipconfig%d", i), v)
	}
	for i, v := range vmedit.Scsi {
		formData.Set(fmt.Sprintf("scsi%d", i), v)
	}

	if vmedit.Cicustom != "" {
		formData.Set("cicustom", fmt.Sprintf("user=cephfs:snippets/%s", vmedit.Cicustom))
		fmt.Printf("user=cephfs:snippets/%s", vmedit.Cicustom)
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
		return errors.Newf("API Error: status code %d, response: %s, body: %s", resp.StatusCode, pveresp.Errors["errors"])
	}

	// クローン作成のUPIDを表示
	log.Printf("VM クローンの作成が開始されました。UPID: %s\n", pveresp.Data)

	return nil
}

func (r *pveRepository) DeleteStorageContest(node, storage, content string) error {
	// フォームデータの作成
	endpoint := fmt.Sprintf("%s/nodes/%s/storage/%s/content/%s", r.pveConf.APIURL, node, storage, content)
	// フォームデータの作成
	formData := url.Values{}

	// 新しいPOSTリクエストの作成
	req, err := http.NewRequest("DELETE", endpoint, strings.NewReader(formData.Encode()))
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
		return errors.Newf("API Error: status code %d, response: %s, body: %s", resp.StatusCode, pveresp.Errors["errors"])
	}

	// クローン作成のUPIDを表示
	log.Printf("VM クローンの作成が開始されました。UPID: %s\n", pveresp.Data)

	return nil
}

// この関数は使用していない
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

func (r *pveRepository) CloneVM(name string, id int, cnode string, cloneid int, tnode string) error {
	// フォームデータの作成
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu/%d/clone", r.pveConf.APIURL, cnode, cloneid)
	fmt.Printf("%s/nodes/%s/qemu/%b/clone\n", r.pveConf.APIURL, cnode, cloneid)
	// フォームデータの作成
	formData := url.Values{}
	formData.Set("name", name)
	formData.Set("newid", strconv.Itoa(id))
	formData.Set("target", tnode)
	formData.Set("full", "0")

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

	fmt.Println("fiewajfioejwaofjipoeawjfpojewaio;fj")
	fmt.Printf("%+v", pveresp.Data)
	return pveresp.Data, nil
}

func (r *pveRepository) ResizeDisk(node string, disk string, size int, vmid int) error {
	// フォームデータの作成
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu/%d/resize", r.pveConf.APIURL, node, vmid)
	fmt.Printf("%s/nodes/%s/qemu/%b/clone\n", r.pveConf.APIURL, node, vmid)
	// フォームデータの作成
	formData := url.Values{}
	formData.Set("disk", disk)
	formData.Set("size", fmt.Sprintf("%dG", size))

	// 新しいPOSTリクエストの作成
	req, err := http.NewRequest("PUT", endpoint, bytes.NewBufferString(formData.Encode()))
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
	log.Printf("resize  vm disk %s\n", pveresp.Data)

	return nil
}

func (r *pveRepository) Boot(node string, vmid int) error {
	// フォームデータの作成
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu/%d/status/start", r.pveConf.APIURL, node, vmid)
	// フォームデータの作成

	// 新しいPOSTリクエストの作成
	req, err := http.NewRequest("POST", endpoint, nil)
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
	log.Printf("start vm%s\n", pveresp.Data)

	return nil
}

func (r *pveRepository) Shutdown(node string, vmid int) error {
	// フォームデータの作成
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu/%d/status/shutdown", r.pveConf.APIURL, node, vmid)
	// フォームデータの作成

	// 新しいPOSTリクエストの作成
	req, err := http.NewRequest("POST", endpoint, nil)
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
	log.Printf("shutdown vm%s\n", pveresp.Data)

	return nil

}

func (r *pveRepository) Template(node string, vmid int) error {
	// フォームデータの作成
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu/%d/template", r.pveConf.APIURL, node, vmid)
	// フォームデータの作成

	// 新しいPOSTリクエストの作成
	req, err := http.NewRequest("POST", endpoint, nil)
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
	log.Printf("template vm%s\n", pveresp.Data)

	return nil
}

func (r *pveRepository) GetNetIntFormQumeAgent(node string, vmid int) ([]model.NetworkIntQumeAgent, error) {
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu/%d/agent/network-get-interfaces", r.pveConf.APIURL, node, vmid)
	// formData := url.Values{}
	fmt.Println(endpoint)

	req, err := http.NewRequest("GET", endpoint, nil)
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

	var getVMconfig model.ResponsePVE[struct {
		Result []model.NetworkIntQumeAgent `json:"result"`
	}]

	err = json.Unmarshal([]byte(body), &getVMconfig)
	if err != nil {
		return nil, xerrors.Errorf("json Unmarshal error: %w", err)
	}
	return getVMconfig.Data.Result, nil
}
