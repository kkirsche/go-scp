package scpConfig

import (
	"github.com/kkirsche/go-scp/libscp/scpAgent"
	"golang.org/x/crypto/ssh"
)

// Agent is used to create the SSH Client Configuration when
// using the SSH Agent for the authentication mechanism
func Agent(u string) (*ssh.ClientConfig, error) {
	a, err := scpAgent.Get()
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: u,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(a.Signers),
		},
	}
	return config, nil
}
