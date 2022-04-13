package sechub

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
)

type Controls struct {
	Enable  []string `yaml:"enable,flow,omitempty"`
	Disable []string `yaml:"disable,flow,omitempty"`
}

type Standard struct {
	Key              string    `yaml:"key,omitempty"`
	Enable           *bool     `yaml:"enable,omitempty"`
	Controls         *Controls `yaml:"controls,omitempty"`
	enabledByDefault bool
}

type Standards []*Standard

type Regions []*SecHub

type SecHub struct {
	Region     string `yaml:"-"`
	AutoEnable *bool  `yaml:"autoEnable,omitempty"`
	Standards  Standards
	Regions    Regions
}

func New(r string) *SecHub {
	return &SecHub{
		Region: r,
	}
}

func (sh *SecHub) Fetch(ctx context.Context, c *securityhub.Client) error {
	stds := Standards{}
	r, err := c.DescribeStandards(ctx, &securityhub.DescribeStandardsInput{})
	if err != nil {
		return err
	}
	for _, s := range r.Standards {
		key := key(*s.StandardsArn)
		stds = append(stds, &Standard{
			Key:              key,
			Enable:           aws.Bool(false),
			enabledByDefault: s.EnabledByDefault,
		})
	}
	hub, err := c.DescribeHub(ctx, &securityhub.DescribeHubInput{})
	if err != nil {
		return err
	}
	sh.AutoEnable = aws.Bool(hub.AutoEnableControls)
	enabled, err := c.GetEnabledStandards(ctx, &securityhub.GetEnabledStandardsInput{})
	if err != nil {
		return err
	}
	for _, s := range enabled.StandardsSubscriptions {
		std := stds.findByKey(key(*s.StandardsArn))
		std.Enable = aws.Bool(true)
		var nt *string
		for {
			ctrls, err := c.DescribeStandardsControls(ctx, &securityhub.DescribeStandardsControlsInput{
				StandardsSubscriptionArn: s.StandardsSubscriptionArn,
				NextToken:                nt,
			})
			if err != nil {
				return err
			}
			if std.Controls == nil && len(ctrls.Controls) > 0 {
				std.Controls = &Controls{}
			}
			for _, ctrl := range ctrls.Controls {
				switch ctrl.ControlStatus {
				case types.ControlStatusEnabled:
					std.Controls.Enable = append(std.Controls.Enable, *ctrl.ControlId)
				case types.ControlStatusDisabled:
					std.Controls.Disable = append(std.Controls.Enable, *ctrl.ControlId)
				}
			}
			nt = ctrls.NextToken
			if ctrls.NextToken == nil {
				break
			}
		}
	}
	sh.Standards = stds

	return nil
}

func Intersect(a, b *SecHub) *SecHub {
	i := &SecHub{}
	// AutoEnable
	if a.AutoEnable != nil && b.AutoEnable != nil && *a.AutoEnable == *b.AutoEnable {
		i.AutoEnable = a.AutoEnable
	} else {
		i.AutoEnable = aws.Bool(true)
	}

	// Standards
	i.Standards = Standards{}
	ikeys := intersect(a.Standards.keys(), b.Standards.keys())
	for _, k := range ikeys {
		is := &Standard{
			Key: k,
		}
		as := a.Standards.findByKey(k)
		bs := b.Standards.findByKey(k)
		// Standards.Enable
		if as.Enable == bs.Enable {
			is.Enable = as.Enable
		} else {
			is.Enable = aws.Bool(as.enabledByDefault)
		}
		// Standards.Controls
		if as.Controls != nil && bs.Controls != nil {
			is.Controls = &Controls{}
			is.Controls.Enable = intersect(as.Controls.Enable, bs.Controls.Enable)
			is.Controls.Disable = intersect(as.Controls.Disable, bs.Controls.Disable)
		}

		i.Standards = append(i.Standards, is)
	}

	return i
}

func Diff(base, a *SecHub) *SecHub {
	d := &SecHub{Region: a.Region}
	// AutoEnable
	if base.AutoEnable != nil && a.AutoEnable != nil && *base.AutoEnable == *a.AutoEnable {
		d.AutoEnable = nil
	} else {
		d.AutoEnable = a.AutoEnable
	}
	// Standards
	d.Standards = Standards{}
	for _, std := range a.Standards {
		bstd := base.Standards.findByKey(std.Key)
		if bstd == nil {
			d.Standards = append(d.Standards, std)
			continue
		}
		dstd := &Standard{Key: std.Key}
		// Standards.Enable
		if bstd.Enable != nil && std.Enable != nil && *bstd.Enable == *std.Enable {
			dstd.Enable = nil
		} else {
			dstd.Enable = std.Enable
		}

		if dstd.Enable == nil && bstd.Enable != nil && *bstd.Enable == false {
			continue
		}
		if dstd.Enable != nil && *dstd.Enable == false {
			d.Standards = append(d.Standards, dstd)
			continue
		}

		// Standards.Controls
		if bstd.Controls != nil && std.Controls != nil {
			dstd.Controls = &Controls{}
			if len(bstd.Controls.Enable) == len(std.Controls.Enable) && len(intersect(bstd.Controls.Enable, std.Controls.Enable)) == len(std.Controls.Enable) {
				dstd.Controls.Enable = nil
			} else {
				dstd.Controls.Enable = diff(bstd.Controls.Enable, std.Controls.Enable)
			}
			if len(bstd.Controls.Disable) == len(std.Controls.Disable) && len(intersect(bstd.Controls.Disable, std.Controls.Disable)) == len(std.Controls.Disable) {
				dstd.Controls.Disable = nil
			} else {
				dstd.Controls.Disable = diff(bstd.Controls.Disable, std.Controls.Disable)
			}
		} else {
			dstd.Controls = std.Controls
		}
		if len(dstd.Controls.Enable) == 0 && len(dstd.Controls.Disable) == 0 {
			continue
		}

		d.Standards = append(d.Standards, dstd)
	}

	if d.AutoEnable == nil && len(d.Standards) == 0 {
		return nil
	}

	return d
}

func (stds Standards) findByKey(key string) *Standard {
	for _, std := range stds {
		if std.Key == key {
			return std
		}
	}
	return nil
}

func (stds Standards) keys() []string {
	keys := []string{}
	for _, std := range stds {
		keys = append(keys, std.Key)
	}
	return keys
}

func intersect(a, b []string) []string {
	i := []string{}
	for _, e := range a {
		if contains(b, e) {
			i = append(i, e)
		}
	}
	return i
}

func diff(base, a []string) []string {
	i := []string{}
	for _, e := range a {
		if !contains(base, e) {
			i = append(i, e)
		}
	}
	return i
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func key(arn string) string {
	splitted := strings.SplitN(arn, "/", 2)
	return splitted[1]
}
