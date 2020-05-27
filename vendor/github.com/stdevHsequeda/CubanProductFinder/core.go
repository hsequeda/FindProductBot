package storeClient

import (
	"context"
	"time"
)

type Section interface {
	GetName() string
	GetUrl() string
	GetParent() string
	GetStore() *Store
	GetPriority(string) int
	GetReadyTime() time.Time
	SetReadyTime(readyTime time.Time)
}

type Product interface {
	GetName() string
	GetPrice() string
	GetLink() string
	GetSection() Section
}

type Worker interface {
	GetArgs() context.Context
	Task(ctx context.Context)
}
