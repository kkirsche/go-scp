package libscp

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

// Key represents where an SSH Key should be read from. This is commonly
// used when the SSH agent is not used.
type Key struct {
	P string
	N string
}

// NewKey is used to create a new SSHKey object and validate that the key file
// actually exists on the system
func NewKey(path, name string) (*Key, error) {
	k := &Key{
		P: path,
		N: name,
	}

	if !k.Exists() {
		return nil, fmt.Errorf("Key file with path `%s` does not exist", k.Path())
	}

	return k, nil
}

// Exists validates that the path to the SSH keyfile is valid and that the key
// actually exists
func (s *Key) Exists() bool {
	_, err := os.Stat(s.Path())
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// Path returns the full path to the key file
func (s *Key) Path() string {
	return fmt.Sprintf("%s/%s", s.P, s.N)
}

// KeyConfig is used to create the SSH Client Configuration when
// using raw SSH key files rather than the SSH Agent for the authentication
// mechanism
func KeyConfig(u string, k *Key) (*ssh.ClientConfig, error) {
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
