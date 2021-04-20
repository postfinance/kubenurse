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
	k8s                kubernetes.Interface
	nodeCache          *nodeCache
	allowUnschedulable bool
}

// NodeSchedulability determines if the kubernetes node is in schedulable mode
// or not.
type NodeSchedulability string

// Constants to define the NodeSchedulability of kubernetes Nodes
const (
	NodeSchedulabilityUnknown NodeSchedulability = "Unknown"
	NodeSchedulable           NodeSchedulability = "Schedulable"
	NodeUnschedulable         NodeSchedulability = "Unschedulable"
)

// Neighbour represents a kubenurse which should be reachable
type Neighbour struct {
	PodName         string
	PodIP           string
	HostIP          string
	NodeName        string
	NodeSchedulable NodeSchedulability
	Phase           string // Pod Phase
}

// New creates a new kubediscovery client. The context is used to stop the k8s watchers/informers.
// When allowUnschedulable is true, no node watcher is created and kubenurses
// on unschedulable nodes are considered as neighbours.
func New(ctx context.Context, allowUnschedulable bool) (*Client, error) {
	// create in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("creating in-cluster configuration: %w", err)
	}

	cliset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating clientset: %w", err)
	}

	var nc *nodeCache

	// Watch nodes only if we do not consider kubenurses on unschedulable nodes
	if !allowUnschedulable {
		nc, err = watchNodes(ctx, cliset)
		if err != nil {
			return nil, fmt.Errorf("starting node watcher: %w", err)
		}
	}

	return &Client{
		k8s:                cliset,
		nodeCache:          nc,
		allowUnschedulable: allowUnschedulable,
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

	var neighbours = make([]Neighbour, len(pods.Items))

	// process pods
	for idx := range pods.Items {
		pod := pods.Items[idx]

		// If we allow unschedulable kubenurses, we set the schedulability
		// to unknown in order not to have to set up a node watcher.
		sched := NodeSchedulabilityUnknown
		if !c.allowUnschedulable {
			sched = NodeUnschedulable
			if c.nodeCache.isSchedulable(pod.Spec.NodeName) {
				sched = NodeSchedulable
			}
		}

		n := Neighbour{
			PodName:         pod.Name,
			PodIP:           pod.Status.PodIP,
			HostIP:          pod.Status.HostIP,
			Phase:           string(pod.Status.Phase),
			NodeName:        pod.Spec.NodeName,
			NodeSchedulable: sched,
		}
		neighbours[idx] = n
	}

	return neighbours, nil
}
