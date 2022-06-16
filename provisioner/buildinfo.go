package provisioner

// https://www.packer.io/docs/templates/legacy_json_templates/engine
type BuildInfo struct {
	// depending on the cloud povider, the type changes
	ID       interface{} `json:"ID"`
	ConnType string      `json:"ConnType"`
	Host     string      `json:"Host"`
	Port     int         `json:"Port"`
	User     string      `json:"User"`

	PackerHTTPAddr string `json:"PackerHTTPAddr"`
	PackerHTTPIP   string `json:"PackerHTTPIP"`
	PackerHTTPPort string `json:"PackerHTTPPort"`
	PackerRunUUID  string `json:"PackerRunUUID"`
	Password       string `json:"Password"`

	SSHAgentAuth      bool   `json:"SSHAgentAuth"`
	SSHPrivateKey     string `json:"SSHPrivateKey"`
	SSHPrivateKeyFile string `json:"SSHPrivateKeyFile"`
	SSHPublicKey      string `json:"SSHPrivateKeyFile"`
}
