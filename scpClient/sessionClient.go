package scpClient

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/kkirsche/go-scp/scpFile"

	"golang.org/x/crypto/ssh"
)

// SessionClient is used to hold the writer for stdout, reader from stdin,
// wait group for the concurrency, and error channel
type SessionClient struct {
	writer io.WriteCloser
	reader io.Reader
	wg     *sync.WaitGroup
	errors chan error
}

// NewSessionClient creates a new SessionClient structure from the ssh.Session
// pointer
func NewSessionClient(s *ssh.Session) (*SessionClient, error) {
	writer, err := s.StdinPipe()
	if err != nil {
		return nil, err
	}

	reader, err := s.StdoutPipe()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	return &SessionClient{
		writer: writer,
		reader: reader,
		wg:     &wg,
		errors: make(chan error),
	}, nil
}

// FileSink is used to receive a file from the remote machine and save it to the local machine
func (c *SessionClient) FileSink(fp, fn string) {
	// We must close the channel for the main thread to work properly. Defer ensures
	// when the function ends, this is closed. We also want to be sure we mark this
	// attempt as completed
	defer close(c.errors)
	defer c.wg.Done()

	logrus.Debugln("Beginning transfer")
	successfulByte := []byte{0}
	// Send a null byte saying that we are ready to receive the data
	c.writer.Write(successfulByte)

	// We want to first receive the command input from remote machine
	// e.g. C0644 113828 test.csv
	scpCommandArray := make([]byte, 500)
	bytesRead, err := c.reader.Read(scpCommandArray)
	if err != nil {
		if err == io.EOF {
			//no problem.
		} else {
			c.errors <- err
			return
		}
	}

	scpStartLine := string(scpCommandArray[:bytesRead])
	scpStartLineArray := strings.Split(scpStartLine, " ")

	filePermission := scpStartLineArray[0][1:]
	fileSize := scpStartLineArray[1]
	fileName := scpStartLineArray[2]

	logrus.Debugf("File with permissions: %s, File Size: %s, File Name: %s", filePermission, fileSize, fileName)

	// Confirm to the remote host that we have received the command line
	c.writer.Write(successfulByte)

	// Now we want to start receiving the file itself from the remote machine
	// one byte at a time
	fileContents := make([]byte, 1)

	var file *os.File
	if fn == "" {
		file, err = scpFile.Create(fmt.Sprintf("%s/%s", fp, fn))
		if err != nil {
			c.errors <- err
			return
		}
	} else {
		file, err = scpFile.Create(fmt.Sprintf("%s/%s", fp, fn))
		if err != nil {
			c.errors <- err
			return
		}
	}

	more := true
	for more {
		bytesRead, err = c.reader.Read(fileContents)
		if err != nil {
			if err == io.EOF {
				more = false
			} else {
				c.errors <- err
				return
			}
		}
		_, err = scpFile.WriteBytes(file, fileContents[:bytesRead])
		if err != nil {
			c.errors <- err
			return
		}
		c.writer.Write(successfulByte)
	}
	err = file.Sync()
	if err != nil {
		c.errors <- err
		return
	}
}

// FileSource allows us to acting as the machine sending a file to the remote host
func (c *SessionClient) FileSource(p string) {
	startByte := []byte{0}
	defer close(c.errors)
	defer c.wg.Done()

	logrus.Debugln("Opening file to send")
	f, err := os.Open(p)
	if err != nil {
		c.errors <- err
		return
	}
	defer f.Close()

	logrus.Debugln("Getting file information")
	i, err := f.Stat()
	if err != nil {
		c.errors <- err
		return
	}

	logrus.WithFields(logrus.Fields{
		"directory":         i.IsDir(),
		"modification_time": i.ModTime().String(),
		"mode":              i.Mode().String(),
		"name":              i.Name(),
		"size":              i.Size(),
	}).Debugln("File information retrieved")

	pa := strings.Split(p, " ")
	begin := []byte(fmt.Sprintf("%s %d %s", "C0644", i.Size(), pa[len(pa)-1]))
	logrus.WithField("statement", string(begin)).Debugln("Beginning transfer")
	c.writer.Write(begin)
	c.writer.Write(startByte)
	ready := make([]byte, 1)
	_, err = c.reader.Read(ready)
	if err != nil {
		c.errors <- err
		return
	}
	logrus.WithField("response", ready).Debugln("Ready to begin transfer")
	r := bufio.NewReader(f)
	logrus.Debugln("Beginning to read file")

fileLoop:
	for err != io.EOF {
		b, err := r.ReadBytes('\n')
		logrus.WithField("line", string(b)).Debugln("Read line")
		if err == io.EOF {
			logrus.WithError(err).Debugln("EOF found")
			break fileLoop
		} else if err != nil {
			logrus.WithError(err).Debugln("Error occurred in fileloop")
			c.errors <- err
			return
		}
		if len(b) > 0 {
			w, err := c.writer.Write(b)
			if err != nil {
				c.errors <- err
				return
			}

			logrus.WithField("bytes", w).Debugf("Wrote %d bytes", w)
			rcvd := make([]byte, 1)
			_, err = c.reader.Read(rcvd)
			if err != nil {
				logrus.WithError(err).Debugln("Error getting response")
			}
			logrus.WithField("response", rcvd).Debugln("Recieved response")
		}
	}

	logrus.Debugln("Transfer complete")
	c.writer.Write([]byte("\x00")) // transfer end with \x00\
	return
}
