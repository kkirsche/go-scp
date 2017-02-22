package scpClient

import "bytes"

// ExecuteCommand is used to
func (c *Client) ExecuteCommand(cmd string) (string, error) {
	err := c.VerifyClient()
	if err != nil {
		return "", err
	}

	// Don't allocate something we may not need, e.g. if we can't build a client
	var b bytes.Buffer

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	s, err := c.client.NewSession()
	if err != nil {
		return "", err
	}

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	s.Stdout = &b

	err = s.Run(cmd)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
