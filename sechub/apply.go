package sechub

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	"github.com/rs/zerolog/log"
)

func (sh *SecHub) Apply(ctx context.Context, cfg aws.Config) error {
	region := cfg.Region
	c := securityhub.NewFromConfig(cfg)
	d := sh.Regions.findByRegionName(region)
	a := sh
	if d != nil {
		a = Override(sh, d)
	}
	current := New(region)
	if err := current.Fetch(ctx, cfg); err != nil {
		return err
	}
	if !current.enabled {
		log.Info().Str("Region", region).Msg("Skip because Security Hub is not enabled")
		return nil
	}
	diff := Diff(current, a)
	if diff == nil {
		log.Info().Str("Region", region).Msg("No changes")
		return nil
	}
	update := false

	// AutoEnable
	if diff.AutoEnable != nil {
		if *diff.AutoEnable == true {
			log.Debug().Str("Region", region).Msg("Enable auto-enable-controls")
		} else {
			log.Debug().Str("Region", region).Msg("Disable auto-enable-controls")
		}
		if _, err := c.UpdateSecurityHubConfiguration(ctx, &securityhub.UpdateSecurityHubConfigurationInput{
			AutoEnableControls: *diff.AutoEnable,
		}); err != nil {
			return err
		}
	}

	// Standards
	stds, err := standards(ctx, c)
	if err != nil {
		return err
	}

	for _, std := range diff.Standards {
		key := std.Key
		s := stds.findByKey(key)
		if s == nil {
			return fmt.Errorf("could not find standard on %s: %s", region, key)
		}

		// Standards.Enable
		if std.Enable != nil {
			update = true
			switch *std.Enable {
			case true:
				log.Debug().Str("Region", region).Str("Standard", key).Msg("Enable standard")
				o, err := c.BatchEnableStandards(ctx, &securityhub.BatchEnableStandardsInput{
					StandardsSubscriptionRequests: []types.StandardsSubscriptionRequest{
						types.StandardsSubscriptionRequest{
							StandardsArn: s.arn,
						},
					},
				})
				if err != nil {
					return err
				}
				s.subscriptionArn = o.StandardsSubscriptions[0].StandardsSubscriptionArn

				// ref: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/securityhub@v1.20.0#pkg-overview
				// * BatchEnableStandards - RateLimit of 1 request per second, BurstLimit of 1 request per second.
				time.Sleep(1 * time.Second)
			case false:
				log.Debug().Str("Region", region).Str("Standard", key).Msg("Disable standard")
				if _, err := c.BatchDisableStandards(ctx, &securityhub.BatchDisableStandardsInput{
					StandardsSubscriptionArns: []string{*s.subscriptionArn},
				}); err != nil {
					return err
				}
				continue
			}
		}

		// Standards.Controls
		if std.Controls == nil {
			log.Debug().Str("Region", region).Str("Standard", key).Msg("Skip controls as there is no difference")
			continue
		}
		cs, err := ctrls(ctx, c, s.subscriptionArn)
		if err != nil {
			return err
		}
		for _, id := range std.Controls.Enable {
			arn, ok := cs.arns[id]
			if !ok {
				log.Debug().Str("Region", region).Str("Standard", key).Str("Control", id).Msg("Skip control")
				continue
			}
			if contains(cs.Enable, id) {
				continue
			}
			update = true
			log.Debug().Str("Region", region).Str("Standard", key).Str("Control", id).Msg("Enable control")
			if _, err := c.UpdateStandardsControl(ctx, &securityhub.UpdateStandardsControlInput{
				StandardsControlArn: arn,
				ControlStatus:       types.ControlStatusEnabled,
			}); err != nil {
				return err
			}
			// ref: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/securityhub@v1.20.0#pkg-overview
			// * UpdateStandardsControl - RateLimit of 1 request per second, BurstLimit of 5 requests per second.
			time.Sleep(1 * time.Second)
		}
		for _, id := range std.Controls.Disable {
			arn, ok := cs.arns[id]
			if !ok {
				log.Debug().Str("Region", region).Str("Standard", key).Str("Control", id).Msg("Skip control")
				continue
			}
			if contains(cs.Disable, id) {
				continue
			}
			update = true
			log.Debug().Str("Region", region).Str("Standard", key).Str("Control", id).Msg("Disable control")
			if _, err := c.UpdateStandardsControl(ctx, &securityhub.UpdateStandardsControlInput{
				StandardsControlArn: arn,
				ControlStatus:       types.ControlStatusDisabled,
			}); err != nil {
				return err
			}
			// ref: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/securityhub@v1.20.0#pkg-overview
			// * UpdateStandardsControl - RateLimit of 1 request per second, BurstLimit of 5 requests per second.
			time.Sleep(1 * time.Second)
		}
	}

	if !update {
		log.Info().Str("Region", region).Msg("No changes")
	}

	return nil
}
