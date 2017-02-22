package libscp

import (
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// AgentConfig is used to create the SSH Client Configuration when
// using the SSH Agent for the authentication mechanism
func AgentConfig(u string) (*ssh.ClientConfig, error) {
	a, err := Get()
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

// Get is used to retrieve the agent connection for use via SSH
func Get() (agent.Agent, error) {
	agentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	return agent.NewClient(agentConn), err
}
