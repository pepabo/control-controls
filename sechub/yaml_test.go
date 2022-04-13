package sechub

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/goccy/go-yaml"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		in   *SecHub
		want string
	}{
		{
			&SecHub{
				AutoEnable: aws.Bool(true),
				Standards: []*Standard{
					&Standard{
						Key:    "cis-aws-foundations-benchmark/v/1.2.0",
						Enable: aws.Bool(true),
						Controls: &Controls{
							Enable: []string{"CIS.1.1", "CIS.1.2"},
						},
					},
				},
			},
			`autoEnable: true
standards:
  cis-aws-foundations-benchmark/v/1.2.0:
    enable: true
    controls:
      enable: [CIS.1.1, CIS.1.2]
`},
	}
	for _, tt := range tests {
		b, err := yaml.Marshal(tt.in)
		if err != nil {
			t.Fatal(err)
		}
		got := string(b)
		if got != tt.want {
			t.Errorf("got %v\nwant %v", got, tt.want)
		}

	}
}
