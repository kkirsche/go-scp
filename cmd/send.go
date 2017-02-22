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
	"github.com/kkirsche/go-scp/libscp"
	"github.com/spf13/cobra"
)

var (
	addr,
	port,
	username,
	password,
	fp string
	verbose bool
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a file from the local machine to a remote machine",
	Long:  `Use to send a file from the local host to a remote host`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetLevel(logrus.InfoLevel)
		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		}

		for _, arg := range args {
			libscp.SendFileWithAgent(username, arg, port)
		}
	},
}

func init() {
	RootCmd.AddCommand(sendCmd)

	u, err := user.Current()
	if err != nil {
		logrus.WithError(err).Errorln("Failed to get current user")
	}

	sendCmd.PersistentFlags().StringVarP(&port, "port", "p", "22", "The port to connect to the remote host on")
	sendCmd.PersistentFlags().StringVarP(&username, "username", "u", u.Username, "The username to connect to the remote host with")
	sendCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode")
}
