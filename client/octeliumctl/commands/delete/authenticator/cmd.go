// Copyright Octelium Labs, LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authenticator

import (
	"context"

	"github.com/octelium/octelium/apis/main/corev1"
	"github.com/octelium/octelium/apis/main/metav1"
	"github.com/octelium/octelium/client/common/client"
	"github.com/octelium/octelium/client/common/cliutils"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "authenticator [NAME]",
	Short: "Delete an Authenticator",
	Example: `
octeliumctl delete authn fido-123456
	`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"authn", "authen"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return doCmd(cmd, args)
	},
}

func doCmd(cmd *cobra.Command, args []string) error {
	i, err := cliutils.GetCLIInfo(cmd, args)
	if err != nil {
		return err
	}

	conn, err := client.GetGRPCClientConn(context.Background(), i.Domain)
	if err != nil {
		return err
	}

	defer conn.Close()
	c := corev1.NewMainServiceClient(conn)

	ctx := context.Background()

	if _, err := c.DeleteAuthenticator(ctx, &metav1.DeleteOptions{Name: i.FirstArg()}); err != nil {
		return err
	}

	cliutils.LineInfo("Authenticator `%s` successfully deleted\n", i.FirstArg())

	return nil
}
