// Package kubediscovery implements a discovery mechanism to find other k8s resources.
package kubediscovery

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Client provides the kubediscovery client methods.
type Client struct {
	k8s       kubernetes.Interface
	nodeCache *nodeCache
}

// Neighbour represents a kubenurse which should be reachable
type Neighbour struct {
	PodName         string
	PodIP           string
	HostIP          string
	NodeName        string
	NodeSchedulable bool
	Phase           string // Pod Phase
}

// New creates a new kubediscovery client. The context is used to stop the k8s watchers/informers.
func New(ctx context.Context) (*Client, error) {
	// create in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("creating in-cluster configuration: %w", err)
	}

	cliset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating clientset: %w", err)
	}

	nodeCache, err := watchNodes(ctx, cliset)
	if err != nil {
		return nil, fmt.Errorf("starting node watcher: %w", err)
	}

	return &Client{
		k8s:       cliset,
		nodeCache: nodeCache,
	}, nil
}

// GetNeighbours returns a slice of neighbour kubenurses for the given namespace and labelSelector.
func (c *Client) GetNeighbours(ctx context.Context, namespace, labelSelector string) ([]Neighbour, error) {
	// Get all pods
	pods, err := c.k8s.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("list pods: %w", err)
	}

	var neighbours []Neighbour

	// process pods
	for _, pod := range pods.Items {
		n := Neighbour{
			PodName:         pod.Name,
			PodIP:           pod.Status.PodIP,
			HostIP:          pod.Status.HostIP,
			Phase:           string(pod.Status.Phase),
			NodeName:        pod.Spec.NodeName,
			NodeSchedulable: c.nodeCache.isSchedulable(pod.Spec.NodeName),
		}
		neighbours = append(neighbours, n)
	}

	return neighbours, nil
}
