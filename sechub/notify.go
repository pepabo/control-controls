package sechub

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/antonmedv/expr"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	"github.com/goccy/go-yaml"
	"github.com/k1LoW/expand"
)

const (
	defaultHeader      = "'*AWS Security Hub Notification*'"
	defaultMessageTmpl = "Notified because condition *'%s'* is met."
	defaultConsoleURL  = "https://ap-northeast-1.console.aws.amazon.com/securityhub/home?region=ap-northeast-1#/findings?search=RecordState%3D%255Coperator%255C%253AEQUALS%255C%253AACTIVE%26WorkflowStatus%3D%255Coperator%255C%253AEQUALS%255C%253ANEW%26WorkflowStatus%3D%255Coperator%255C%253AEQUALS%255C%253ANOTIFIED"
)

var defaultTemplate = map[string]interface{}{
	"blocks": []interface{}{
		map[string]interface{}{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": "{{ header }}",
			},
		},
		map[string]interface{}{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": "{{ message }}",
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

func (sh *SecHub) Notify(ctx context.Context, cfg aws.Config, findings []NotifyFinding, dryrun bool) error {
	urep := strings.NewReplacer("ap-northeast-1", cfg.Region)
	now := time.Now()
	env := map[string]interface{}{
		"region":     sh.region,
		"consoleURL": urep.Replace(defaultConsoleURL),
		"month":      int(now.Month()),
		"day":        now.Day(),
		"hour":       now.Hour(),
		"weekday":    int(now.Weekday()),
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
		if n.Header == "" {
			n.Header = defaultHeader
		}
		if n.Message == "" {
			n.Message = fmt.Sprintf(defaultMessageTmpl, n.If)
		}
		if n.If == "" {
			return errors.New("no cond")
		}
		if !dryrun && n.WebhookURL == "" {
			return errors.New("no webhookURL")
		}
		env["header"] = n.Header
		env["cond"] = n.If
		env["message"] = n.Message
		tf, err := expr.Eval(fmt.Sprintf("(%s) == true", n.If), env)
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

		if dryrun {
			var out bytes.Buffer
			if err := json.Indent(&out, b, "", "  "); err != nil {
				return err
			}
			fmt.Println(out.String())
			continue
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
		if err := resp.Body.Close(); err != nil {
			return err
		}
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
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(ee); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
