package kubediscovery

import (
	"context"
	"fmt"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const (
	resyncPeriod = time.Hour * 1
)

type nodeCache struct {
	nodes map[string]bool
	mu    *sync.RWMutex
}

// watchNodes starts an informer to watch v1.Node resource, the context can be used to stop the informer
func watchNodes(ctx context.Context, client kubernetes.Interface) (*nodeCache, error) {
	nc := nodeCache{
		nodes: make(map[string]bool),
		mu:    new(sync.RWMutex),
	}

	informer := informers.NewSharedInformerFactory(client, resyncPeriod).Core().V1().Nodes().Informer()

	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    nc.add,
			UpdateFunc: nc.update,
			DeleteFunc: nc.delete,
		},
	)

	go informer.Run(ctx.Done())

	if ok := cache.WaitForCacheSync(ctx.Done(), informer.HasSynced); !ok {
		return nil, fmt.Errorf("watching nodes: initial cache sync not successful")
	}

	return &nc, nil
}

func (nc *nodeCache) add(obj interface{}) {
	node := obj.(*corev1.Node)

	nc.mu.Lock()
	nc.nodes[node.Name] = node.Spec.Unschedulable
	nc.mu.Unlock()
}

func (nc *nodeCache) delete(obj interface{}) {
	node := obj.(*corev1.Node)

	nc.mu.Lock()
	delete(nc.nodes, node.Name)
	nc.mu.Unlock()
}

func (nc *nodeCache) update(_, obj interface{}) {
	node := obj.(*corev1.Node)

	nc.mu.Lock()
	nc.nodes[node.Name] = node.Spec.Unschedulable
	nc.mu.Unlock()
}

func (nc *nodeCache) isSchedulable(node string) bool {
	nc.mu.RLock()
	defer nc.mu.RUnlock()

	return !nc.nodes[node]
}
