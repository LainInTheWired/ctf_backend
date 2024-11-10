package handler

import (
	"fmt"
	"net/http"

	"github.com/LainInTheWired/ctf-backend/pveapi/model"
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
	CPUs        int    `json"cpu" validate:"required"`
	IP          string `json:"ip" validate:"required,cidr"`
	Gateway     string `json:"ip" validate:"required,ip"`
	OS          string `json"os" validate:"required"`
	Description string `json"description" validate:"required"`
}
type CloneQuestionsRequest struct {
	QID   int    `json"qid" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Teams []TeamRequest
}
type TeamRequest struct {
	ID      int      `json:"id"`
	Sshkeys []string `json:"sshkeys"`
}
type deleteVMRequest struct {
	ID int `json"id" validate:"required"`
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

// func (h *PVEHandler) GetVM(c echo.Context) error {
// 	_, err := h.serv.zGetVM("200")
// 	if err != nil {
// 		wrappedErr := xerrors.Errorf(": %w", err)
// 		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 	}
// 	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "create vm success"))
// }

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

	conf := &model.VMEdit{
		Memory:   fmt.Sprintf("%dmb", req.Memory),
		Cores:    req.CPUs,
		Node:     "PVE01",
		Ipconfig: []string{fmt.Sprintf("ip=%s,gw=%s", req.IP, req.Gateway)},
	}
	err := h.serv.CreateCloudinitVM(req.Name, conf)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})

	}
	return c.JSON(http.StatusOK, fmt.Sprintf("message", "success create vm"))
}

func (h *PVEHandler) DeleteVM(c echo.Context) error {
	// リクエストから構造体にデータをコピー
	var req deleteVMRequest
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

	if err := h.serv.DeleteVMByVmid(req.ID); err != nil {
		return err
	}
	return nil
}

// func (h *PVEHandler) GETTestHander(c echo.Context) error {
// 	err := h.serv.EditVM(service.VMEdit{})
// 	if err != nil {
// 		wrappedErr := xerrors.Errorf(": %w", err)
// 		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 	}

// 	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "test success"))
// }

// func (h *PVEHandler) GetNodeTestHander(c echo.Context) error {
// 	nodes, err := h.serv.GetNodeList()
// 	if err != nil {
// 		wrappedErr := xerrors.Errorf(": %w", err)
// 		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 	}
// 	for _, n := range nodes {
// 		h.serv.GetVMList(&n)
// 	}
// 	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "test success"))
// }

// func (h *PVEHandler) CloneQuestions(c echo.Context) error {
// 	// リクエストから構造体にデータをコピー
// 	var req CloneQuestionsRequest
// 	if err := c.Bind(&req); err != nil {
// 		wrappedErr := xerrors.Errorf(": %w", err)
// 		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 	}
// 	// データをバリデーションにかける
// 	if err := c.Validate(req); err != nil {
// 		wrappedErr := xerrors.Errorf(": %w", err)
// 		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 	}

// 	var wg sync.WaitGroup

// 	wg.Add(len(req.Teams))

// 	for _, team := range req.Teams {
// 		hostname := fmt.Sprintf("ctf-%s-%d", "contestName", team.ID)
// 		users := []service.User{{
// 			Name:              "user",
// 			Sudo:              "ALL=(ALL) NOPASSWD:ALL",
// 			Shell:             "/bin/bash",
// 			Passwd:            "user",
// 			SshAuthorizedKeys: team.Sshkeys,
// 		},
// 		}
// 		filename := fmt.Sprintf("ctf-%d-%d.yml", req.QID, team.ID)
// 		err := h.serv.CloudinitGenerator(filename, hostname, users)
// 		if err != nil {
// 			wrappedErr := xerrors.Errorf(": %w", err)
// 			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 		}
// 		err = h.serv.TransferFileViaSCP(filename)
// 		if err != nil {
// 			wrappedErr := xerrors.Errorf(": %w", err)
// 			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 		}

// 		nextid, err := h.serv.NextVMID()
// 		if err != nil {
// 			wrappedErr := xerrors.Errorf(": %w", err)
// 			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 		}
// 		fmt.Println(nextid)

// 		nid, err := h.serv.NextVMID()
// 		if err != nil {
// 			wrappedErr := xerrors.Errorf(": %w", err)
// 			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 		}

// 		name := fmt.Sprintf("ctf-%d-%d", req.QID, team.ID)
// 		n, err := strconv.Atoi(nid)
// 		if err != nil {
// 			wrappedErr := xerrors.Errorf(": %w", err)
// 			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 		}

// 		if err = h.serv.CloneVM(name, n, "PVE01", 9000); err != nil {
// 			wrappedErr := xerrors.Errorf(": %w", err)
// 			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 		}
// 		conf := service.VMEdit{
// 			Vmid:     n,
// 			Memory:   "2mb",
// 			Cores:    2,
// 			Node:     "PVE01",
// 			Ipconfig: []string{"ip=10.0.10.90/24,gw=10.0.10.1"},
// 			Scsi:     []string{fmt.Sprintf("vmdisk:vm-%d-disk-0,size=16G", n)},
// 			Cicustom: filename,
// 		}
// 		for i := 0; i < 5; i++ {
// 			fmt.Println(i)
// 			if err = h.serv.EditVM(conf); err != nil {
// 				if i > 4 {
// 					wrappedErr := xerrors.Errorf(": %w", err)
// 					log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 					return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 				}
// 			} else {
// 				break
// 			}
// 		}

// 	}

// 	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "create vm success"))
// }
