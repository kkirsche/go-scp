package libscp

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/kkirsche/go-scp/libscp/scpAuth"
	"github.com/kkirsche/go-scp/libscp/scpClient"
	"github.com/kkirsche/go-scp/libscp/scpFile"
)

func SendFileWithAgent(username, arg, port string) (logrus.Fields, error) {
	res := strings.Split(arg, ":")
	fname := res[1]

	fp, err := scpFile.ExpandPath(fname)
	if err != nil {
		logrus.WithError(err).Errorln("Failed to get current directory")
		return logrus.Fields{}, nil
	}

	creds := scpAuth.NewCredentials(username, "")
	a := scpClient.NewAgentClient(res[0], port, creds)
	err = a.SendFileToRemote(fp)
	if err != nil {
		return logrus.Fields{
			"address":  res[0],
			"port":     port,
			"file":     fp,
			"username": username,
		}, err
	}

	return logrus.Fields{}, nil
}
