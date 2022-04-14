package sechub

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	"github.com/goccy/go-yaml"
)

type Controls struct {
	Enable  []string `yaml:"enable,flow,omitempty"`
	Disable []string `yaml:"disable,flow,omitempty"`
	arns    map[string]*string
}

type Standard struct {
	Key              string    `yaml:"key,omitempty"`
	Enable           *bool     `yaml:"enable,omitempty"`
	Controls         *Controls `yaml:"controls,omitempty"`
	arn              *string
	subscriptionArn  *string
	enabledByDefault bool
}

type Standards []*Standard

type Regions []*SecHub

type SecHub struct {
	AutoEnable *bool `yaml:"autoEnable,omitempty"`
	Standards  Standards
	Regions    Regions
	region     string
	enabled    bool
}

func New(r string) *SecHub {
	return &SecHub{
		region: r,
	}
}

func Load(p string) (*SecHub, error) {
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	hub := &SecHub{}
	if err := yaml.Unmarshal(b, hub); err != nil {
		return nil, err
	}
	return hub, err
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
		if as.Enable != nil && bs.Enable != nil && *as.Enable == *bs.Enable {
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
	d := New(a.region)
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

func Override(base, a *SecHub) *SecHub {
	o := deepcopy(base)
	// AutoEnable
	if a.AutoEnable != nil {
		o.AutoEnable = a.AutoEnable
	}

	// Standards
	for _, std := range o.Standards {
		as := a.Standards.findByKey(std.Key)
		if as == nil {
			continue
		}
		// Standards.Enable
		if as.Enable != nil {
			std.Enable = as.Enable
		}
		// Standards.Controls
		if as.Controls != nil {
			if len(as.Controls.Enable) > 0 {
				std.Controls.Enable = append(std.Controls.Enable, as.Controls.Enable...)
			}
			if len(as.Controls.Disable) > 0 {
				std.Controls.Disable = append(std.Controls.Disable, as.Controls.Disable...)
			}
		}
	}
	for _, k := range intersect(base.Standards.keys(), a.Standards.keys()) {
		as := a.Standards.findByKey(k)
		if as == nil {
			continue
		}
		base.Standards = append(base.Standards, as)
	}

	return o
}

func (rs Regions) findByRegionName(name string) *SecHub {
	for _, r := range rs {
		if r.region == name {
			return r
		}
	}
	return nil
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

func standards(ctx context.Context, c *securityhub.Client) (Standards, error) {
	stds := Standards{}
	r, err := c.DescribeStandards(ctx, &securityhub.DescribeStandardsInput{})
	if err != nil {
		return nil, err
	}
	for _, s := range r.Standards {
		key := key(*s.StandardsArn)
		stds = append(stds, &Standard{
			Key:              key,
			Enable:           aws.Bool(false),
			arn:              s.StandardsArn,
			enabledByDefault: s.EnabledByDefault,
		})
	}
	enabled, err := c.GetEnabledStandards(ctx, &securityhub.GetEnabledStandardsInput{})
	if err != nil {
		return nil, err
	}
	for _, s := range enabled.StandardsSubscriptions {
		std := stds.findByKey(key(*s.StandardsArn))
		std.Enable = aws.Bool(true)
		std.subscriptionArn = s.StandardsSubscriptionArn
	}

	return stds, nil
}

func ctrls(ctx context.Context, c *securityhub.Client, subscriptionArn *string) (*Controls, error) {
	cs := &Controls{
		arns: map[string]*string{},
	}
	var nt *string
	for {
		ctrls, err := c.DescribeStandardsControls(ctx, &securityhub.DescribeStandardsControlsInput{
			StandardsSubscriptionArn: subscriptionArn,
			NextToken:                nt,
		})
		if err != nil {
			return nil, err
		}
		for _, ctrl := range ctrls.Controls {
			cs.arns[*ctrl.ControlId] = ctrl.StandardsControlArn
			switch ctrl.ControlStatus {
			case types.ControlStatusEnabled:
				cs.Enable = append(cs.Enable, *ctrl.ControlId)
			case types.ControlStatusDisabled:
				cs.Disable = append(cs.Enable, *ctrl.ControlId)
			}
		}
		nt = ctrls.NextToken
		if ctrls.NextToken == nil {
			break
		}
	}
	return cs, nil
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

func deepcopy(in *SecHub) *SecHub {
	in.Regions = nil
	b, _ := yaml.Marshal(in)
	out := &SecHub{}
	_ = yaml.Unmarshal(b, out)
	return out
}

func key(arn string) string {
	splitted := strings.SplitN(arn, "/", 2)
	return splitted[1]
}
