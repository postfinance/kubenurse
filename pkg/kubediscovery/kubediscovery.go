package kubediscovery

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getClientset() (*kubernetes.Clientset, error) {
	// create in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// create clientset
	return kubernetes.NewForConfig(config)
}

type Neighbour struct {
	PodName  string
	PodIP    string
	HostIP   string
	NodeName string
	Phase    string
}

func GetNeighbourhood(namespace, labelSelector string) ([]Neighbour, error) {
	clientset, err := getClientset()
	if err != nil {
		return nil, err
	}

	// Get all pods
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: labelSelector,
	})

	if err != nil {
		return nil, err
	}

	// Process pods
	res := []Neighbour{}
	for _, pod := range pods.Items {
		n := Neighbour{
			PodName:  pod.Name,
			PodIP:    pod.Status.PodIP,
			HostIP:   pod.Status.HostIP,
			Phase:    string(pod.Status.Phase),
			NodeName: pod.Spec.NodeName,
		}
		res = append(res, n)
	}

	return res, nil
}
