package servicecheck

import (
	"container/heap"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//nolint:gochecknoglobals // used during testing
var (
	osHostname  = os.Hostname
	currentNode string
)

// Neighbour represents a kubenurse which should be reachable
type Neighbour struct {
	PodName  string
	PodIP    string
	HostIP   string
	NodeName string
	NodeHash uint64
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

		if pod.Name == hostname { // only query other pods, not the currently running pod
			currentNode = pod.Spec.NodeName
			continue
		}

		n := Neighbour{
			PodName:  pod.Name,
			PodIP:    pod.Status.PodIP,
			HostIP:   pod.Status.HostIP,
			NodeName: pod.Spec.NodeName,
			NodeHash: sha256Uint64(pod.Spec.NodeName),
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
		check := func(ctx context.Context) (string, error) {
			if c.UseTLS {
				return c.doRequest(ctx, "https://"+neighbour.PodIP+":8443/alwayshappy")
			}

			return c.doRequest(ctx, "http://"+neighbour.PodIP+":8080/alwayshappy")
		}

		_, _ = c.measure(check, "path_"+neighbour.NodeName)
	}
}

type Uint64Heap []uint64

func (h Uint64Heap) Len() int           { return len(h) }
func (h Uint64Heap) Less(i, j int) bool { return h[i] > h[j] } // we want a max-heap, therefore the inversed condition
func (h Uint64Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *Uint64Heap) Push(x any) {
	*h = append(*h, x.(uint64))
}

func (h *Uint64Heap) Pop() any {
	n := len(*h)
	x := (*h)[n-1]
	*h = (*h)[0 : n-1]

	return x
}

func (c *Checker) filterNeighbours(nh []*Neighbour) []*Neighbour {
	m := make(map[uint64]*Neighbour, c.NeighbourLimit+1)

	sl := make(Uint64Heap, 0, c.NeighbourLimit+1)
	h := &sl
	currentNodeHash := sha256Uint64(currentNode)

	heap.Init(h)

	for _, n := range nh {
		adjHash := n.NodeHash - currentNodeHash
		m[adjHash] = n

		heap.Push(h, adjHash)

		if len(*h) > c.NeighbourLimit {
			p := heap.Pop(h).(uint64)
			delete(m, p)
		}
	}

	filteredNeighbours := make([]*Neighbour, 0, c.NeighbourLimit)

	for _, n := range m {
		filteredNeighbours = append(filteredNeighbours, n)
	}

	return filteredNeighbours
}

func sha256Uint64(s string) uint64 {
	h := sha256.Sum256([]byte(s))
	return binary.BigEndian.Uint64(h[:8])
}
