package dlock

import "context"

type DistributedLockService interface {
	Acquire(ctx context.Context, name string) (DistributedLock, error)
}

type DistributedLock interface {
	Renew(ctx context.Context) error
	Release(ctx context.Context) error
	TTL(ctx context.Context) (int64, error)
}
