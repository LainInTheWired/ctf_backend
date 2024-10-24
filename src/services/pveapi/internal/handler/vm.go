package handler

import (
	"fmt"
	"net/http"

	"github.com/LainInTheWired/ctf-backend/pveapi/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/xerrors"
)

type PVEParam struct {
	name string
	ip   string
}

// 依存関係用の構造体
type PVEHandler struct {
	serv service.PVEService
}

type createVMRequest struct {
	Name        string `json"name" validate:"required"`
	Memory      string `json"memory" validate:"required" `
	CPUs        string `json"cpu" validate:"required"`
	OS          string `json"os" validate:"required"`
	Description string `json"description" validate:"required"`
}

type deleteVMRequest struct {
	Name string `json"name" validate:"required"`
}

// // VMConfig は作成するVMの設定を保持します
// type VMConfig struct {
// 	VMID           string `json:"vmid"`
// 	Name           string `json:"name"`
// 	Memory         string `json:"memory"` // MB単位
// 	CPUs           string `json:"cores"`
// 	Net0           string `json:"net0"`        // 例: "virtio=DE:AD:BE:EF:00:00,bridge=vmbr0"
// 	Scsi0          string `json:"scsi0"`       // 例: "kingston_1tb:vm-200-disk-0,size=16G"
// 	Boot           string `json:"boot"`        // 例: "c"
// 	Ide2           string `json:"ide2"`        // 例: "local:iso/AlmaLinux-9.3-x86_64-boot.iso"
// 	OSType         string `json:"ostype"`      // 例: "l26" (Linux 2.6/3.x/4.x)
// 	SCSIController string `json:"scsihw"`      // 例: "virtio-scsi-single"
// 	Description    string `json:"description"` // VMの説明（オプション）
// }

func NewPVEAPI(s service.PVEService) *PVEHandler {
	return &PVEHandler{
		serv: s,
	}
}

func (h *PVEHandler) GetVM(c echo.Context) error {
	_, err := h.serv.GetVM("200")
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "create vm success"))
}

func (h *PVEHandler) CreateCloudinitVM(c echo.Context) error {
	// リクエストから構造体にデータをコピー
	var req createVMRequest
	if err := c.Bind(&req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	nodes, err := h.serv.GetNodeList()
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	IDs := map[int]bool{}

	for _, n := range nodes {
		vms, err := h.serv.GetVMList(&n)
		if err != nil {
			wrappedErr := xerrors.Errorf(": %w", err)
			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
		}
		for _, v := range vms {
			// i, err := strconv.Atoi(v.Vmid)
			// if err != nil {
			// 	wrappedErr := xerrors.Errorf(": %w", err)
			// 	log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
			// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
			// }
			IDs[v.Vmid] = true
		}
	}

	cvmid := 100
	for {
		if !IDs[cvmid] {
			break
		}
		cvmid++
	}
	err = h.serv.CloneVM(req.Name, cvmid, "pve02", 9000)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	conf := service.VMEdit{
		Vmid:     cvmid,
		Memory:   "2mb",
		Cores:    2,
		Node:     "pve02",
		Ipconfig: []string{"ip=192.168.11.90/24,gw=192.168.11.1"},
		Scsi:     []string{fmt.Sprintf("local-lvm:vm-%d-disk-0,size=16G", cvmid)},
	}

	err = h.serv.EditVM(conf)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "create vm success"))
}

func (h *PVEHandler) DeleteVM(c echo.Context) error {
	// リクエストから構造体にデータをコピー
	var req createVMRequest
	if err := c.Bind(&req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	conf := &service.VMDelete{
		Vmid: 107,
		Node: "pve02",
	}

	if err := h.serv.DeleteVM(conf); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return nil
}

func (h *PVEHandler) GETTestHander(c echo.Context) error {
	err := h.serv.EditVM(service.VMEdit{})
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "test success"))
}

func (h *PVEHandler) GetNodeTestHander(c echo.Context) error {
	nodes, err := h.serv.GetNodeList()
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	for _, n := range nodes {
		h.serv.GetVMList(&n)
	}
	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "test success"))
}
