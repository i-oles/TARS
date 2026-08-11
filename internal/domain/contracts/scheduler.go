package contracts

import "context"

type IScheduler interface {
	Run(ctx context.Context) error
}
