package sechub

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIntersectFindingGroups(t *testing.T) {
	tests := []struct {
		a    FindingGroups
		b    FindingGroups
		want FindingGroups
	}{
		{nil, nil, FindingGroups{}},
		{
			FindingGroups{
				&FindingGroup{
					ControlID: "IAM.1",
					Resources: FindingResources{
						&FindingResource{
							Arn:    "arn:aws:iam::1234567890:user/user-a",
							Status: "SUPPRESSED",
							Note:   "This is suppressed",
						},
						&FindingResource{
							Arn:    "arn:aws:iam::1234567890:user/user-b",
							Status: "RESOLVED",
							Note:   "This is resolved",
						},
					},
				},
			},
			FindingGroups{
				&FindingGroup{
					ControlID: "IAM.1",
					Resources: FindingResources{
						&FindingResource{
							Arn:    "arn:aws:iam::1234567890:user/user-a",
							Status: "SUPPRESSED",
							Note:   "This is suppressed",
						},
					},
				},
				&FindingGroup{
					ControlID: "IAM.2",
					Resources: FindingResources{
						&FindingResource{
							Arn:    "arn:aws:iam::1234567890:user/user-a",
							Status: "SUPPRESSED",
							Note:   "This is suppressed",
						},
					},
				},
			},
			FindingGroups{
				&FindingGroup{
					ControlID: "IAM.1",
					Resources: FindingResources{
						&FindingResource{
							Arn:    "arn:aws:iam::1234567890:user/user-a",
							Status: "SUPPRESSED",
							Note:   "This is suppressed",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		got := intersectFindingGroups(tt.a, tt.b)
		if diff := cmp.Diff(got, tt.want, nil); diff != "" {
			t.Errorf("%s", diff)
		}
	}
}

func TestDiffFindingGroups(t *testing.T) {
	tests := []struct {
		base FindingGroups
		a    FindingGroups
		want FindingGroups
	}{
		{nil, nil, FindingGroups{}},
		{
			FindingGroups{
				&FindingGroup{
					ControlID: "IAM.1",
					Resources: FindingResources{
						&FindingResource{
							Arn:    "arn:aws:iam::1234567890:user/user-a",
							Status: "SUPPRESSED",
							Note:   "This is suppressed",
						},
						&FindingResource{
							Arn:    "arn:aws:iam::1234567890:user/user-b",
							Status: "RESOLVED",
							Note:   "This is resolved",
						},
					},
				},
			},
			FindingGroups{
				&FindingGroup{
					ControlID: "IAM.1",
					Resources: FindingResources{
						&FindingResource{
							Arn:    "arn:aws:iam::1234567890:user/user-a",
							Status: "SUPPRESSED",
							Note:   "This is suppressed",
						},
					},
				},
				&FindingGroup{
					ControlID: "IAM.2",
					Resources: FindingResources{
						&FindingResource{
							Arn:    "arn:aws:iam::1234567890:user/user-a",
							Status: "SUPPRESSED",
							Note:   "This is suppressed",
						},
					},
				},
			},
			FindingGroups{
				&FindingGroup{
					ControlID: "IAM.2",
					Resources: FindingResources{
						&FindingResource{
							Arn:    "arn:aws:iam::1234567890:user/user-a",
							Status: "SUPPRESSED",
							Note:   "This is suppressed",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		got := diffFindingGroups(tt.base, tt.a)
		if diff := cmp.Diff(got, tt.want, nil); diff != "" {
			t.Errorf("%s", diff)
		}
	}
}
