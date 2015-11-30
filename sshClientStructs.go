package goScp

type SshCredentials struct {
	Username string
	Password string
}

type RemoteMachine struct {
	Host string
	Port string
}

type SshKeyfile struct {
	Path     string
	Filename string
}
