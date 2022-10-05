package sechub

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	"github.com/goccy/go-yaml"
	"github.com/k1LoW/expand"
)

const defaultConsoleURL = "https://ap-northeast-1.console.aws.amazon.com/securityhub/home?region=ap-northeast-1#/findings?search=RecordState%3D%255Coperator%255C%253AEQUALS%255C%253AACTIVE%26WorkflowStatus%3D%255Coperator%255C%253AEQUALS%255C%253ANEW%26WorkflowStatus%3D%255Coperator%255C%253AEQUALS%255C%253ANOTIFIED"

var defaultTemplate = map[string]interface{}{
	"blocks": []interface{}{
		map[string]interface{}{
			"type": "header",
			"text": map[string]interface{}{
				"type":  "plain_text",
				"text":  "AWS Security Hub Notification",
				"emoji": true,
			},
		},
		map[string]interface{}{
			"type": "section",
			"text": map[string]interface{}{
				"type":  "plain_text",
				"text":  "Notifying because condition '{{ cond }}' is met.",
				"emoji": true,
			},
		},
		map[string]interface{}{
			"type": "section",
			"fields": []interface{}{
				map[string]interface{}{
					"type": "mrkdwn",
					"text": "*CRITICAL:*\n{{ critical - critical_resolved - critical_suppressed }}",
				},
				map[string]interface{}{
					"type": "mrkdwn",
					"text": "*HIGH:*\n{{ high - high_resolved - high_suppressed }}",
				},
			},
		},
		map[string]interface{}{
			"type": "section",
			"fields": []interface{}{
				map[string]interface{}{
					"type": "mrkdwn",
					"text": "*MEDIUM:*\n{{ medium - medium_resolved - medium_suppressed }}",
				},
				map[string]interface{}{
					"type": "mrkdwn",
					"text": "*LOW:*\n{{ low - low_resolved - low_suppressed }}",
				},
			},
		},
		map[string]interface{}{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": "<{{ consoleURL }}|View findings>",
			},
		},
	},
}

type NotifyFinding struct {
	SeverityLabel  types.SeverityLabel
	WorkflowStatus types.WorkflowStatus
}

func (sh *SecHub) Notify(ctx context.Context, findings []NotifyFinding) error {
	urep := strings.NewReplacer("ap-northeast-1", sh.region)
	env := map[string]interface{}{
		"region":     sh.region,
		"consoleURL": urep.Replace(defaultConsoleURL),
	}
	for _, sl := range types.SeverityLabelCritical.Values() {
		slkey := strings.ToLower(string(sl))
		env[slkey] = 0
		for _, ws := range types.WorkflowStatusNew.Values() {
			wskey := strings.ToLower(string(ws))
			env[wskey] = 0
			key := fmt.Sprintf("%s_%s", slkey, wskey)
			env[key] = 0
		}
	}
	for _, f := range findings {
		slkey := strings.ToLower(string(f.SeverityLabel))
		wskey := strings.ToLower(string(f.WorkflowStatus))
		key := fmt.Sprintf("%s_%s", slkey, wskey)
		env[slkey] = env[slkey].(int) + 1
		env[wskey] = env[wskey].(int) + 1
		env[key] = env[key].(int) + 1
	}
	for _, n := range sh.Notifications {
		if n.Cond == "" {
			return errors.New("no cond")
		}
		if n.WebhookURL == "" {
			return errors.New("no webhookURL")
		}
		env["cond"] = n.Cond
		tf, err := expr.Eval(fmt.Sprintf("(%s) == true", n.Cond), env)
		if err != nil {
			return err
		}
		if !tf.(bool) {
			continue
		}
		if n.Template == nil {
			n.Template = defaultTemplate
		}
		b, err := expandBody(n.Template, env)
		if err != nil {
			return err
		}
		req, err := http.NewRequest(
			http.MethodPost,
			n.WebhookURL,
			bytes.NewBuffer(b),
		)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
	}
	return nil
}

func expandBody(tmpl, env interface{}) ([]byte, error) {
	const (
		delimStart = "{{"
		delimEnd   = "}}"
	)
	b, err := yaml.Marshal(tmpl)
	if err != nil {
		return nil, err
	}
	e, err := expand.ReplaceYAML(string(b), expand.ExprRepFn(delimStart, delimEnd, env), false)
	if err != nil {
		return nil, err
	}
	var ee interface{}
	if err := yaml.Unmarshal([]byte(e), &ee); err != nil {
		return nil, err
	}
	return json.Marshal(ee)
}
