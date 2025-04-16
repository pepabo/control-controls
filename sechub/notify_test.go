package sechub

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	"github.com/k1LoW/httpstub"
	"github.com/tenntenn/golden"
)

func TestNotify(t *testing.T) {
	tests := []struct {
		name         string
		notification *Notification
		findings     []NotifyFinding
		notify       bool
	}{
		{
			"use default template",
			&Notification{
				If: "true",
			},
			[]NotifyFinding{
				{
					SeverityLabel:  types.SeverityLabelCritical,
					WorkflowStatus: types.WorkflowStatusNew,
				},
			},
			true,
		},
		{
			"cond false",
			&Notification{
				If: "false",
			},
			[]NotifyFinding{
				{
					SeverityLabel:  types.SeverityLabelCritical,
					WorkflowStatus: types.WorkflowStatusNew,
				},
			},
			false,
		},
		{
			"notify critical",
			&Notification{
				If: "critical > 0",
			},
			[]NotifyFinding{
				{
					SeverityLabel:  types.SeverityLabelCritical,
					WorkflowStatus: types.WorkflowStatusNew,
				},
			},
			true,
		},
		{
			"not notify critical",
			&Notification{
				If: "critical > 0",
			},
			[]NotifyFinding{
				{
					SeverityLabel:  types.SeverityLabelHigh,
					WorkflowStatus: types.WorkflowStatusNew,
				},
			},
			false,
		},
		{
			"use custom template",
			&Notification{
				If: "true",
				Template: map[string]interface{}{
					"critical": "CRITICAL: {{ critical }}",
				},
			},
			[]NotifyFinding{
				{
					SeverityLabel:  types.SeverityLabelCritical,
					WorkflowStatus: types.WorkflowStatusNew,
				},
			},
			true,
		},
		{
			"change header",
			&Notification{
				Header: "Notification!!",
				If:     "true",
			},
			[]NotifyFinding{
				{
					SeverityLabel:  types.SeverityLabelCritical,
					WorkflowStatus: types.WorkflowStatusNew,
				},
			},
			true,
		},
		{
			"change message",
			&Notification{
				Message: "Notice!!",
				If:      "true",
			},
			[]NotifyFinding{
				{
					SeverityLabel:  types.SeverityLabelCritical,
					WorkflowStatus: types.WorkflowStatusNew,
				},
			},
			true,
		},
		{
			"dryrun true",
			&Notification{
				If: "true",
			},
			[]NotifyFinding{
				{
					SeverityLabel:  types.SeverityLabelCritical,
					WorkflowStatus: types.WorkflowStatusNew,
				},
			},
			false,
		},
	}
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}
	region := "dummy-ap-1"
	cfg.Region = region
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpstub.NewRouter(t)
			r.Method(http.MethodPost).Header("Content-Type", "application/json").ResponseString(http.StatusOK, ``)
			ts := r.Server()
			t.Cleanup(func() {
				ts.Close()
			})
			tt.notification.WebhookURL = ts.URL
			sh := New(region)
			sh.Notifications = append(sh.Notifications, tt.notification)
			dryrun := tt.name == "dryrun true"
			if err := sh.Notify(ctx, cfg, tt.findings, dryrun); err != nil {
				t.Error(err)
			}
			if len(r.Requests()) == 0 {
				if tt.notify {
					t.Error("want notify")
				}
				return
			}
			got := r.Requests()[0].Body
			t.Cleanup(func() {
				if err := r.Requests()[0].Body.Close(); err != nil {
					t.Error(err)
				}
			})
			key := strings.Replace(tt.name, " ", "_", -1)
			if os.Getenv("UPDATE_GOLDEN") != "" {
				golden.Update(t, "testdata", key, got)
				return
			}
			if diff := golden.Diff(t, "testdata", key, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
