package goScp

// SSHCredentials are the SSH credentials that should be used to connect to the
// remote host. This is for use with the SSH Agent.
type SSHCredentials struct {
	Username string
	Password string
}

// RemoteHost is the remote machine that should be connected to. Specifically,
// what hostname and port.
type RemoteHost struct {
	Host string
	Port string
}

// SSHKeyfile represents where an SSH Key should be read from. This is used when
// the SSH agent is not used.
type SSHKeyfile struct {
	Path     string
	Filename string
}
