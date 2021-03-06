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
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/fatih/color"
	"github.com/pepabo/control-controls/sechub"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan [CONFIG_FILE]",
	Short: "plan",
	Long:  `plan.`,
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
		regions, err := regions(ctx, cfg)
		if err != nil {
			return err
		}

		changes := []*sechub.Change{}
		for _, r := range regions {
			cfg.Region = r
			log.Info().Msg(fmt.Sprintf("Checking %s", r))
			c, err := hub.Plan(ctx, cfg, disabledReason)
			if err != nil {
				return err
			}
			changes = append(changes, c...)
		}

		cmd.Println("")

		if len(changes) == 0 {
			cmd.Println("No changes. Controls are up-to-date.")
		} else {
			green := color.New(color.FgGreen).PrintfFunc()
			red := color.New(color.FgRed).PrintfFunc()
			yellow := color.New(color.FgYellow).PrintfFunc()
			enable := 0
			disable := 0
			change := 0
			for _, c := range changes {
				switch c.ChangeType {
				case sechub.ENABLE:
					enable += 1
					green("%s %s\n", c.ChangeType, c.Key)
				case sechub.DISABLE:
					disable += 1
					if c.DisabledReason == "" {
						red("%s %s\n", c.ChangeType, c.Key)
					} else {
						red("%s %s (disabled reason: %s)\n", c.ChangeType, c.Key, c.DisabledReason)
					}
				case sechub.CHANGE:
					change += 1
					yellow("%s %s %s\n", c.ChangeType, c.Key, c.Changed)
				}
			}
			cmd.Println("")
			cmd.Printf("Plan: %d to enable, %d to change, %d to disable\n", enable, change, disable)
			os.Exit(2)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.Flags().StringVarP(&disabledReason, "disabled-reason", "", "", "A description of the reason why you are disabling a security standard control.")
	planCmd.Flags().StringSliceVarP(&overlays, "overlay", "", []string{}, "patch file or directory for overlaying")
}
