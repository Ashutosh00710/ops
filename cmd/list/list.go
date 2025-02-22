package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/getnoops/ops/pkg/brain"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List active deployments.",
		Long:  "List all deployments that have a status of either `PENDING` or `RUNNING`.",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = MustNewConfig(viper.GetViper())

			return ListActiveDeployments()
		},
	}

	return cmd
}

func ListActiveDeployments() error {
	res, err := brain.Client.ListActiveDeploymentsWithResponse(context.Background())
	if err != nil {
		return err
	}

	var activeDeploymentsResponse []brain.ActiveDeployment
	json.Unmarshal(res.Body, &activeDeploymentsResponse)

	for _, d := range activeDeploymentsResponse {
		fmt.Printf("\n - %s (%s): %s", d.Status, d.EnvironmentName, d.DeploymentId)
	}

	return nil
}
