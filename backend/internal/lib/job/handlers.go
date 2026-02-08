package job

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anuragShingare30/go-boilerplate/internal/config"
	"github.com/anuragShingare30/go-boilerplate/internal/lib/email"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

// @dev This handles all the background jobs (for now, just email related)
// @dev Now, to handle different job, we can create new file(like email_tasks.go) and add the functionality and add handler function here similar to (handleWelcomeEmailTask())

var emailClient *email.Client

func (j *JobService) Init(cfg *config.Config, logger *zerolog.Logger) {
	emailClient = email.NewClient(cfg, logger)
}

func (j *JobService) handleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var p WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal welcome email payload: %w", err)
	}

	j.logger.Info().
		Str("type", "welcome").
		Str("to", p.To).
		Msg("Processing welcome email task")

	err := emailClient.SendWelcomeEmail(
		p.To,
		p.FirstName,
	)
	if err != nil {
		j.logger.Error().
			Str("type", "welcome").
			Str("to", p.To).
			Err(err).
			Msg("Failed to send welcome email")
		return err
	}

	j.logger.Info().
		Str("type", "welcome").
		Str("to", p.To).
		Msg("Successfully sent welcome email")
	return nil
}
