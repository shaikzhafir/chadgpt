package conversor

import "context"

type Conversor interface {
	Ask(ctx context.Context, userMessage, sessionID string) (string, string, error)
}
