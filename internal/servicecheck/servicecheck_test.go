package servicecheck

import (
	"context"
	"testing"
	"time"

	"github.com/VictoriaMetrics/metrics"
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

	checker, err := New(fakeClient, false, 3*time.Second, func(s string) Histogram {
		return metrics.GetOrCreatePrometheusHistogram(s)
	})
	checker.SkipCheckAPIServerDNS = true
	checker.SkipCheckAPIServerDirect = true

	checker.ExtraChecks = map[string]string{
		"check_number_two": "http://interesting.endpoint:8080/abcd",
		"google":           "http://google.ch/",
	}
	r.NoError(err)
	r.NotNil(checker)

	t.Run("run", func(t *testing.T) {
		r := require.New(t)
		checker.Run(context.Background())

		r.Equal(okStr, checker.LastCheckResult[NeighbourhoodState])
		r.Equal(okStr, checker.LastCheckResult["google"])
		r.Equal(errStr, checker.LastCheckResult["check_number_two"])
	})

	// var bb bytes.Buffer
	// metrics.WritePrometheus(&bb, false)
	// fmt.Println(bb.String())
}
