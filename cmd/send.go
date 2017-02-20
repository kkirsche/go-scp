// Copyright © 2017 Kevin Kirsche <kev.kirsche[at]gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os/user"

	"github.com/Sirupsen/logrus"
	"github.com/kkirsche/go-scp/scpAuth"
	"github.com/kkirsche/go-scp/scpClient"
	"github.com/spf13/cobra"
)

var (
	addr,
	port,
	username,
	password,
	fp string
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetLevel(logrus.DebugLevel)
		creds := scpAuth.NewCredentials(username, "")
		a := scpClient.NewAgentClient(addr, port, creds)
		err := a.SendFileToRemote(fp)
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"address":  addr,
				"port":     port,
				"file":     fp,
				"username": username,
			}).Errorln("Failed to send file to host")
		}
	},
}

func init() {
	RootCmd.AddCommand(sendCmd)

	u, err := user.Current()
	if err != nil {
		logrus.WithError(err).Errorln("Failed to get current user")
	}

	sendCmd.PersistentFlags().StringVarP(&addr, "address", "a", "", "The remote IP address to send to")
	sendCmd.PersistentFlags().StringVarP(&port, "port", "p", "22", "The port to connect to the remote host on")
	sendCmd.PersistentFlags().StringVarP(&fp, "filepath", "f", "", "The path to the file to send")
	sendCmd.PersistentFlags().StringVarP(&username, "username", "u", u.Username, "The username to connect to the remote host with")
}
