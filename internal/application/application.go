package application

import "context"

type ITask interface {
	Name() string
	Run(ctx context.Context) error
}
