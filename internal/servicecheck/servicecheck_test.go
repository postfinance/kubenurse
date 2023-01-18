package servicecheck

import (
	"context"
	"testing"
	"time"

	"github.com/postfinance/kubenurse/internal/kubediscovery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var fakeNeighbourPod = v1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name: "kubenurse-dummy",
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
	},
}

func TestCombined(t *testing.T) {
	r := require.New(t)

	// fake client, with a dummy neighbour pod
	fakeClient := fake.NewSimpleClientset()
	_, err := fakeClient.CoreV1().Pods("kube-system").Create(context.Background(), &fakeNeighbourPod, metav1.CreateOptions{})
	r.NoError(err)

	discovery, err := kubediscovery.New(context.Background(), fakeClient, false)
	r.NoError(err)

	checker, err := New(context.Background(), discovery, prometheus.NewRegistry(), false, 3*time.Second, prometheus.DefBuckets)
	r.NoError(err)
	r.NotNil(checker)

	t.Run("run", func(t *testing.T) {
		r := require.New(t)
		result, hadError := checker.Run()
		r.True(hadError)
		r.Len(result.Neighbourhood, 1)
	})

	t.Run("scheduled", func(t *testing.T) {
		stopped := make(chan struct{})

		go func() {
			// blocks until StopScheduled()
			checker.RunScheduled(time.Second * 5)

			close(stopped)
		}()

		checker.StopScheduled()

		<-stopped
	})
}
