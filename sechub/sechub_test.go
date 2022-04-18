package sechub

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
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
							Disable: []string{},
						},
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
