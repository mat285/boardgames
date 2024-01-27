package persist

import (
	"context"

	"github.com/blend/go-sdk/uuid"
)

type Interface interface {
	CheckAndSet(context.Context, Object) (*Object, error)
	Load(context.Context, uuid.UUID) (*Object, error)
}
