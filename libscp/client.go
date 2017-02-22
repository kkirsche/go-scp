package libscp

import (
	"net"

	"golang.org/x/crypto/ssh"
)

// Client is the remote machine that should be connected to for the transfer.
// Specifically, what hostname and port.
type Client struct {
	addr        string
	port        string
	useAgent    bool
	key         *Key
	credentials *Credentials
	client      *ssh.Client
}

// NewAgentClient creates a new host object that will connect using the SSH
// Agent signers
func NewAgentClient(addr, port string, creds *Credentials) *Client {
	return &Client{
		addr:        addr,
		port:        port,
		useAgent:    true,
		credentials: creds,
	}
}

// NewKeyClient creates a new host object that will connect using the SSH
// Agent signers
func NewKeyClient(addr, port string, creds *Credentials, key *Key) *Client {
	return &Client{
		addr:        addr,
		port:        port,
		useAgent:    false,
		credentials: creds,
		key:         key,
	}
}

// Connect takes the host object and connects, creating an SSH Client connection
func (c *Client) Connect() error {
	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig.
	var co *ssh.ClientConfig
	var err error

	if c.useAgent {
		co, err = AgentConfig(c.credentials.Username)
		if err != nil {
			return err
		}
	} else {
		co, err = KeyConfig(c.credentials.Username, c.key)
		if err != nil {
			return err
		}
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(c.addr, c.port), co)
	if err != nil {
		return err
	}

	c.client = client

	return nil
}

// VerifyClient checks if we have a client, and if not attempts to connect
func (c *Client) VerifyClient() error {
	// If we don't have a client yet, we should try to create one
	if c.client == nil {
		err := c.Connect()
		if err != nil {
			return err
		}
	}

	return nil
}
