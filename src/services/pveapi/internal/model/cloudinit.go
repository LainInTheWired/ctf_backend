package model

// CloudConfig represents the top-level cloud-init configuration
type CloudConfig struct {
	Hostname string   `yaml:"hostname,omitempty"`
	Users    []User   `yaml:"users,omitempty"`
	Packages []string `yaml:"packages,omitempty"`
	RunCmd   []string `yaml:"runcmd,omitempty"`
}

// User represents a user configuration in cloud-init
type User struct {
	Name              string   `yaml:"name"`
	Sudo              string   `yaml:"sudo,omitempty"`
	Passwd            string   `yaml:"passwd"`
	Groups            string   `yaml:"groups,omitempty"`
	Shell             string   `yaml:"shell,omitempty"`
	SshAuthorizedKeys []string `yaml:"ssh-authorized-keys,omitempty"`
}
