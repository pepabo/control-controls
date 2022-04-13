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
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/goccy/go-yaml"
	"github.com/pepabo/control-controls/sechub"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export",
	Long:  `export.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}
		ec2s := ec2.NewFromConfig(cfg)
		regions, err := ec2s.DescribeRegions(ctx, &ec2.DescribeRegionsInput{AllRegions: aws.Bool(false)})
		if err != nil {
			return err
		}
		var base *sechub.SecHub
		hubs := []*sechub.SecHub{}
		for _, r := range regions.Regions {
			cfg.Region = *r.RegionName
			log.Info().Msg(fmt.Sprintf("Fetching controls from %s", cfg.Region))
			sh := sechub.New(*r.RegionName)
			if err := sh.Fetch(ctx, cfg); err != nil {
				return err
			}
			if base == nil {
				base = sh
			} else {
				base = sechub.Intersect(base, sh)
			}
			hubs = append(hubs, sh)
		}

		for _, h := range hubs {
			d := sechub.Diff(base, h)
			if d != nil {
				base.Regions = append(base.Regions, d)
			}
		}

		b, err := yaml.Marshal(base)
		if err != nil {
			return err
		}

		cmd.Println(string(b))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
