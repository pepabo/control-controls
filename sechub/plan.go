package sechub

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	"github.com/rs/zerolog/log"
)

type ChangeType string

const (
	ENABLE  ChangeType = "+"
	DISABLE ChangeType = "-"
	CHANGE  ChangeType = "~"
)

type Change struct {
	Key            string
	ChangeType     ChangeType
	DisabledReason string
	Changed        interface{}
}

func (sh *SecHub) Plan(ctx context.Context, cfg aws.Config, reason string) ([]*Change, error) {
	changes := []*Change{}
	region := cfg.Region
	c := securityhub.NewFromConfig(cfg)
	d := sh.Regions.findByRegionName(region)
	a, err := contextcopy(sh)
	if err != nil {
		return nil, err
	}
	if d != nil {
		a, err = Override(sh, d)
		if err != nil {
			return nil, err
		}
	}
	current := New(region)
	if err := current.Fetch(ctx, cfg); err != nil {
		return nil, err
	}
	if !current.enabled {
		log.Info().Str("Region", region).Msg("Skip because Security Hub is not enabled")
		return changes, nil
	}
	diff, err := Diff(current, a)
	if err != nil {
		return nil, err
	}

	if diff == nil {
		return changes, nil
	}

	// AutoEnable
	if diff.AutoEnable != nil {
		if *diff.AutoEnable {
			changes = append(changes, &Change{
				Key:        fmt.Sprintf("%s::%s", region, "auto-enable-controls"),
				ChangeType: ENABLE,
			})
		} else {
			changes = append(changes, &Change{
				Key:        fmt.Sprintf("%s::%s", region, "auto-enable-controls"),
				ChangeType: DISABLE,
			})
		}
	}

	// Standards
	stds, err := standards(ctx, c)
	if err != nil {
		return nil, err
	}

	for _, std := range diff.Standards {
		key := std.Key
		s := stds.findByKey(key)
		a, err := arn.Parse(*s.subscriptionArn)
		if err != nil {
			return nil, err
		}
		if s == nil {
			return nil, fmt.Errorf("could not find standard on %s: %s", region, key)
		}

		// Standards.Enable
		if std.Enable != nil {
			switch *std.Enable {
			case true:
				changes = append(changes, &Change{
					Key:        fmt.Sprintf("%s::standards::%s", region, key),
					ChangeType: ENABLE,
				})
			case false:
				changes = append(changes, &Change{
					Key:        fmt.Sprintf("%s::standards::%s", region, key),
					ChangeType: DISABLE,
				})
			}
			continue
		}

		// Standards.Controls
		if std.Controls != nil {
			cs, err := ctrls(ctx, c, s.subscriptionArn)
			if err != nil {
				return nil, err
			}
			for _, id := range std.Controls.Enable {
				_, ok := cs.arns[id]
				if !ok {
					continue
				}
				changes = append(changes, &Change{
					Key:            fmt.Sprintf("%s::standards::%s::controls::%s", region, key, id),
					ChangeType:     ENABLE,
					DisabledReason: "",
				})
			}
			for _, d := range std.Controls.Disable {
				id := d.Key.(string)
				if d.Value.(string) != "" {
					reason = d.Value.(string)
				}
				_, ok := cs.arns[id]
				if !ok {
					continue
				}
				changes = append(changes, &Change{
					Key:            fmt.Sprintf("%s::standards::%s::controls::%s", region, key, id),
					ChangeType:     DISABLE,
					DisabledReason: reason,
				})
			}
		}

		// Standards.Findings
		if std.Findings != nil {
			cs, err := ctrls(ctx, c, s.subscriptionArn)
			if err != nil {
				return nil, err
			}
			for _, fg := range std.Findings {
				for _, r := range fg.Resources {
					aa, err := arn.Parse(r.Arn)
					if err != nil {
						return nil, err
					}
					if region != "" && aa.Region != "" && aa.Region != region {
						continue
					}
					cArn, ok := cs.arns[fg.ControlID]
					if !ok {
						return nil, fmt.Errorf("not found: %s", fg.ControlID)
					}
					got, err := c.GetFindings(ctx, &securityhub.GetFindingsInput{
						Filters: &types.AwsSecurityFindingFilters{
							AwsAccountId:  []types.StringFilter{types.StringFilter{Comparison: types.StringFilterComparisonEquals, Value: aws.String(a.AccountID)}},
							ResourceId:    []types.StringFilter{types.StringFilter{Comparison: types.StringFilterComparisonEquals, Value: aws.String(r.Arn)}},
							ProductName:   []types.StringFilter{types.StringFilter{Comparison: types.StringFilterComparisonEquals, Value: aws.String("Security Hub")}},
							ProductFields: []types.MapFilter{types.MapFilter{Comparison: types.MapFilterComparisonEquals, Key: aws.String("StandardsControlArn"), Value: cArn}},
						},
					})
					if err != nil {
						return nil, err
					}
					if len(got.Findings) != 1 {
						if len(got.Findings) == 0 && aa.Region == "" {
							// eg. arn:aws:s3:::
							continue
						}
						return nil, fmt.Errorf("not found: %s", r.Arn)
					}
					status := string(got.Findings[0].Workflow.Status)
					note := ""
					if got.Findings[0].Note != nil && got.Findings[0].Note.Text != nil {
						note = *got.Findings[0].Note.Text
					}
					if r.Status != status || r.Note != note {
						changed := fmt.Sprintf("%s -> %s (note: %s)", status, r.Status, r.Note)
						if r.Note == "" {
							changed = fmt.Sprintf("%s -> %s", status, r.Status)
						}
						changes = append(changes, &Change{
							Key:        *cArn,
							ChangeType: CHANGE,
							Changed:    changed,
						})
					}
				}
			}
		}
	}

	return changes, nil
}
