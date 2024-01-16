package dlock

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type cosmosDistributedLockService struct {
	client    *azcosmos.Client
	container *azcosmos.ContainerClient
	options   *CosmosDistributedLockOptions
}

type CosmosDistributedLockOptions struct {
	Database               string
	Container              string
	LeaseTimeToLiveSeconds int32
	PartitionKeyPath       string
	RenewalIntervalSeconds int32
}

func NewCosmosDistributedLockService(client *azcosmos.Client, configureOptions func(opt *CosmosDistributedLockOptions)) (DistributedLockService, error) {
	options := &CosmosDistributedLockOptions{
		Database:  "dlock",
		Container: "locks",
	}
	configureOptions(options)
	db, err := client.NewDatabase(options.Database)
	if err != nil {
		return nil, err
	}

	_, err = db.Read(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	container, err := client.NewContainer(options.Database, options.Container)
	if err != nil {
		return nil, err
	}

	_, err = container.Read(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &cosmosDistributedLockService{
		client:    client,
		container: container,
		options:   options,
	}, nil
}

func (s *cosmosDistributedLockService) Acquire(ctx context.Context, name string) (DistributedLock, error) {
	renewalCtx, cancelFunc := context.WithCancel(context.Background())
	lock := &cosmosDistributedLock{
		container:     s.container,
		Id:            name,
		TimeToLive:    s.options.LeaseTimeToLiveSeconds,
		partitionKey:  partitionKey{},
		RequestedAt:   time.Now(),
		renewalCtx:    renewalCtx,
		cancelRenewal: cancelFunc,
	}

	if s.options.PartitionKeyPath != "" && s.options.PartitionKeyPath != "id" {
		lock.partitionKey[s.options.PartitionKeyPath] = "lock"
	}

	lockJson, err := json.Marshal(lock)
	if err != nil {
		return nil, err
	}

	response, err := s.container.CreateItem(ctx, azcosmos.NewPartitionKeyString("lock"), lockJson, nil)
	if err != nil {
		return nil, err
	}

	lock.ETag = response.ETag

	lock.startRenewal(s.options.RenewalIntervalSeconds)

	return lock, nil
}

type partitionKey map[string]string

type cosmosDistributedLock struct {
	partitionKey
	container     *azcosmos.ContainerClient
	ctx           context.Context
	renewalCtx    context.Context
	cancelRenewal context.CancelFunc
	Id            string      `json:"id"`
	TimeToLive    int32       `json:"ttl"`
	RequestedAt   time.Time   `json:"requestedAt"`
	LastRenewed   time.Time   `json:"lastRenewed"`
	ETag          azcore.ETag `json:"_etag"`
}

func (l *cosmosDistributedLock) Release(ctx context.Context) error {
	l.renewalCtx.Done()
	panic("implement me")
}

func (l *cosmosDistributedLock) TTL(ctx context.Context) (int64, error) {

}

func (l *cosmosDistributedLock) startRenewal(interval int32) {
	if interval == 0 {
		return
	}

	go func() {
		for {
			select {
			case <-l.ctx.Done():
				return
			case <-l.renewalCtx.Done():
				return
			case <-time.After(time.Duration(interval) * time.Second):
				l.renew()
			}
		}
	}()
}

func (l *cosmosDistributedLock) renew() {
	l.TimeToLive = int32(time.Now().Sub(l.RequestedAt).Seconds())

	l.LastRenewed = time.Now()
	patchOps := azcosmos.PatchOperations{}
	patchOps.AppendReplace("/lastRenewed", l.LastRenewed)

	response, err := l.container.PatchItem(l.renewalCtx, azcosmos.NewPartitionKeyString("lock"), l.Id, patchOps, &azcosmos.ItemOptions{
		IfMatchEtag: &l.ETag,
	})

	if err != nil {
		return
	}

	l.ETag = response.ETag
}
