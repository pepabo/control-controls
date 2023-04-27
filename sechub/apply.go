package sechub

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	"github.com/rs/zerolog/log"
)

const noteUpdateBy = "control-controls"

func (sh *SecHub) Apply(ctx context.Context, cfg aws.Config, reason string) error {
	region := cfg.Region
	c := securityhub.NewFromConfig(cfg)
	d := sh.Regions.findByRegionName(region)
	a, err := contextcopy(sh)
	if err != nil {
		return err
	}
	if d != nil {
		a, err = Override(sh, d)
		if err != nil {
			return err
		}
	}
	current := New(region)
	if err := current.Fetch(ctx, cfg); err != nil {
		return err
	}
	if !current.enabled {
		log.Info().Str("Region", region).Msg("Skip because Security Hub is not enabled")
		return nil
	}
	diff, err := Diff(current, a)
	if err != nil {
		return err
	}
	if diff == nil {
		log.Info().Str("Region", region).Msg("No changes")
		return nil
	}
	update := false

	// AutoEnable
	if diff.AutoEnable != nil {
		if *diff.AutoEnable {
			log.Info().Str("Region", region).Msg("Enable auto-enable-controls")
		} else {
			log.Info().Str("Region", region).Msg("Disable auto-enable-controls")
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
		a, err := arn.Parse(*s.subscriptionArn)
		if err != nil {
			return err
		}
		if s == nil {
			return fmt.Errorf("could not find standard on %s: %s", region, key)
		}

		// Standards.Enable
		if std.Enable != nil {
			update = true
			switch *std.Enable {
			case true:
				log.Info().Str("Region", region).Str("Standard", key).Msg("Enable standard")
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
				log.Info().Str("Region", region).Str("Standard", key).Msg("Disable standard")
				if _, err := c.BatchDisableStandards(ctx, &securityhub.BatchDisableStandardsInput{
					StandardsSubscriptionArns: []string{*s.subscriptionArn},
				}); err != nil {
					return err
				}
				continue
			}
		}

		if std.Controls == nil && std.Findings == nil {
			log.Debug().Str("Region", region).Str("Standard", key).Msg("Skip controls as there is no difference")
			continue
		}

		// Standards.Controls
		if std.Controls != nil {
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
				log.Info().Str("Region", region).Str("Standard", key).Str("Control", id).Msg("Enable control")
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
			for _, d := range std.Controls.Disable {
				id := d.Key.(string)
				if d.Value.(string) != "" {
					reason = d.Value.(string)
				}
				arn, ok := cs.arns[id]
				if !ok {
					log.Debug().Str("Region", region).Str("Standard", key).Str("Control", id).Msg("Skip control")
					continue
				}
				if containsMapSlice(cs.Disable, id, reason) {
					continue
				}
				update = true
				log.Info().Str("Region", region).Str("Standard", key).Str("Control", id).Str("Reason", reason).Msg("Disable control")
				if _, err := c.UpdateStandardsControl(ctx, &securityhub.UpdateStandardsControlInput{
					StandardsControlArn: arn,
					ControlStatus:       types.ControlStatusDisabled,
					DisabledReason:      &reason,
				}); err != nil {
					return err
				}
				// ref: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/securityhub@v1.20.0#pkg-overview
				// * UpdateStandardsControl - RateLimit of 1 request per second, BurstLimit of 5 requests per second.
				time.Sleep(1 * time.Second)
			}
		}

		// ControlFindingGenerator
		hub, err := c.DescribeHub(ctx, &securityhub.DescribeHubInput{})
		if err != nil {
			return err
		}
		ctrlfg := string(hub.ControlFindingGenerator)

		// Standards.Findings
		if std.Findings != nil {
			cs, err := ctrls(ctx, c, s.subscriptionArn)
			if err != nil {
				return err
			}
			for _, fg := range std.Findings {
				for _, r := range fg.Resources {
					aa, err := arn.Parse(r.Arn)
					if err != nil {
						return err
					}
					if region != "" && aa.Region != "" && aa.Region != region {
						continue
					}
					cArn, ok := cs.arns[fg.ControlID]
					if !ok {
						return fmt.Errorf("not found: %s", fg.ControlID)
					}

					findingFilters := &types.AwsSecurityFindingFilters{
						AwsAccountId: []types.StringFilter{types.StringFilter{Comparison: types.StringFilterComparisonEquals, Value: aws.String(a.AccountID)}},
						ResourceId:   []types.StringFilter{types.StringFilter{Comparison: types.StringFilterComparisonEquals, Value: aws.String(r.Arn)}},
						ProductName:  []types.StringFilter{types.StringFilter{Comparison: types.StringFilterComparisonEquals, Value: aws.String("Security Hub")}},
						RecordState:  []types.StringFilter{types.StringFilter{Comparison: types.StringFilterComparisonEquals, Value: aws.String("ACTIVE")}},
					}
					if ctrlfg == "SECURITY_CONTROL" {
						findingFilters.ComplianceSecurityControlId = []types.StringFilter{types.StringFilter{Comparison: types.StringFilterComparisonEquals, Value: aws.String(fg.ControlID)}}
						findingFilters.ComplianceAssociatedStandardsId = []types.StringFilter{types.StringFilter{Comparison: types.StringFilterComparisonEquals, Value: aws.String(fmt.Sprintf("standards/%s", key))}}
					} else {
						findingFilters.ProductFields = []types.MapFilter{types.MapFilter{Comparison: types.MapFilterComparisonEquals, Key: aws.String("StandardsControlArn"), Value: cArn}}
					}
					got, err := c.GetFindings(ctx, &securityhub.GetFindingsInput{
						Filters: findingFilters,
					})
					if err != nil {
						return err
					}
					if len(got.Findings) != 1 {
						if len(got.Findings) == 0 && aa.Region == "" {
							// eg. arn:aws:s3:::
							continue
						}
						return fmt.Errorf("not found: %s", r.Arn)
					}
					gotFg := got.Findings[0]
					status := string(gotFg.Workflow.Status)
					note := ""
					if gotFg.Note != nil && gotFg.Note.Text != nil {
						note = *gotFg.Note.Text
					}
					if r.Status != status || r.Note != note {
						log.Info().Str("Region", region).Str("Standard", key).Str("Control", fg.ControlID).Str("Resource ID", r.Arn).Str("Status", r.Status).Str("Note", r.Note).Msg("Change workfow status")

						input := &securityhub.BatchUpdateFindingsInput{
							FindingIdentifiers: []types.AwsSecurityFindingIdentifier{{
								Id:         gotFg.Id,
								ProductArn: gotFg.ProductArn,
							}},
							Workflow: &types.WorkflowUpdate{
								Status: types.WorkflowStatus(r.Status),
							},
						}
						if r.Note != "" {
							input.Note = &types.NoteUpdate{
								Text:      aws.String(r.Note),
								UpdatedBy: aws.String(noteUpdateBy),
							}
						}
						if _, err := c.BatchUpdateFindings(ctx, input); err != nil {
							return err
						}
					}
				}
			}
		}

	}

	if !update {
		log.Info().Str("Region", region).Msg("No changes")
	}

	return nil
}
