package watcher

import (
	"context"
	"strings"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"filter/src/controller/resources"
	filterv1 "filter/src/pkg/apis/filtercontroller/v1"
)

type Watcher struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
	api      *resources.Api
}

func NewWatcher(queue workqueue.RateLimitingInterface, indexer cache.Indexer, informer cache.Controller, api *resources.Api) *Watcher {
	return &Watcher{
		queue:    queue,
		indexer:  indexer,
		informer: informer,
		api:      api,
	}
}

func (c *Watcher) AddQueue(key string) {
	c.queue.Add(key)
}

func (c *Watcher) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if !c.processNextItem() {
					return
				}
			}
		}
	}()

	<-ctx.Done()
}

func (c *Watcher) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncHandler(key.(string))
	if err == nil {
		c.queue.Forget(key)
	} else {
		c.queue.AddRateLimited(key)
	}

	return true
}

func (c *Watcher) syncHandler(key string) error {
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		return err
	}

	if !exists {
		name := strings.Split(key, "/")[1]
		err = c.api.DeleteResources(name)
		return err
	}

	resource := obj.(*filterv1.Filter)

	err = c.api.CreateResources(resource)

	return err
}
