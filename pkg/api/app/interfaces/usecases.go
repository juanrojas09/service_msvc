package interfaces

import "context"

type (
	UseCases interface {
		Handle(ctx context.Context, params ...interface{}) (interface{}, error)
	}
)
