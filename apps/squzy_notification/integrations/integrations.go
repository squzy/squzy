package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/squzy/squzy/apps/squzy_notification/config"
	"github.com/squzy/squzy/apps/squzy_notification/database"
	"github.com/squzy/squzy/internal/httptools"
	api "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"net/http"
	"time"
)

type Integrations interface {
	Slack(ctx context.Context, ruleName string, incident *api.Incident, config *database.SlackConfig)
	Webhook(ctx context.Context, ruleName string, incident *api.Incident, config *database.WebHookConfig)
}

type WebhookRequest struct {
	Id        string             `json:"id"`
	RuleId    string             `json:"string"`
	Status    api.IncidentStatus `json:"status"`
	CreatedAt string             `json:"createdAt,omitempty"`
	UpdatedAt string             `json:"updatedAt,omitempty"`
	Link      string             `json:"link,omitempty"`
}

type integrations struct {
	httpTools httptools.HTTPTool
	cfg       config.Config
}

func (i *integrations) Slack(ctx context.Context, ruleName string, incident *api.Incident, config *database.SlackConfig) {
	createdAt, updatedAt := getIncidentTime(incident)
	msg := &slack.WebhookMessage{
		Attachments: []slack.Attachment{
			{
				Blocks: slack.Blocks{
					BlockSet: []slack.Block{
						slack.SectionBlock{
							Type: slack.MBTSection,
							Text: &slack.TextBlockObject{
								Type: slack.MarkdownType,
								Text: fmt.Sprintf("Incident *%s* with new status *%s*", incident.Id, incident.Status.String()),
							},
						},
						slack.DividerBlock{
							Type: slack.MBTDivider,
						},
						slack.SectionBlock{
							Type: slack.MBTSection,
							Text: &slack.TextBlockObject{
								Type: slack.MarkdownType,
								Text: fmt.Sprintf("Rule: *%s* \nCreatedAt: *%s* \n UpdatedAt: *%s* \n", ruleName, createdAt, updatedAt),
							},
						},
						slack.DividerBlock{
							Type: slack.MBTDivider,
						},
						slack.SectionBlock{
							Type: slack.MBTSection,
							Text: &slack.TextBlockObject{
								Type: slack.MarkdownType,
								Text: fmt.Sprintf("<%s|LINK>", fmt.Sprintf("%s/incidents/%s", i.cfg.GetDashboardHost(), incident.Id)),
							},
						},
					},
				},
			},
		},
	}
	_ = slack.PostWebhookContext(ctx, config.Url, msg)
}

func (i *integrations) Webhook(ctx context.Context, ruleName string, incident *api.Incident, config *database.WebHookConfig) {
	createdAt, updatedAt := getIncidentTime(incident)
	webHook := &WebhookRequest{
		Id:        incident.Id,
		Status:    incident.Status,
		CreatedAt: createdAt,
		RuleId:    incident.RuleId,
		UpdatedAt: updatedAt,
		Link:      fmt.Sprintf("%s/incidents/%s", i.cfg.GetDashboardHost(), incident.Id),
	}
	// Skip error because we create that structure
	body, _ := json.Marshal(webHook)
	req, _ := http.NewRequest(http.MethodPost, config.Url, bytes.NewBuffer(body))

	_, _, _ = i.httpTools.SendRequest(req.WithContext(ctx))
}

func getIncidentTime(incident *api.Incident) (createdAt string, updatedAt string) {
	if len(incident.Histories) > 0 {
		createdAtTime := incident.Histories[0].Timestamp.AsTime()
		if incident.Histories[0].Timestamp.CheckValid() == nil {
			createdAt = createdAtTime.Format(time.RFC3339)
		}
		updatedTime := incident.Histories[len(incident.Histories)-1].Timestamp.AsTime()
		if incident.Histories[len(incident.Histories)-1].Timestamp.CheckValid() == nil {
			updatedAt = updatedTime.Format(time.RFC3339)
		}
	}

	return createdAt, updatedAt
}

func New(httpTools httptools.HTTPTool, cfg config.Config) Integrations {
	return &integrations{
		httpTools: httpTools,
		cfg:       cfg,
	}
}
