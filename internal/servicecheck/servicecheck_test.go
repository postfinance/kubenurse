package servicecheck

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var fakeNeighbourPod = v1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "kubenurse-dummy",
		Namespace: "kube-system",
		Labels: map[string]string{
			"app": "kubenurse",
		},
	},
	Spec: v1.PodSpec{
		NodeName: "dummy",
	},
	Status: v1.PodStatus{
		HostIP: "127.0.0.1",
		PodIP:  "127.0.0.1",
		Phase:  v1.PodRunning,
	},
}

func TestCombined(t *testing.T) {
	r := require.New(t)

	// fake client, with a dummy neighbour pod
	fakeClient := fake.NewFakeClient(&fakeNeighbourPod)

	checker, err := New(fakeClient, prometheus.NewRegistry(), false, 3*time.Second, prometheus.DefBuckets)
	r.NoError(err)
	r.NotNil(checker)

	t.Run("run", func(t *testing.T) {
		r := require.New(t)
		checker.Run(context.Background())

		r.Equal(okStr, checker.LastCheckResult[NeighbourhoodState])
	})
}
