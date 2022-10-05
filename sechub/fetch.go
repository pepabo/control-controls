package sechub

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
)

func (sh *SecHub) Fetch(ctx context.Context, cfg aws.Config) error {
	c := securityhub.NewFromConfig(cfg)
	stds, err := standards(ctx, c)
	if err != nil {
		return err
	}
	hub, err := c.DescribeHub(ctx, &securityhub.DescribeHubInput{})
	if err != nil {
		return err
	}
	if hub.SubscribedAt != nil {
		sh.enabled = true
	}
	sh.AutoEnable = aws.Bool(hub.AutoEnableControls)
	for _, std := range stds {
		if std.Enable == nil || !*std.Enable {
			continue
		}
		cs, err := ctrls(ctx, c, std.subscriptionArn)
		if err != nil {
			return err
		}
		std.Controls = cs
	}
	sh.Standards = stds

	return nil
}
