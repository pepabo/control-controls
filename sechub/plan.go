package sechub

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/rs/zerolog/log"
)

type ChangeType string

const (
	ENABLE  ChangeType = "+"
	DISABLE ChangeType = "-"
)

type Change struct {
	Key        string
	ChangeType ChangeType
}

func (sh *SecHub) Plan(ctx context.Context, cfg aws.Config) ([]*Change, error) {
	changes := []*Change{}
	region := cfg.Region
	c := securityhub.NewFromConfig(cfg)
	d := sh.Regions.findByRegionName(region)
	a := sh
	if d != nil {
		a = Override(sh, d)
	}
	current := New(region)
	if err := current.Fetch(ctx, cfg); err != nil {
		return nil, err
	}
	if !current.enabled {
		log.Info().Str("Region", region).Msg("Skip because Security Hub is not enabled")
		return changes, nil
	}
	diff := Diff(current, a)
	if diff == nil {
		return changes, nil
	}

	// AutoEnable
	if diff.AutoEnable != nil {
		if *diff.AutoEnable == true {
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
		if std.Controls == nil {
			continue
		}
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
				Key:        fmt.Sprintf("%s::standards::%s::controls::%s", region, key, id),
				ChangeType: ENABLE,
			})
		}
		for _, id := range std.Controls.Disable {
			_, ok := cs.arns[id]
			if !ok {
				continue
			}
			changes = append(changes, &Change{
				Key:        fmt.Sprintf("%s::standards::%s::controls::%s", region, key, id),
				ChangeType: DISABLE,
			})
		}
	}

	return changes, nil
}
