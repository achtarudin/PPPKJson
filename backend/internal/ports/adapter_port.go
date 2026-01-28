package ports

import "context"

type AdapterPort interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	IsReady() bool
	Value() any
}
