package model

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
	Memory   int
	Scsi     []string
	Cicustom string
}

type VMDelete struct {
	Vmid int
	Node string
}

type ResponsePVE[T any] struct {
	Data   T                 `json:"data"`
	Errors map[string]string `json:"errors,omitempty"`
}

type ClusterResources struct {
	Uptime    int     `json:"uptime"`
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Maxmem    int64   `json:"maxmem"`
	Node      string  `json:"node"`
	Status    string  `json:"status"`
	Maxcpu    int     `json:"maxcpu"`
	Netin     int     `json:"netin"`
	Mem       int     `json:"mem"`
	Template  int     `json:"template"`
	Diskread  int     `json:"diskread"`
	Type      string  `json:"type"`
	Diskwrite int     `json:"diskwrite"`
	Maxdisk   int64   `json:"maxdisk"`
	CPU       float64 `json:"cpu"`
	Disk      int     `json:"disk"`
	Netout    int     `json:"netout"`
	Vmid      int     `json:"vmid"`
}

type NetworkIntQumeAgent struct {
	Statistics struct {
		RxBytes   int `json:"rx-bytes"`
		RxDropped int `json:"rx-dropped"`
		RxErrs    int `json:"rx-errs"`
		RxPackets int `json:"rx-packets"`
		TxPackets int `json:"tx-packets"`
		TxErrs    int `json:"tx-errs"`
		TxDropped int `json:"tx-dropped"`
		TxBytes   int `json:"tx-bytes"`
	} `json:"statistics"`
	Name            string `json:"name"`
	HardwareAddress string `json:"hardware-address"`
	IPAddresses     []struct {
		IPAddressType string `json:"ip-address-type"`
		Prefix        int    `json:"prefix"`
		IPAddress     string `json:"ip-address"`
	} `json:"ip-addresses"`
}

// type ClusterResources struct {
// 	Maxcpu    int     `json:"maxcpu"`
// 	Vmid      int     `json:"vmid"`
// 	Name      string  `json:"name"`
// 	Maxmem    int64   `json:"maxmem"`
// 	Disk      int     `json:"disk"`
// 	Diskwrite int64   `json:"diskwrite"`
// 	Maxdisk   int64   `json:"maxdisk"`
// 	Type      string  `json:"type"`
// 	ID        string  `json:"id"`
// 	Node      string  `json:"node"`
// 	Uptime    int     `json:"uptime"`
// 	Diskread  int     `json:"diskread"`
// 	Netin     int64   `json:"netin"`
// 	Status    string  `json:"status"`
// 	CPU       float64 `json:"cpu"`
// 	Netout    int64   `json:"netout"`
// 	Mem       int64   `json:"mem"`
// 	Template  int     `json:"template"`
// }

// type ClusterResources struct {
// 	Uptime int64 `json:"uptime"`
// 	// ID        string  `json:"id"`
// 	Name      string  `json:"name"`
// 	Maxmem    int64   `json:"maxmem"`
// 	Node      string  `json:"node"`
// 	Status    string  `json:"status"`
// 	Maxcpu    int     `json:"maxcpu"`
// 	Netin     int64   `json:"netin"`
// 	Mem       int64   `json:"mem"`
// 	Template  int     `json:"template"`
// 	Diskread  int64   `json:"diskread"`
// 	Type      string  `json:"type"`
// 	Diskwrite int64   `json:"diskwrite"`
// 	Maxdisk   int64   `json:"maxdisk"`
// 	CPU       float64 `json:"cpu"`
// 	Disk      int     `json:"disk"`
// 	Netout    int64   `json:"netout"`
// 	Vmid      int     `json:"vmid"`
// }
