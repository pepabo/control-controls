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

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pepabo/control-controls/sechub"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply [CONFIG_FILE]",
	Short: "apply",
	Long:  `apply.`,
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

		regions, err := regions(ctx, cfg)
		if err != nil {
			return err
		}

		for _, r := range regions {
			cfg.Region = r
			log.Info().Msg(fmt.Sprintf("Applying to %s", r))
			if err := hub.Apply(ctx, cfg); err != nil {
				return err
			}
		}

		log.Info().Msg("Apply complete")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
