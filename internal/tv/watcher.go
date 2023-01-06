package tv

import "context"

type Watcher interface {
	StartWatching(ctx context.Context, errorChan chan error, responseChan chan Alarm)
}
