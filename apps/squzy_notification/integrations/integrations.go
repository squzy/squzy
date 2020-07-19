package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/slack-go/slack"
	api "github.com/squzy/squzy_generated/generated/proto/v1"
	"net/http"
	"squzy/apps/squzy_notification/config"
	"squzy/apps/squzy_notification/database"
	"squzy/internal/httptools"
	"time"
)

type Integrations interface {
	Slack(ctx context.Context, incident *api.Incident, config *database.SlackConfig)
	Webhook(ctx context.Context, incident *api.Incident, config *database.WebHookConfig)
}

type WebhookRequest struct {
	Id        string             `json:"id"`
	Status    api.IncidentStatus `json:"status"`
	CreatedAt string             `json:"createdAt,omitempty"`
	UpdatedAt string             `json:"updatedAt,omitempty"`
}

type integrations struct {
	httpTools httptools.HTTPTool
	cfg config.Config
}

func (i *integrations) Slack(ctx context.Context, incident *api.Incident, config *database.SlackConfig) {
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
								Text: fmt.Sprintf("CreatedAt: %s \n UpdatedAt: %s \n", createdAt, updatedAt),
							},
						},
						slack.SectionBlock{
							Type: slack.MBTSection,
							Text: &slack.TextBlockObject{
								Type: slack.MarkdownType,
								Text: fmt.Sprintf("<%s|LINK>", fmt.Sprintf("%s/incidents/%s",i.cfg.GetDashboardHost(), incident.Id)),
							},
						},
					},
				},
			},
		},
	}
	_ = slack.PostWebhook(config.Url, msg)
}

func (i *integrations) Webhook(ctx context.Context, incident *api.Incident, config *database.WebHookConfig) {
	createdAt, updatedAt := getIncidentTime(incident)
	webHook := &WebhookRequest{
		Id:        incident.Id,
		Status:    incident.Status,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	// Skip error because we create that structure
	body, _ := json.Marshal(webHook)
	req, _ := http.NewRequest(http.MethodPost, config.Url, bytes.NewBuffer(body))

	_, _, _ = i.httpTools.SendRequest(req.WithContext(ctx))
}

func getIncidentTime(incident *api.Incident) (createdAt string, updatedAt string) {
	if len(incident.Histories) > 0 {
		createdAtTime, err := ptypes.Timestamp(incident.Histories[0].Timestamp)
		if err == nil {
			createdAt = createdAtTime.Format(time.RFC3339)
		}
		updatedTime, err := ptypes.Timestamp(incident.Histories[len(incident.Histories)-1].Timestamp)
		if err == nil {
			updatedAt = updatedTime.Format(time.RFC3339)
		}
	}

	return createdAt, updatedAt
}

func New(httpTools httptools.HTTPTool, cfg config.Config) Integrations {
	return &integrations{
		httpTools: httpTools,
		cfg:cfg,
	}
}
