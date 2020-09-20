package slack

import (
	"fmt"

	"github.com/slack-go/slack"

	"github.com/tuxofil/p24fetch/config"
	"github.com/tuxofil/p24fetch/schema"
)

type Slack struct {
	// Configuration used to create the instance.
	config config.Config
	// Slack client
	client *slack.Client
}

// Create new Slack interface instance
func New(cfg *config.Config) (*Slack, error) {
	s := &Slack{config: *cfg}
	if token := cfg.SlackToken; token != "" {
		s.client = slack.New(cfg.SlackToken)
	}
	return s, nil
}

// Send Slack notifications for unsorted transactions
func (s *Slack) ReportUnsorted(trans []schema.Transaction) {
	if !s.IsActive() {
		return
	}
	for _, tran := range trans {
		_, _, err := s.client.PostMessage(s.config.SlackChannel,
			slack.MsgOptionText(fmt.Sprintf(
				"Unsorted transaction from `%s`:\n```%s```",
				s.config.MerchantName, tran.String()), false))
		if err != nil {
			s.config.Logf("post to Slack: %s", err)
		}
	}
}

// Send message to a configured channel.
func (s *Slack) Send(message string) error {
	if !s.IsActive() {
		return nil
	}
	_, _, _, err := s.client.SendMessage(s.config.SlackChannel,
		slack.MsgOptionText(message, false))
	return err
}

// Send message to a configured channel.
func (s *Slack) Sendf(format string, args ...interface{}) error {
	return s.Send(fmt.Sprintf(format, args...))
}

func (s *Slack) IsActive() bool {
	return s.client != nil && s.config.SlackChannel != ""
}
