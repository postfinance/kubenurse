package servicecheck

import (
	"context"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Neighbour represents a kubenurse which should be reachable
type Neighbour struct {
	PodName  string
	PodIP    string
	HostIP   string
	NodeName string
}

// GetNeighbours returns a slice of neighbour kubenurses for the given namespace and labelSelector.
func (c *Checker) GetNeighbours(ctx context.Context, namespace, labelSelector string) ([]Neighbour, error) {
	// Get all pods
	pods := v1.PodList{}
	selector, _ := labels.Parse(labelSelector)
	err := c.client.List(ctx, &pods, &client.ListOptions{
		LabelSelector: selector,
		Namespace:     namespace,
	})

	if err != nil {
		return nil, fmt.Errorf("list pods: %w", err)
	}

	var neighbours = make([]Neighbour, 0, len(pods.Items))

	var hostname, _ = os.Hostname()

	// process pods
	for idx := range pods.Items {
		pod := pods.Items[idx]

		if !c.allowUnschedulable { // if we disallow unschedulable nodes, we have to check their status
			n := v1.Node{}
			if err := c.client.Get(ctx, types.NamespacedName{Name: pod.Spec.NodeName}, &n); err == nil {
				if n.Spec.Unschedulable { // node unschedulable, we do not include this pod in the neighbour list
					continue
				}
			}
		}

		if pod.Status.Phase != v1.PodRunning || // only query running pods (excludes pending ones)
			pod.DeletionTimestamp != nil { // exclude terminating pods
			continue
		}

		if pod.Name == hostname { // only quey other pods, not the currently running pod
			continue
		}

		n := Neighbour{
			PodName:  pod.Name,
			PodIP:    pod.Status.PodIP,
			HostIP:   pod.Status.HostIP,
			NodeName: pod.Spec.NodeName,
		}
		neighbours = append(neighbours, n)
	}

	return neighbours, nil
}

// checkNeighbours checks the /alwayshappy endpoint from every discovered kubenurse neighbour. Neighbour pods on nodes
// which are not schedulable are excluded from this check to avoid possible false errors.
func (c *Checker) checkNeighbours(nh []Neighbour) {
	for _, neighbour := range nh {
		neighbour := neighbour // pin

		check := func(ctx context.Context) (string, error) {
			if c.UseTLS {
				return c.doRequest(ctx, "https://"+neighbour.PodIP+":8443/alwayshappy")
			}

			return c.doRequest(ctx, "http://"+neighbour.PodIP+":8080/alwayshappy")
		}

		_, _ = c.measure(check, "path_"+neighbour.NodeName)
	}
}
