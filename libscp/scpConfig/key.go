package scpConfig

import (
	"io/ioutil"

	"github.com/kkirsche/go-scp/libscp/scpAuth"
	"golang.org/x/crypto/ssh"
)

// Key is used to create the SSH Client Configuration when
// using raw SSH key files rather than the SSH Agent for the authentication
// mechanism
func Key(u string, k *scpAuth.Key) (*ssh.ClientConfig, error) {
	contents, err := ioutil.ReadFile(k.Path())
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(contents)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: u,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	return config, nil
}
