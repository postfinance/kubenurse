package kubediscovery

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	kubenursePod = v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kubenurse-dummy",
			Labels: map[string]string{
				"app": "kubenurse",
			},
		},
	}
	differentPod = v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "different",
			Labels: map[string]string{
				"app": "different",
			},
		},
	}
)

func TestGetNeighbours(t *testing.T) {
	r := require.New(t)
	fakeClient := fake.NewSimpleClientset()

	createFakePods(fakeClient)

	client, err := New(context.Background(), fakeClient, false)
	r.NoError(err)

	neighbours, err := client.GetNeighbours(context.Background(), "kube-system", "app=kubenurse")
	r.NoError(err)
	r.Len(neighbours, 1)
	r.Equal(kubenursePod.ObjectMeta.Name, neighbours[0].PodName)
}

func createFakePods(k8s kubernetes.Interface) {
	for _, pod := range []v1.Pod{kubenursePod, differentPod} {
		_, err := k8s.CoreV1().Pods("kube-system").Create(context.Background(), &pod, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
	}
}
