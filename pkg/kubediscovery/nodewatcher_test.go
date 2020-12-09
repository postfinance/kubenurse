package kubediscovery

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNodeWatcher(t *testing.T) {
	node := &corev1.Node{}
	node.Name = "testnode"
	node.Spec.Unschedulable = false

	r := require.New(t)
	fakeClient := fake.NewSimpleClientset()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start the informer
	nc, err := watchNodes(ctx, fakeClient)
	r.NoError(err)
	r.NotNil(nc)

	_, err = fakeClient.CoreV1().Nodes().Create(ctx, node, metav1.CreateOptions{})
	r.NoError(err)
	r.True(nc.isSchedulable(node.Name), "node is schedulable")

	node.Spec.Unschedulable = true

	_, err = fakeClient.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	r.NoError(err)
	time.Sleep(100 * time.Millisecond) // the informer needs some time...
	r.False(nc.isSchedulable(node.Name), "node is not schedulable")

	err = fakeClient.CoreV1().Nodes().Delete(ctx, node.Name, metav1.DeleteOptions{})
	r.NoError(err)

	r.True(nc.isSchedulable("unknown"), "node not in cache")
}
