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

func withAgentSshConfig(username string) (*ssh.ClientConfig, error) {
	agent, err := getAgent()
	if err != nil {
		return &ssh.ClientConfig{}, err
	}
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agent.Signers),
		},
	}
	return config, nil
}

func withoutAgentSshConfig(username string, sshKeyFile SshKeyfile) (*ssh.ClientConfig, error) {
	keyFilePath := fmt.Sprintf("%s/%s", sshKeyFile.Path, sshKeyFile.Filename)
	keyFileContents, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return &ssh.ClientConfig{}, err
	}
	signer, err := ssh.ParsePrivateKey(keyFileContents)
	if err != nil {
		return &ssh.ClientConfig{}, err
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	return config, nil
}

func Connect(sshKeyFile SshKeyfile, sshCredentials SshCredentials, remoteMachine RemoteMachine, usingSshAgent bool) (*ssh.Client, error) {
	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig.
	var config *ssh.ClientConfig
	var err error
	if usingSshAgent {
		config, err = withAgentSshConfig(sshCredentials.Username)
	} else {
		config, err = withoutAgentSshConfig(sshCredentials.Username, sshKeyFile)
	}

	client, err := ssh.Dial("tcp", remoteMachine.Host+":"+remoteMachine.Port, config)

	return client, err
}

func ExecuteCommand(client ssh.Client, cmd string) (string, error) {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: " + err.Error())
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
