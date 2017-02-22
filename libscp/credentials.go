package libscp

// Credentials are the SSH credentials that should be used to connect to the
// remote host. This is for use with the SSH Agent.
type Credentials struct {
	Username string
	Password string
}

// NewCredentials creates a new credential object
func NewCredentials(u, p string) *Credentials {
	return &Credentials{
		Username: u,
		Password: p,
	}
}
