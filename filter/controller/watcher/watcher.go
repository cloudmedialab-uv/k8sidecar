package watcher

import (
	"context"
	"log"
	"strings"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"filter/controller/resources"
	filterv1 "filter/pkg/apis/filtercontroller/v1"
)

// Watcher watches the custom resources and handles add/delete events.
type Watcher struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
	api      *resources.Api
}

// NewWatcher initializes a new Watcher.
func NewWatcher(queue workqueue.RateLimitingInterface, indexer cache.Indexer, informer cache.Controller, api *resources.Api) *Watcher {
	return &Watcher{
		queue:    queue,
		indexer:  indexer,
		informer: informer,
		api:      api,
	}
}

// AddQueue adds a key to the watcher's queue.
func (c *Watcher) AddQueue(key string) {
	c.queue.Add(key)
}

// Run starts the Watcher's loop to process items from the queue.
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

// processNextItem processes the next item from the queue.
func (c *Watcher) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	// Handle the synchronization of the key.
	err := c.syncHandler(key.(string))
	if err != nil {
		log.Printf("Error syncing key: %s. Error: %v", key, err)
		c.queue.AddRateLimited(key)
	} else {
		c.queue.Forget(key)
	}
	return true
}

// syncHandler handles the synchronization of the key with the actual state.
func (c *Watcher) syncHandler(key string) error {
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		log.Printf("Failed to fetch object with key %s from store: %v", key, err)
		return err
	}

	// If the object doesn't exist, delete the associated resources.
	if !exists {
		name := strings.Split(key, "/")[1]
		err := c.api.DeleteResources(name)
		if err != nil {
			log.Printf("Error deleting resources for key %s: %v", key, err)
		}
		return err
	}

	// If the object exists, create or update the associated resources.
	resource := obj.(*filterv1.Filter)
	err = c.api.CreateResources(resource)
	if err != nil {
		log.Printf("Error creating resources for key %s: %v", key, err)
	}
	return err
}
