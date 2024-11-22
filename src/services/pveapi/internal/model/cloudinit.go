package model

// CloudConfig represents the top-level cloud-init configuration
type CloudinitConfig struct {
	Hostname  string   `yaml:"hostname,omitempty"`
	FQDN      string   `yaml:"fqdn,omitempty"`
	SshPwauth int      `yaml:"ssh_pwauth,omitempty"`
	Users     []User   `yaml:"users,omitempty"`
	Packages  []string `yaml:"packages,omitempty"`
	RunCmd    []string `yaml:"runcmd,omitempty"`
}

// User represents a user configuration in cloud-init
type User struct {
	Name string `yaml:"name"`
	Sudo string `yaml:"sudo,omitempty"`
	// Passwd            string   `yaml:"passwd,omitempty"`
	PlainTextPasswd   string   `yaml:"plain_text_passwd,omitempty"`
	Groups            string   `yaml:"groups,omitempty"`
	Shell             string   `yaml:"shell,omitempty"`
	LockPasswd        bool     `yaml:"lock_passwd"`
	SshAuthorizedKeys []string `yaml:"ssh-authorized-keys,omitempty"`
	SshPwauth         string   `yaml:"ssh_pwauth,omitempty"`
}

type QuesionRequest struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CategoryID  int      `json:"category_id"`
	Env         string   `json:"env"`
	Sshkeys     []string `json:"sshkeys"`
	Memory      int      `json"memory"`
	CPUs        int      `json"cpu"`
	Disk        int      `json"disk"`
	IP          string   `json:"ip" validate:"cidr"`
	Gateway     string   `json:"gateway" validate:"ip"`
	Filename    string   `json:"filename"`
}
