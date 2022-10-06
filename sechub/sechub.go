package sechub

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	"github.com/goccy/go-yaml"
	"github.com/k1LoW/expand"
)

type Controls struct {
	Enable  []string      `yaml:"enable,flow,omitempty"`
	Disable yaml.MapSlice `yaml:"disable,omitempty"`
	arns    map[string]*string
}

type Standard struct {
	Key      string        `yaml:"key,omitempty"`
	Enable   *bool         `yaml:"enable,omitempty"`
	Controls *Controls     `yaml:"controls,omitempty"`
	Findings FindingGroups `yaml:"-"`

	arn              *string
	subscriptionArn  *string
	enabledByDefault bool
}

type Notification struct {
	Cond       string `yaml:"if"`
	WebhookURL string `yaml:"webhookURL"`
	Template   interface{}
}

type Standards []*Standard

type Regions []*SecHub

type Notifications []*Notification

type SecHub struct {
	AutoEnable    *bool `yaml:"autoEnable,omitempty"`
	Standards     Standards
	Regions       Regions
	Notifications Notifications `yaml:"notifications,omitempty"`
	region        string        // current region
	enabled       bool          // whether Security Hub is enabled in current region
}

func New(r string) *SecHub {
	return &SecHub{
		region: r,
	}
}

func Load(p string) (*SecHub, error) {
	b, err := os.ReadFile(filepath.Clean(p))
	if err != nil {
		return nil, err
	}
	hub := &SecHub{}
	if err := yaml.Unmarshal(expand.ExpandenvYAMLBytes(b), hub); err != nil {
		return nil, err
	}
	if err := hub.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %s: %w", p, err)
	}
	return hub, err
}

func Intersect(a, b *SecHub) *SecHub {
	i := &SecHub{}
	// AutoEnable
	if a.AutoEnable != nil && b.AutoEnable != nil && *a.AutoEnable == *b.AutoEnable {
		i.AutoEnable = a.AutoEnable
	} else {
		// default: true
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
			is.Controls.Disable = intersectMapSlice(as.Controls.Disable, bs.Controls.Disable)
		}

		// Standards.Findngs
		is.Findings = intersectFindingGroups(as.Findings, bs.Findings)

		i.Standards = append(i.Standards, is)
	}

	return i
}

func Diff(base, a *SecHub) (*SecHub, error) {
	b, err := contextcopy(base)
	if err != nil {
		return nil, err
	}
	d := New(a.region)
	// AutoEnable
	if b.AutoEnable != nil && a.AutoEnable != nil && *b.AutoEnable == *a.AutoEnable {
		d.AutoEnable = nil
	} else {
		d.AutoEnable = a.AutoEnable
	}
	// Standards
	d.Standards = Standards{}
	for _, std := range a.Standards {
		bstd := b.Standards.findByKey(std.Key)
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

		if dstd.Enable == nil && bstd.Enable != nil && !*bstd.Enable {
			continue
		}
		if dstd.Enable != nil && !*dstd.Enable {
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
			if len(bstd.Controls.Disable) == len(std.Controls.Disable) && len(intersectMapSlice(bstd.Controls.Disable, std.Controls.Disable)) == len(std.Controls.Disable) {
				dstd.Controls.Disable = nil
			} else {
				dstd.Controls.Disable = diffMapSlice(bstd.Controls.Disable, std.Controls.Disable)
			}
		} else {
			dstd.Controls = std.Controls
		}

		if dstd.Enable == nil && dstd.Controls == nil {
			continue
		}

		// Standards.Findings
		dstd.Findings = diffFindingGroups(bstd.Findings, std.Findings)

		if dstd.Enable == nil && dstd.Controls != nil && len(dstd.Controls.Enable) == 0 && len(dstd.Controls.Disable) == 0 && len(dstd.Findings) == 0 {
			continue
		}

		d.Standards = append(d.Standards, dstd)

	}

	if d.AutoEnable == nil && len(d.Standards) == 0 {
		return nil, nil
	}

	return d, nil
}

func Override(base, a *SecHub) (*SecHub, error) {
	o, err := contextcopy(base)
	if err != nil {
		return nil, err
	}
	o.overlay(a)
	return o, nil
}

func (base *SecHub) Overlay(overlay *SecHub) {
	base.overlay(overlay)

	for _, r := range base.Regions {
		if or := overlay.Regions.findByRegionName(r.region); or != nil {
			r.overlay(or)
		}
	}
	for _, or := range overlay.Regions {
		if r := base.Regions.findByRegionName(or.region); r == nil {
			base.Regions = append(base.Regions, or)
		}
	}
}

func (base *SecHub) overlay(overlay *SecHub) {
	// AutoEnable
	if overlay.AutoEnable != nil {
		base.AutoEnable = overlay.AutoEnable
	}

	// Standards
	for _, std := range base.Standards {
		as := overlay.Standards.findByKey(std.Key)
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
				std.Controls.Enable = unique(append(std.Controls.Enable, as.Controls.Enable...))
			}
			if len(as.Controls.Disable) > 0 {
				// If 'Enable' and 'Disable' contain the same key, 'Enable' has priority.
				disable := yaml.MapSlice{}
				for _, d := range append(std.Controls.Disable, as.Controls.Disable...) {
					if !contains(std.Controls.Enable, d.Key.(string)) {
						disable = append(disable, d)
					}
				}
				std.Controls.Disable = uniqueMapSlice(disable)
			}
		}
		// Standards.Findings
		std.Findings = overlayFindingGroups(std.Findings, as.Findings)
	}
	for _, k := range diff(base.Standards.keys(), overlay.Standards.keys()) {
		as := overlay.Standards.findByKey(k)
		if as == nil {
			continue
		}
		base.Standards = append(overlay.Standards, as)
	}

	// Notifications
	ns := overlay.Notifications
	for _, b := range base.Notifications {
		exist := false
		for _, o := range overlay.Notifications {
			if b.Cond == o.Cond && b.WebhookURL == o.WebhookURL {
				exist = true
			}
		}
		if !exist {
			ns = append(ns, b)
		}
	}
	base.Notifications = ns
}

func (sh *SecHub) Validate() error {
	if err := sh.validate(); err != nil {
		return err
	}
	for _, r := range sh.Regions {
		if err := r.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (sh *SecHub) validate() error {
	for _, std := range sh.Standards {
		disableKeys := []string{}
		m := map[string]struct{}{}
		if std.Controls != nil && len(std.Controls.Disable) > 0 {
			for _, d := range std.Controls.Disable {
				key := d.Key.(string)
				if _, ok := m[key]; ok {
					return fmt.Errorf("duplicate key: disable control %s", key)
				}
				disableKeys = append(disableKeys, key)
				m[key] = struct{}{}
			}
		}
		dup := intersect(std.Controls.Enable, disableKeys)
		if len(dup) > 0 {
			return fmt.Errorf("it exists for both enable contorol and disable contorol: %s", dup)
		}
	}
	return nil
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
				cs.Disable = append(cs.Disable, yaml.MapItem{Key: *ctrl.ControlId, Value: *ctrl.DisabledReason})
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

func intersectMapSlice(a, b yaml.MapSlice) yaml.MapSlice {
	i := yaml.MapSlice{}
	for _, e := range a {
		if containsMapSlice(b, e.Key.(string), e.Value.(string)) {
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

func diffMapSlice(base, a yaml.MapSlice) yaml.MapSlice {
	i := yaml.MapSlice{}
	for _, e := range a {
		if !containsMapSlice(base, e.Key.(string), e.Value.(string)) {
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

func containsMapSlice(s yaml.MapSlice, k, v string) bool {
	for _, ss := range s {
		if k == ss.Key.(string) && v == ss.Value.(string) {
			return true
		}
	}
	return false
}

func unique(in []string) []string {
	u := []string{}
	m := map[string]struct{}{}
	for _, s := range in {
		if _, ok := m[s]; ok {
			continue
		}
		u = append(u, s)
		m[s] = struct{}{}
	}
	return u
}

func uniqueMapSlice(in yaml.MapSlice) yaml.MapSlice {
	keys := []string{}
	m := map[string]yaml.MapItem{}
	for _, s := range in {
		if _, ok := m[s.Key.(string)]; !ok {
			keys = append(keys, s.Key.(string))
		}
		m[s.Key.(string)] = s
	}
	u := yaml.MapSlice{}
	for _, k := range keys {
		u = append(u, m[k])
	}
	return u
}

func contextcopy(in *SecHub) (*SecHub, error) {
	b, err := yaml.Marshal(in)
	if err != nil {
		return nil, err
	}
	out := &SecHub{}
	if err := yaml.UnmarshalWithOptions(b, out, yaml.DisallowDuplicateKey()); err != nil {
		return nil, err
	}
	out.Regions = nil
	return out, nil
}

func key(arn string) string {
	splitted := strings.SplitN(arn, "/", 2)
	return splitted[1]
}
