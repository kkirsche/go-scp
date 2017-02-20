package scpAgent

import (
	"net"
	"os"

	"golang.org/x/crypto/ssh/agent"
)

// Get is used to retrieve the agent connection for use via SSH
func Get() (agent.Agent, error) {
	agentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	return agent.NewClient(agentConn), err
}
