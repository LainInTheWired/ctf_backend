package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

type pveService struct {
	pveConf    *PVEConfig
	HTTPClient *http.Client
}

type PVEService interface {
	CloneVM(string, int, string, int) error
	GetVM(string) (*VMConfig, error)
	EditVM(VMEdit) error
	GetNodeList() ([]NodeList, error)
	GetVMList(nodes *NodeList) ([]VMList, error)
	DeleteVM(vmdelete *VMDelete) error
}

// ProxmoxConfig はProxmoxへの接続設定を保持します
type PVEConfig struct {
	APIURL        string
	Authorization string
}

// VMConfig は作成するVMの設定を保持します
//	type VMConfig struct {
//		VMID           string `json:"vmid"`
//		Name           string `json:"name"`
//		Memory         string `json:"memory"` // MB単位
//		CPUs           string `json:"cores"`
//		Net0           string `json:"net0"`        // 例: "virtio=DE:AD:BE:EF:00:00,bridge=vmbr0"
//		Scsi0          string `json:"scsi0"`       // 例: "kingston_1tb:vm-200-disk-0,size=16G"
//		Boot           string `json:"boot"`        // 例: "c"
//		Ide2           string `json:"ide2"`        // 例: "local:iso/AlmaLinux-9.3-x86_64-boot.iso"
//		OSType         string `json:"ostype"`      // 例: "l26" (Linux 2.6/3.x/4.x)
//		SCSIController string `json:"scsihw"`      // 例: "virtio-scsi-single"
//		Description    string `json:"description"` // VMの説明（オプション）
//	}

type VMConfig struct {
	Vmgenid    string `json:"vmgenid"`
	Sockets    int    `json:"sockets"`
	Net0       string `json:"net0"`
	Serial0    string `json:"serial0"`
	Scsi0      string `json:"scsi0"`
	Agent      string `json:"agent"`
	Meta       string `json:"meta"`
	Digest     string `json:"digest"`
	Ide2       string `json:"ide2"`
	CPU        string `json:"cpu"`
	Name       string `json:"name"`
	Nameserver string `json:"nameserver"`
	IPConfig0  string `json:"ipconfig0"`
	CiCustom   string `json:"cicustom"`
	Cores      int    `json:"cores"`
	Boot       string `json:"boot"`
	VGA        string `json:"vga"`
	NUMA       int    `json:"numa"`
	SMBIOS1    string `json:"smbios1"`
	Memory     string `json:"memory"`
	ScsiHW     string `json:"scsihw"`
}

// type NodeList struct {
// 	vmid   string
// 	name   string
// 	status string
// }

type NodeList struct {
	Maxdisk        int64   `json:"maxdisk"`
	Maxcpu         int     `json:"maxcpu"`
	Node           string  `json:"node"`
	Disk           int64   `json:"disk"`
	Mem            int64   `json:"mem"`
	ID             string  `json:"id"`
	Level          string  `json:"level"`
	Status         string  `json:"status"`
	Maxmem         int64   `json:"maxmem"`
	Type           string  `json:"type"`
	SslFingerprint string  `json:"ssl_fingerprint"`
	Uptime         int     `json:"uptime"`
	CPU            float64 `json:"cpu"`
}

type VMList struct {
	Pid       int     `json:"pid"`
	Uptime    int     `json:"uptime"`
	Serial    int     `json:"serial"`
	CPU       float64 `json:"cpu"`
	Diskread  int     `json:"diskread"`
	Cpus      int     `json:"cpus"`
	Diskwrite int     `json:"diskwrite"`
	Name      string  `json:"name"`
	Vmid      int     `json:"vmid"`
	Netout    int     `json:"netout"`
	Mem       int     `json:"mem"`
	Disk      int     `json:"disk"`
	Maxmem    int64   `json:"maxmem"`
	Netin     int     `json:"netin"`
	Status    string  `json:"status"`
	Maxdisk   int64   `json:"maxdisk"`
}
type VMEdit struct {
	Vmid     int
	Node     string
	CPU      string
	Cores    int
	Boot     string
	Bootdisk string
	Ipconfig []string
	Memory   string
	Scsi     []string
}

type VMDelete struct {
	Vmid int
	Node string
}

type ResponsePVE[T any] struct {
	Data   T                 `json:"data"`
	Errors map[string]string `json:"errors,omitempty"`
}

func NewPVEService(conf *PVEConfig, client *http.Client) PVEService {
	return &pveService{
		pveConf:    conf,
		HTTPClient: client,
	}
}
func (s *pveService) GetNodeList() ([]NodeList, error) {
	endpoint := fmt.Sprintf("%s/nodes", s.pveConf.APIURL)

	formData := url.Values{}

	req, err := http.NewRequest("GET", endpoint, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, xerrors.Errorf("can't create http request: %w", err)
	}
	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	// json.Unmarshalでデコード
	var pveresp ResponsePVE[[]NodeList]
	if err := json.Unmarshal(body, &pveresp); err != nil {
		return nil, xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return nil, xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}

	return pveresp.Data, nil
}

func (s *pveService) GetVMList(nodes *NodeList) ([]VMList, error) {
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu", s.pveConf.APIURL, strings.Split(nodes.ID, "/")[1])
	formData := url.Values{}
	req, err := http.NewRequest("GET", endpoint, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, xerrors.Errorf("can't create http request: %w", err)
	}
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	var pveresp ResponsePVE[[]VMList]
	if err := json.Unmarshal(body, &pveresp); err != nil {
		return nil, xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return nil, xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}
	return pveresp.Data, nil
}

func (s *pveService) EditVM(vmedit VMEdit) error {
	// フォームデータの作成
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu/%d/config", s.pveConf.APIURL, vmedit.Node, vmedit.Vmid)
	// フォームデータの作成
	formData := url.Values{}
	formData.Set("cores", "2")
	for i, v := range vmedit.Ipconfig {
		formData.Set(fmt.Sprintf("ipconfig%d", i), v)
	}
	for i, v := range vmedit.Scsi {
		formData.Set(fmt.Sprintf("scsi%d", i), v)
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
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	// json.Unmarshalでデコード
	var pveresp ResponsePVE[string]
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

func (s *pveService) GetVM(id string) (*VMConfig, error) {
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu/%s/config", s.pveConf.APIURL, "pve02", id)
	formData := url.Values{}

	req, err := http.NewRequest("GET", endpoint, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, xerrors.Errorf("can't create http request: %w", err)
	}

	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return nil, xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}

	var getVMconfig ResponsePVE[VMConfig]

	err = json.Unmarshal([]byte(body), &getVMconfig)
	if err != nil {
		return nil, xerrors.Errorf("json Unmarshal error: %w", err)
	}
	fmt.Println(getVMconfig)

	return &getVMconfig.Data, nil
}

func (s *pveService) CloneVM(name string, id int, node string, cloneid int) error {
	// フォームデータの作成
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu/%d/clone", s.pveConf.APIURL, node, cloneid)
	fmt.Printf("%s/nodes/%s/qemu/%b/clone\n", s.pveConf.APIURL, node, cloneid)
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
	resp, err := s.HTTPClient.Do(req)
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
	var pveresp ResponsePVE[string]
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

func (s *pveService) DeleteVM(vmdelete *VMDelete) error {
	endpoint := fmt.Sprintf("%s/nodes/%s/qemu/%d/", s.pveConf.APIURL, vmdelete.Node, vmdelete.Vmid)
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return xerrors.Errorf("can't create http request: %w", err)
	}
	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// リクエストの送信
	resp, err := s.HTTPClient.Do(req)
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
	var pveresp ResponsePVE[string]
	if err := json.Unmarshal(body, &pveresp); err != nil {
		return xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}

	return nil

}
