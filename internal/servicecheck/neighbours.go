package servicecheck

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"slices"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var osHostname = os.Hostname //nolint:gochecknoglobals // used during testing

// Neighbour represents a kubenurse which should be reachable
type Neighbour struct {
	PodName  string
	PodIP    string
	HostIP   string
	NodeName string
	NodeHash string
}

// GetNeighbours returns a slice of neighbour kubenurses for the given namespace and labelSelector.
func (c *Checker) GetNeighbours(ctx context.Context, namespace, labelSelector string) ([]*Neighbour, error) {
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

	var neighbours = make([]*Neighbour, 0, len(pods.Items))

	var hostname, _ = osHostname()

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
			NodeHash: sha256String(pod.Spec.NodeName),
		}
		neighbours = append(neighbours, &n)
	}

	return neighbours, nil
}

// checkNeighbours checks the /alwayshappy endpoint from every discovered kubenurse neighbour. Neighbour pods on nodes
// which are not schedulable are excluded from this check to avoid possible false errors.
func (c *Checker) checkNeighbours(nh []*Neighbour) {
	if c.NeighbourLimit > 0 && len(nh) > c.NeighbourLimit {
		nh = c.filterNeighbours(nh)
	}

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

func (c *Checker) filterNeighbours(nh []*Neighbour) []*Neighbour {
	m := make(map[string]*Neighbour, len(nh))
	l := make([]string, 0, len(nh))

	for _, n := range nh {
		m[n.NodeHash] = n
		l = append(l, n.NodeHash)
	}

	slices.Sort(l)

	currentHostName, _ := osHostname()
	hostnameHash := sha256String(currentHostName)

	if m[hostnameHash].NodeName != currentHostName {
		panic("the current hostname hash doesn't match the value in the map")
	}

	idx, _ := slices.BinarySearch(l, hostnameHash)

	filteredNeighbours := make([]*Neighbour, 0, c.NeighbourLimit)

	for i := 0; i < c.NeighbourLimit; i++ {
		hash := l[(idx+i+1)%len(l)]
		filteredNeighbours = append(filteredNeighbours, m[hash])
	}

	return filteredNeighbours
}

func sha256String(s string) string {
	h := sha256.Sum256([]byte(s))
	return string(h[:])
}
