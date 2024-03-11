package servicecheck

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func generateNeighbours(n int) (nh []*Neighbour) {

	nh = make([]*Neighbour, 0, n)

	for i := range n {
		nodeName := fmt.Sprintf("a1-k8s-abcd%03d.domain.tld", i)
		neigh := Neighbour{
			NodeName: nodeName,
			NodeHash: sha256String(nodeName),
		}
		nh = append(nh, &neigh)
	}

	return
}

func TestNodeFiltering(t *testing.T) {

	n := 1_000
	neighbourLimit := 10
	nh := generateNeighbours(n)
	require.NotNil(t, nh)

	checker := Checker{
		NeighbourLimit: neighbourLimit,
	}

	t.Run("all nodes should get NEIGHBOUR_LIMIT checks", func(t *testing.T) {
		counter := make(map[string]int, n)

		for i := range n {
			currentNode = nh[i].NodeName
			filtered := checker.filterNeighbours(nh)
			require.Equal(t, neighbourLimit, len(filtered))

			for _, neigh := range filtered {
				counter[neigh.NodeName]++
			}
		}

		for _, count := range counter {
			require.Equal(t, neighbourLimit, count, "one node didn't receive exactly NEIGHBOUR_LIMIT checks")
		}

	})
}
