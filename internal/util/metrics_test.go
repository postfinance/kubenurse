package util_test

import (
	"crypto/rand"
	"fmt"
	"strings"
	"testing"

	"github.com/postfinance/kubenurse/internal/util"
)

func BenchmarkGenMetrics(b *testing.B) {
	labels := make([]string, 0, 10)
	for range cap(labels) {
		labels = append(labels, rand.Text())
	}

	for b.Loop() {
		util.GenMetricsName("test", labels...)
	}
}

// BenchmarkOptimizedGenMetrics spun out of a discussion on Join() vs
// strings.Builder comparison. the strings.Join() was used in the end for its
// readibility and the <1% overhead.
func BenchmarkOptimizedGenMetrics(b *testing.B) {
	labels := make([]string, 0, 10)
	for range cap(labels) {
		labels = append(labels, rand.Text())
	}

	optimizedGenMetrics := func(name string, kvs ...string) string {
		n := len(kvs)
		var labels string
		b := strings.Builder{}
		b.Grow(30 * n) // let's assume each label/value is 30

		if n > 0 {
			if n%2 != 0 {
				panic("odd number or label tags, cannot construct the metric name")
			}
			for i := 0; i < n; i += 2 {
				b.WriteString(fmt.Sprintf("%s=%q,", kvs[i], kvs[i+1]))
			}
			labels = b.String()
			labels = labels[:len(labels)-1]
		}

		return fmt.Sprintf("%s_%s{%s}", util.MetricsNamespace, name, labels)
	}
	for b.Loop() {
		optimizedGenMetrics("test", labels...)
	}
}
