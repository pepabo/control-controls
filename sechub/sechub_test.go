package sechub

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/goccy/go-yaml"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestIntersect(t *testing.T) {
	tests := []struct {
		a    *SecHub
		b    *SecHub
		want *SecHub
	}{
		{
			&SecHub{},
			&SecHub{},
			&SecHub{
				AutoEnable: aws.Bool(true),
				Standards:  Standards{},
			},
		},
		{
			&SecHub{AutoEnable: aws.Bool(false)},
			&SecHub{AutoEnable: aws.Bool(true)},
			&SecHub{
				AutoEnable: aws.Bool(true),
				Standards:  Standards{},
			},
		},
		{
			&SecHub{AutoEnable: aws.Bool(false)},
			&SecHub{AutoEnable: aws.Bool(false)},
			&SecHub{
				AutoEnable: aws.Bool(false),
				Standards:  Standards{},
			},
		},
		{
			&SecHub{Standards: Standards{
				&Standard{
					Key:              "aws-foundational-security-best-practices/v/1.0.0",
					Enable:           aws.Bool(true),
					enabledByDefault: true,
				},
			}},
			&SecHub{Standards: Standards{
				&Standard{
					Key:              "aws-foundational-security-best-practices/v/1.0.0",
					Enable:           aws.Bool(false),
					enabledByDefault: true,
				},
			}},
			&SecHub{
				AutoEnable: aws.Bool(true),
				Standards: Standards{
					&Standard{
						Key:              "aws-foundational-security-best-practices/v/1.0.0",
						Enable:           aws.Bool(true),
						enabledByDefault: true,
						Findings:         FindingGroups{},
					},
				},
			},
		},
		{
			&SecHub{Standards: Standards{
				&Standard{
					Key:              "aws-foundational-security-best-practices/v/1.0.0",
					enabledByDefault: true,
					Controls: &Controls{
						Enable: []string{"IAM.1", "IAM.2"},
					},
				},
			}},
			&SecHub{Standards: Standards{
				&Standard{
					Key:              "aws-foundational-security-best-practices/v/1.0.0",
					enabledByDefault: true,
					Controls: &Controls{
						Enable: []string{"IAM.1", "IAM.3"},
					},
				},
			}},
			&SecHub{
				AutoEnable: aws.Bool(true),
				Standards: Standards{
					&Standard{
						Key:              "aws-foundational-security-best-practices/v/1.0.0",
						Enable:           aws.Bool(true),
						enabledByDefault: true,
						Controls: &Controls{
							Enable:  []string{"IAM.1"},
							Disable: yaml.MapSlice{},
						},
						Findings: FindingGroups{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		got := Intersect(tt.a, tt.b)
		opt := cmpopts.IgnoreUnexported(SecHub{}, Standard{}, Controls{})
		if diff := cmp.Diff(got, tt.want, opt); diff != "" {
			t.Errorf("%s", diff)
		}
	}
}

func TestOverlay(t *testing.T) {
	tests := []struct {
		base    *SecHub
		overlay *SecHub
		want    *SecHub
	}{
		{
			&SecHub{},
			&SecHub{},
			&SecHub{},
		},
		{
			&SecHub{AutoEnable: aws.Bool(false)},
			&SecHub{AutoEnable: aws.Bool(true)},
			&SecHub{
				AutoEnable: aws.Bool(true),
			},
		},
		{
			&SecHub{Standards: Standards{
				&Standard{
					Key:    "aws-foundational-security-best-practices/v/1.0.0",
					Enable: aws.Bool(true),
					Controls: &Controls{
						Enable: []string{"IAM.1", "IAM.2"},
					},
				},
			}},
			&SecHub{Standards: Standards{
				&Standard{
					Key:    "aws-foundational-security-best-practices/v/1.0.0",
					Enable: aws.Bool(false),
				},
			}},
			&SecHub{Standards: Standards{
				&Standard{
					Key:    "aws-foundational-security-best-practices/v/1.0.0",
					Enable: aws.Bool(false),
					Controls: &Controls{
						Enable: []string{"IAM.1", "IAM.2"},
					},
				},
			}},
		},
		{
			&SecHub{Standards: Standards{
				&Standard{
					Key: "aws-foundational-security-best-practices/v/1.0.0",
					Controls: &Controls{
						Enable: []string{"IAM.1", "IAM.2"},
						Disable: yaml.MapSlice{
							yaml.MapItem{Key: "Redshift.4", Value: "Redshit is not running."},
							yaml.MapItem{Key: "Redshift.6", Value: "Redshit is not running."},
						},
					},
				},
			}},
			&SecHub{Standards: Standards{
				&Standard{
					Key: "aws-foundational-security-best-practices/v/1.0.0",
					Controls: &Controls{
						Enable: []string{"IAM.1", "IAM.3", "Redshift.6"},
						Disable: yaml.MapSlice{
							yaml.MapItem{Key: "Redshift.7", Value: "Redshit is not running."},
						},
					},
				},
			}},
			&SecHub{Standards: Standards{
				&Standard{
					Key: "aws-foundational-security-best-practices/v/1.0.0",
					Controls: &Controls{
						Enable: []string{"IAM.1", "IAM.2", "IAM.3", "Redshift.6"},
						Disable: yaml.MapSlice{
							yaml.MapItem{Key: "Redshift.4", Value: "Redshit is not running."},
							yaml.MapItem{Key: "Redshift.7", Value: "Redshit is not running."},
						},
					},
				},
			}},
		},
	}

	for _, tt := range tests {
		tt.base.Overlay(tt.overlay)
		opt := cmpopts.IgnoreUnexported(SecHub{}, Standard{}, Controls{})
		if diff := cmp.Diff(tt.base, tt.want, opt); diff != "" {
			t.Errorf("%s", diff)
		}
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		base *SecHub
		a    *SecHub
		want *SecHub
	}{
		{
			&SecHub{Standards: Standards{
				&Standard{
					Key:    "aws-foundational-security-best-practices/v/1.0.0",
					Enable: aws.Bool(true),
					Controls: &Controls{
						Enable: []string{"IAM.1", "IAM.2"},
					},
				},
			}},
			&SecHub{Standards: Standards{
				&Standard{
					Key:    "aws-foundational-security-best-practices/v/1.0.0",
					Enable: aws.Bool(true),
					Controls: &Controls{
						Enable: []string{"IAM.1", "IAM.2", "IAM.3"},
					},
				},
			}},
			&SecHub{Standards: Standards{
				&Standard{
					Key:    "aws-foundational-security-best-practices/v/1.0.0",
					Enable: nil,
					Controls: &Controls{
						Enable: []string{"IAM.3"},
					},
					Findings: FindingGroups{},
				},
			}},
		},
		{
			&SecHub{Standards: Standards{
				&Standard{
					Key:      "aws-foundational-security-best-practices/v/1.0.0",
					Enable:   aws.Bool(false),
					Controls: &Controls{},
				},
				&Standard{
					Key:      "pci-dss/v/3.2.1",
					Enable:   aws.Bool(false),
					Controls: &Controls{},
				},
			}},
			&SecHub{Standards: Standards{
				&Standard{
					Key:      "aws-foundational-security-best-practices/v/1.0.0",
					Enable:   aws.Bool(false),
					Controls: &Controls{},
				},
				&Standard{
					Key:      "pci-dss/v/3.2.1",
					Enable:   aws.Bool(true),
					Controls: &Controls{},
				},
			}},
			&SecHub{Standards: Standards{
				&Standard{
					Key:      "pci-dss/v/3.2.1",
					Enable:   aws.Bool(true),
					Controls: &Controls{},
					Findings: FindingGroups{},
				},
			}},
		},
		{
			&SecHub{Standards: Standards{
				&Standard{
					Key:      "aws-foundational-security-best-practices/v/1.0.0",
					Enable:   aws.Bool(false),
					Controls: &Controls{},
				},
				&Standard{
					Key:      "pci-dss/v/3.2.1",
					Enable:   aws.Bool(false),
					Controls: &Controls{},
				},
			}},
			&SecHub{Standards: Standards{
				&Standard{
					Key:      "aws-foundational-security-best-practices/v/1.0.0",
					Enable:   aws.Bool(false),
					Controls: &Controls{},
				},
				&Standard{
					Key:      "pci-dss/v/3.2.1",
					Enable:   aws.Bool(false),
					Controls: &Controls{},
				},
			}},
			nil,
		},
		{
			&SecHub{Standards: Standards{
				&Standard{
					Key:    "aws-foundational-security-best-practices/v/1.0.0",
					Enable: aws.Bool(true),
					Controls: &Controls{
						Enable: []string{"IAM.1", "IAM.2"},
						Disable: yaml.MapSlice{
							yaml.MapItem{Key: "IAM.3", Value: "some reason"},
						},
					},
				},
				&Standard{
					Key:      "cis-aws-foundations-benchmark/v/1.2.0",
					Enable:   aws.Bool(true),
					Controls: &Controls{},
				},
			}},
			&SecHub{Standards: Standards{
				&Standard{
					Key:    "aws-foundational-security-best-practices/v/1.0.0",
					Enable: aws.Bool(true),
					Controls: &Controls{
						Enable: []string{"IAM.1", "IAM.2"},
						Disable: yaml.MapSlice{
							yaml.MapItem{Key: "IAM.3", Value: "some reason"},
						},
					},
				},
				&Standard{
					Key:    "cis-aws-foundations-benchmark/v/1.2.0",
					Enable: aws.Bool(true),
					Controls: &Controls{
						Disable: yaml.MapSlice{
							yaml.MapItem{Key: "CIS.1.1", Value: "some reason"},
						},
					},
				},
			}},
			&SecHub{Standards: Standards{
				&Standard{
					Key:    "cis-aws-foundations-benchmark/v/1.2.0",
					Enable: nil,
					Controls: &Controls{
						Disable: yaml.MapSlice{
							yaml.MapItem{Key: "CIS.1.1", Value: "some reason"},
						},
					},
					Findings: FindingGroups{},
				},
			}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got, err := Diff(tt.base, tt.a)
			if err != nil {
				t.Error(err)
			}
			opt := cmpopts.IgnoreUnexported(SecHub{}, Standard{}, Controls{})
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}
