package service

import (
	"fmt"
	"strconv"

	"github.com/LainInTheWired/ctf-backend/pveapi/model"
	"github.com/LainInTheWired/ctf-backend/pveapi/repository"
	"github.com/cockroachdb/errors"
)

type pveService struct {
	pveRepo repository.PVERepository
}

type PVEService interface {
	CreateCloudinitVM(name string, vmconf *model.VMEdit) error
	DeleteVMByVmid(vmid int) error
}

func NewPVEService(r repository.PVERepository) PVEService {
	return &pveService{
		pveRepo: r,
	}
}

// SelectLeastLoadedNode は最も負荷が低いノードを選択します
func (p *pveService) SelectNodeByCPU() (string, error) {
	nodes, err := p.pveRepo.GetNodeList()
	if err != nil {
		return "", errors.Wrap(err, "can't get nodes")
	}
	minLoad := 1.0 // CPU負荷は通常0.0〜1.0の範囲
	var selectedNode string
	for _, node := range nodes {
		// オンラインのノードのみ対象
		if node.Status != "online" {
			continue
		}

		// CPU負荷が最小のノードを選択
		if node.CPU < minLoad {
			minLoad = node.CPU
			selectedNode = node.Node
		}
	}

	if selectedNode == "" {
		return "", errors.New("not found online node")
	}
	return selectedNode, nil
}

func (p *pveService) CreateCloudinitVM(name string, vmconf *model.VMEdit) error {
	svmid, err := p.pveRepo.NextVMID()
	if err != nil {
		return errors.Wrap(err, "can't get next VID")
	}
	vmid, err := strconv.Atoi(svmid)
	if err != nil {
		return errors.Wrap(err, "can't cast vmid")
	}
	vmconf.Scsi = []string{fmt.Sprintf("vmdisk:vm-%d-disk-0,size=16G", vmid)}

	vmconf.Vmid = vmid
	err = p.pveRepo.CloneVM(name, vmid, "PVE01", 9000)
	if err != nil {
		return errors.Wrap(err, "can't clone vm")
	}

	// err = p.EditVM(*vmconf)
	err = p.fiveEditVM(vmconf)
	if err != nil {
		return errors.Wrap(err, "can't edit vm")
	}

	return nil
}

func (p *pveService) fiveEditVM(conf *model.VMEdit) error {
	for i := 0; i < 5; i++ {
		if err := p.pveRepo.EditVM(*conf); err != nil {
			if i > 4 {
				return errors.Wrap(err, "can't error")
			}
		} else {
			break
		}
	}
	return nil
}

func (p *pveService) SearchNodeByVmid(vmid int) (string, error) {
	res, err := p.pveRepo.GetClusterResourcesList()
	if err != nil {
		return "", err
	}
	for _, v := range res {
		id, err := strconv.Atoi(v.ID)
		if err != nil {
			return "", errors.Wrap(err, "can't ID Atoi")
		}
		if vmid == id {
			return v.Node, nil
		}
	}
	return "", errors.Wrap(err, "not fount this vmid in cluster")
}
func (p *pveService) DeleteVMByVmid(vmid int) error {
	n, err := p.SearchNodeByVmid(vmid)
	if err != nil {
		return errors.Wrap(err, "can't search node")
	}
	conf := &model.VMDelete{
		Vmid: vmid,
		Node: n,
	}

	if err := p.pveRepo.DeleteVM(conf); err != nil {
		return err
	}
	return nil

}
