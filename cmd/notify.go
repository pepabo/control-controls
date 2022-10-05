/*
Copyright © 2022 GMO Pepabo, inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	"github.com/pepabo/control-controls/sechub"
	"github.com/spf13/cobra"
)

var notifyCmd = &cobra.Command{
	Use:   "notify [CONFIG_FILE]",
	Short: "notify",
	Long:  `notify.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}
		hub, err := sechub.Load(args[0])
		if err != nil {
			return err
		}
		for _, o := range overlays {
			oo, err := sechub.Load(o)
			if err != nil {
				return err
			}
			hub.Overlay(oo)
		}
		if len(hub.Notifications) == 0 {
			cmd.Println("no notifications")
			return nil
		}
		region, err := detectAggregationRegion(ctx, cfg)
		if err != nil {
			return err
		}
		cfg.Region = region
		findings, err := collectActiveFindings(ctx, cfg)
		if err != nil {
			return err
		}
		if err := hub.Notify(ctx, cfg, findings); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(notifyCmd)
	notifyCmd.Flags().StringSliceVarP(&overlays, "overlay", "", []string{}, "patch file or directory for overlaying")
}

func collectActiveFindings(ctx context.Context, cfg aws.Config) ([]sechub.NotifyFinding, error) {
	c := securityhub.NewFromConfig(cfg)
	hub, err := c.DescribeHub(ctx, &securityhub.DescribeHubInput{})
	if err != nil {
		return nil, err
	}
	a, err := arn.Parse(*hub.HubArn)
	if err != nil {
		return nil, err
	}
	var (
		nt       *string
		findings []types.AwsSecurityFinding
	)
	for {
		o, err := c.GetFindings(ctx, &securityhub.GetFindingsInput{
			Filters: &types.AwsSecurityFindingFilters{
				AwsAccountId:  []types.StringFilter{{Comparison: types.StringFilterComparisonEquals, Value: aws.String(a.AccountID)}},
				ProductName:   []types.StringFilter{{Comparison: types.StringFilterComparisonEquals, Value: aws.String("Security Hub")}},
				RecordState:   []types.StringFilter{{Comparison: types.StringFilterComparisonEquals, Value: aws.String("ACTIVE")}},
				SeverityLabel: []types.StringFilter{{Comparison: types.StringFilterComparisonNotEquals, Value: aws.String("INFORMATIONAL")}},
			},
			MaxResults: int32(100),
			NextToken:  nt,
		})
		if err != nil {
			return nil, err
		}
		findings = append(findings, o.Findings...)
		if o.NextToken == nil {
			break
		}
		nt = o.NextToken
	}
	nf := []sechub.NotifyFinding{}
	for _, f := range findings {
		nf = append(nf, sechub.NotifyFinding{
			SeverityLabel:  f.Severity.Label,
			WorkflowStatus: f.Workflow.Status,
		})
	}
	return nf, nil
}
