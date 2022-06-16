package provisioner

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildInfo(t *testing.T) {

	data := `
		{"ConnType":"ssh","Host":"127.0.0.1","ID":"packer-alpine-1599666417","PackerHTTPAddr":"10.0.2.2:8644","PackerHTTPIP":"10.0.2.2","PackerHTTPPort":"8644","PackerRunUUID":"6ef27e2e-f4ea-be8e-7dc9-b1089998c001","Password":"vagrant","Port":2463,"SSHAgentAuth":false,"SSHPrivateKey":"","SSHPrivateKeyFile":"","SSHPublicKey":"","User":"vagrant","WinRMPassword":""}
	`
	var buildInfo BuildInfo
	err := json.Unmarshal([]byte(data), &buildInfo)
	require.NoError(t, err)
	assert.Equal(t, "packer-alpine-1599666417", buildInfo.ID)
}
