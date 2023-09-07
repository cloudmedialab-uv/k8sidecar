package main

import (
	"context"
	"fmt"
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
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	config, err := kubeconfig.ClientConfig()
	if err != nil {
		fmt.Print(err)
	}

	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		fmt.Print(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	api := resources.NewApi(client)

	informer := informers.NewSharedInformerFactory(clientset, time.Second*30).Filtercontroller().V1().Filters()

	controller := watcher.NewWatcher(
		workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Queue"),

		informer.Informer().GetIndexer(),

		informer.Informer().GetController(),

		api,
	)

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {

			key, err := cache.MetaNamespaceKeyFunc(obj)

			fmt.Println("Add " + key)
			if err == nil {
				controller.AddQueue(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)

			fmt.Println("Delete " + key)
			if err == nil {
				controller.AddQueue(key)
			}
		},
	})

	go informer.Informer().Run(context.Background().Done())

	controller.Run(context.Background())
}
