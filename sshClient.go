package goScp

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

const (
	VERSION = "0.0.1"
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

func ExecuteCommand(client *ssh.Client, cmd string) (string, error) {
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

func CopyRemoteFileToLocal(client *ssh.Client, remoteFilePath string, filename string) error {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: " + err.Error())
	}
	defer session.Close()

	writer, err := session.StdinPipe()
	if err != nil {
		return err
	}

	reader, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	doneChannel := make(chan bool)

	go func(writer io.WriteCloser, reader io.Reader, doneChannel chan bool) {
		writer.Write([]byte{0})
		// Length of
		scpCommandArray := make([]byte, 100)
		bytes_read, err := reader.Read(scpCommandArray)
		log.Print(string(scpCommandArray[:bytes_read]))
		if err != nil {
			if err == io.EOF {
				//no problem.
			} else {
				log.Fatalf("Error reading standard input: %s", err.Error())
			}
		}

		writer.Write([]byte{0})
		fileContents := make([]byte, 100)
		more := true
		for more {
			bytes_read, err = reader.Read(fileContents)
			// log.Print(string(fileContents[:bytes_read]))
			if err != nil {
				if err == io.EOF {
					more = false
				} else {
					log.Fatalf("Error reading standard input: %s", err.Error())
				}
			}
			writer.Write([]byte{0})
		}
		doneChannel <- true
	}(writer, reader, doneChannel)

	log.Print("Attempting: /usr/bin/scp -f " + remoteFilePath + "/" + filename)
	session.Run("/usr/bin/scp -f " + remoteFilePath + "/" + filename)
	<-doneChannel
	writer.Close()
	return nil
}

func CopyLocalFileToRemote(client *ssh.Client, localFilePath string, filename string) error {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: " + err.Error())
	}
	defer session.Close()

	writer, err := session.StdinPipe()
	if err != nil {
		return err
	}
	defer writer.Close()

	go func() {
		fileContents, _ := ioutil.ReadFile(localFilePath + "/" + filename)
		content := string(fileContents)
		fmt.Fprintln(writer, "C0644", len(content), filename)
		fmt.Fprint(writer, content)
		fmt.Fprintln(writer, "\x00") // transfer end with \x00\
	}()

	session.Run("/usr/bin/scp -t ./")
	return nil
}
