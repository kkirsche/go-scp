package goScp

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func getAgent() (agent.Agent, error) {
	agentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	return agent.NewClient(agentConn), err
}

func withAgentSshConfig(username string) *ssh.ClientConfig {
	agent, err := getAgent()
	if err != nil {
		log.Println("Failed to connect to SSH_AUTH_SOCK:", err)
		os.Exit(1)
	}
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agent.Signers),
		},
	}
	return config
}

func withoutAgentSshConfig(username string, sshKeyFile SshKeyfile) *ssh.ClientConfig {
	keyFilePath := fmt.Sprintf("%s/%s", sshKeyFile.Path, sshKeyFile.Filename)
	keyFileContents, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	signer, err := ssh.ParsePrivateKey(keyFileContents)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	return config
}

func Connect(sshKeyFile SshKeyfile, sshCredentials SshCredentials, remoteMachine RemoteMachine, usingSshAgent bool) (*ssh.Client, error) {
	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig.
	var config *ssh.ClientConfig
	if usingSshAgent {
		config = withAgentSshConfig(sshCredentials.Username)
	} else {
		config = withoutAgentSshConfig(sshCredentials.Username, sshKeyFile)
	}

	client, err := ssh.Dial("tcp", remoteMachine.Host+":"+remoteMachine.Port, config)

	return client, err
}

func ExecuteCommand(client ssh.Client, cmd string) (string, error) {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(cmd); err != nil {
		return "", err
	}

	return b.String(), nil
}
