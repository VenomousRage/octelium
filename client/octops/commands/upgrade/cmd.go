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

package upgrade

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-version"
	"github.com/octelium/octelium/apis/main/corev1"
	"github.com/octelium/octelium/apis/main/metav1"
	"github.com/octelium/octelium/client/common/client"
	"github.com/octelium/octelium/client/common/cliutils"
	"github.com/octelium/octelium/client/octops/commands/initcmd"
	"github.com/octelium/octelium/pkg/utils/ldflags"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

type args struct {
	KubeConfigFilePath string
	KubeContext        string
	Version            string
	CheckMode          bool
	CheckModeClient    bool
	Wait               bool
}

var examples = `
octops upgrade example.com
octops upgrade octelium.example.com --kubeconfig /path/to/kueconfig
octops upgrade sub.octelium.example.com  --kubeconfig /path/to/kueconfig
`

var Cmd = &cobra.Command{
	Use:     "upgrade [DOMAIN]",
	Short:   "Upgrade your Octelium Cluster",
	Args:    cobra.ExactArgs(1),
	Example: examples,
	RunE: func(cmd *cobra.Command, args []string) error {
		return doCmd(cmd, args)
	},
}

var cmdArgs args

func init() {
	Cmd.PersistentFlags().StringVar(&cmdArgs.KubeConfigFilePath, "kubeconfig", "", "kubeconfig file path")
	Cmd.PersistentFlags().StringVar(&cmdArgs.KubeContext, "kubecontext", "", "kubecontext")

	Cmd.PersistentFlags().BoolVar(&cmdArgs.CheckMode, "check", false,
		"Check whether there is a more recent latest release without actually upgrading the Cluster")
	Cmd.PersistentFlags().BoolVar(&cmdArgs.CheckModeClient, "check-client", false,
		"Check whether there is a more recent latest release for Octelium CLIs")
	Cmd.PersistentFlags().StringVar(&cmdArgs.Version, "version", "",
		`The desired Octelium Cluster version. By default it is set to "latest"`)

	Cmd.PersistentFlags().BoolVar(&cmdArgs.Wait, "wait", false,
		"Wait until the upgrade is complete")
}

func doCmd(cmd *cobra.Command, args []string) error {

	ctx := cmd.Context()

	clusterDomain := args[0]

	switch {
	case cmdArgs.CheckMode:
		return doCheck(ctx, clusterDomain)
	case cmdArgs.CheckModeClient:
		return doCheckClient(ctx)
	}

	cfg, err := initcmd.BuildConfigFromFlags("", cmdArgs.KubeConfigFilePath)
	if err != nil {
		return err
	}

	k8sC, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	if err := cliutils.RunPromptConfirm("Confirm to proceed with the Cluster upgrade"); err != nil {
		return err
	}

	if err := createGenesis(ctx, k8sC, clusterDomain, cmdArgs.Version); err != nil {
		return err
	}

	cliutils.LineNotify("Upgrading the Cluster has started.\n")

	if cmdArgs.Wait {
		s := cliutils.NewSpinner(os.Stdout)
		s.SetSuffix("Waiting for the Cluster upgrade to finish")
		s.Start()

		conn, err := client.GetGRPCClientConn(ctx, clusterDomain)
		if err != nil {
			return err
		}
		defer conn.Close()

		c := corev1.NewMainServiceClient(conn)

		versionInit, err := getCurrentClusterVersion(ctx, c)
		if err != nil {
			return err
		}

		for range 1000 {
			versionCurrent, err := getCurrentClusterVersion(ctx, c)
			if err == nil && !versionCurrent.Equal(versionInit) {
				s.Stop()
				cliutils.LineNotify("Cluster has been upgraded.\n")
				return nil
			} else if err != nil {
				zap.L().Debug("Could not getCurrentClusterVersion", zap.Error(err))
			} else {
				zap.L().Debug("Versions still match",
					zap.String("current", versionCurrent.String()),
					zap.String("init", versionInit.String()))
			}
			time.Sleep(1 * time.Second)
		}

		s.Stop()
	}

	return nil

}

func doCheck(ctx context.Context, domain string) error {

	conn, err := client.GetGRPCClientConn(ctx, domain)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := corev1.NewMainServiceClient(conn)

	latestVersion, err := getLatestVersion(ctx)
	if err != nil {
		return err
	}

	zap.L().Debug("Latest release", zap.String("version", latestVersion.String()))

	currentVersion, err := getCurrentClusterVersion(ctx, c)
	if err != nil {
		return errors.Errorf("Could not parse current Cluster version. Possibly not a production semVer release")
	}

	zap.L().Debug("Current release", zap.String("version", currentVersion.String()))

	if latestVersion.LessThanOrEqual(currentVersion) {
		cliutils.LineNotify("No Cluster upgrade is needed.\n")
		return nil
	}

	cliutils.LineNotify("Cluster can be upgraded\n")
	cliutils.LineNotify("Current Cluster Version: %s.\n", currentVersion.String())
	cliutils.LineNotify("Latest Cluster Version: %s.\n", latestVersion.String())

	return nil
}

func doCheckClient(ctx context.Context) error {

	latestVersion, err := getLatestVersion(ctx)
	if err != nil {
		return err
	}

	currentVersion, err := version.NewSemver(ldflags.SemVer)
	if err != nil {
		return err
	}

	if latestVersion.LessThanOrEqual(currentVersion) {
		cliutils.LineNotify("Your client version is up-to-date.\n")
		return nil
	}

	cliutils.LineNotify("Current Client Version: %s.\n", currentVersion.String())
	cliutils.LineNotify("Latest Client Version: %s.\n", latestVersion.String())
	cliutils.LineNotify("Your Octelium CLIs can be upgraded using the following command:\n")

	switch runtime.GOOS {
	case "linux", "darwin":
		cliutils.LineNotify("curl -fsSL https://octelium.com/install.sh | bash\n")
	case "windows":
		cliutils.LineNotify("iwr https://octelium.com/install.ps1 -useb | iex\n")
	default:
		return errors.Errorf("Unsupported OS platform")
	}

	return nil
}

func getLatestVersion(ctx context.Context) (*version.Version, error) {
	resp, err := resty.New().SetDebug(ldflags.IsDev()).
		R().
		SetContext(ctx).
		Get("https://raw.githubusercontent.com/octelium/octelium/refs/heads/main/unsorted/latest_release")
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, errors.Errorf("Could not get latest Octelium version release")
	}

	return version.NewSemver(string(resp.Body()))
}

func getCurrentClusterVersion(ctx context.Context, c corev1.MainServiceClient) (*version.Version, error) {
	rgn, err := c.GetRegion(ctx, &metav1.GetOptions{
		Name: "default",
	})
	if err != nil {
		return nil, err
	}

	return version.NewSemver(rgn.Status.Version)
}
