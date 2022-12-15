package processing

import (
	"context"
	"sync"

	contracts "github.com/punk-link/platform-contracts"
)

type QueueProcessor interface {
	Process(ctx context.Context, wg *sync.WaitGroup, platformer contracts.Platformer)
}
