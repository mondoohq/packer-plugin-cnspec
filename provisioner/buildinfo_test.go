package provisioner

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildInfo(t *testing.T) {

	data := `
		{"ConnType":"ssh","Host":"127.0.0.1","ID":"packer-alpine-1599666417","PackerHTTPAddr":"10.0.2.2:8644","PackerHTTPIP":"10.0.2.2","PackerHTTPPort":"8644","PackerRunUUID":"6ef27e2e-f4ea-be8e-7dc9-b1089998c001","Password":"vagrant","Port":2463,"SSHAgentAuth":false,"SSHPrivateKey":"priv-key","SSHPublicKey":"pub-key","User":"vagrant","WinRMPassword":""}
	`
	var buildInfo BuildInfo
	err := json.Unmarshal([]byte(data), &buildInfo)
	require.NoError(t, err)
	assert.Equal(t, "packer-alpine-1599666417", buildInfo.ID)
	assert.Equal(t, "ssh", buildInfo.ConnType)
	assert.Equal(t, "127.0.0.1", buildInfo.Host)
	assert.Equal(t, 2463, buildInfo.Port)
	assert.Equal(t, "vagrant", buildInfo.User)
	assert.Equal(t, "vagrant", buildInfo.Password)
	assert.Equal(t, "10.0.2.2:8644", buildInfo.PackerHTTPAddr)
	assert.Equal(t, "10.0.2.2", buildInfo.PackerHTTPIP)
	assert.Equal(t, "8644", buildInfo.PackerHTTPPort)
	assert.Equal(t, "6ef27e2e-f4ea-be8e-7dc9-b1089998c001", buildInfo.PackerRunUUID)
	assert.False(t, buildInfo.SSHAgentAuth)
	assert.Equal(t, "priv-key", buildInfo.SSHPrivateKey)
	assert.Equal(t, "pub-key", buildInfo.SSHPublicKey)
}
