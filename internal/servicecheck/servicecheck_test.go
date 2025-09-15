package servicecheck

import (
	"context"
	"net/http"
	"net/http/httptest"
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/not-found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// fake client, with a dummy neighbour pod
	fakeClient := fake.NewFakeClient(&fakeNeighbourPod)

	checker, err := New(fakeClient, false, 3*time.Second, func(s string) Histogram {
		return metrics.GetOrCreatePrometheusHistogram(s)
	})
	checker.SkipCheckAPIServerDNS = true
	checker.SkipCheckAPIServerDirect = true

	checker.ExtraChecks = map[string]string{
		"check_not_found": server.URL + "/not-found",
		"check_ok":        server.URL + "/ok",
		"check_ipv6":      "https://ipv6.google.com/",
	}
	r.NoError(err)
	r.NotNil(checker)

	t.Run("run", func(t *testing.T) {
		r := require.New(t)
		checker.Run(context.Background())

		r.Equal(okStr, checker.LastCheckResult[NeighbourhoodState])
		r.Equal(okStr, checker.LastCheckResult["check_ok"])
		// r.Equal(okStr, checker.LastCheckResult["check_ipv6"]) // gh-action doesn't support IPv6 yet
		r.Equal("404 Not Found", checker.LastCheckResult["check_not_found"])
	})
}
