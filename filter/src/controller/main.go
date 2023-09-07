package main

import (
	"context"
	"log"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"

	"filter/src/controller/resources"
	watcher "filter/src/controller/watcher"

	clientset "filter/src/pkg/generated/clientset/versioned"
	informers "filter/src/pkg/generated/informers/externalversions"
)

func main() {
	// Set up the kubeconfig for connecting to a Kubernetes cluster.
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Fetch the actual config (kubeconfig)
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %s", err.Error())
	}

	// Initialize a new clientset for our custom resource.
	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %s", err.Error())
	}

	// Initialize the standard Kubernetes client.
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create kubernetes client: %s", err.Error())
	}

	// Create a new API resource manager using the Kubernetes client.
	api := resources.NewApi(client)

	// Set up informers to watch for changes in our custom resource.
	informer := informers.NewSharedInformerFactory(clientset, time.Second*30).Filtercontroller().V1().Filters()

	// Initialize our custom controller (watcher).
	controller := watcher.NewWatcher(
		workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Queue"),
		informer.Informer().GetIndexer(),
		informer.Informer().GetController(),
		api,
	)

	// Set up event handlers for when our custom resource is added or deleted.
	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			log.Printf("Resource added: %s", key)
			if err == nil {
				controller.AddQueue(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			log.Printf("Resource deleted: %s", key)
			if err == nil {
				controller.AddQueue(key)
			}
		},
	})

	// Start the informer to begin watching for changes.
	go informer.Informer().Run(context.Background().Done())

	// Run the custom controller (watcher) to process events.
	controller.Run(context.Background())
}
