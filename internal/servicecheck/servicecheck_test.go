package servicecheck

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	promRegistry := prometheus.NewRegistry()
	checker, err := New(fakeClient, promRegistry, false, 3*time.Second, prometheus.DefBuckets)
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
		r.Equal(errStr, checker.LastCheckResult["google"]) // test extra endpoint functionality
	})

	var bb bytes.Buffer
	metrics.WritePrometheus(&bb, false)

	fmt.Println(bb.String())

	rw := httptest.NewRecorder()
	promhttp.HandlerFor(promRegistry,promhttp.HandlerOpts{}).ServeHTTP(rw, &http.Request{})
	fmt.Println(rw.Body.String())
}
