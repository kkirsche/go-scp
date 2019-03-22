package libscp

import "github.com/sirupsen/logrus"

// ReceiveFileFromRemote is used to receive a file using the SCP protocol from a
// remote host / machine
func (c *Client) ReceiveFileFromRemote(remoteFilePath string, remoteFilename string, localFilePath string, localFileName string) error {
	err := c.VerifyClient()
	if err != nil {
		return err
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	sc, err := NewSessionClient(session)
	if err != nil {
		return err
	}
	defer sc.writer.Close()

	sc.wg.Add(1)
	go sc.FileSink(localFilePath, localFileName)
	session.Run("/usr/bin/scp -f " + remoteFilePath + "/" + remoteFilename)

	for err := range sc.errors {
		logrus.WithError(err).Errorln("Error receiving the file")
	}
	sc.wg.Wait()

	return nil
}

// SendFileToRemote is used to send a file using the SCP protocol to a remote
// host or machine
func (c *Client) SendFileToRemote(fp string) error {
	logrus.Debugln("Verifying client")
	err := c.VerifyClient()
	if err != nil {
		return err
	}
	logrus.Debugln("Client passed verification")

	logrus.Debugln("Creating session")
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	s, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer s.Close()

	logrus.Debugln("Creating session client")
	sc, err := NewSessionClient(s)
	if err != nil {
		return err
	}
	defer sc.writer.Close()

	sc.wg.Add(1)
	go sc.FileSource(fp)
	logrus.Infoln("Beginning transfer")
	go s.Run("/usr/bin/scp -t ./")
	logrus.Debugln("Waiting for transfer to complete...")
	sc.wg.Wait()
	for err := range sc.errors {
		logrus.WithError(err).Errorln("Error sending the file")
	}
	return nil
}
