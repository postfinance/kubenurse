package util

import (
	"fmt"
	"strings"
)

const MetricsNamespace = "kubenurse"

func GenMetricsName(name string, kvs ...string) string {
	n := len(kvs)
	labels := make([]string, n/2)

	if n > 0 {
		if n%2 != 0 {
			panic("odd number or label tags, cannot construct the metric name")
		}
		for i := 0; i < n; i += 2 {
			labels[i/2] = fmt.Sprintf("%s=%q", kvs[i], kvs[i+1])
		}
	}

	return fmt.Sprintf("%s_%s{%s}", MetricsNamespace, name, strings.Join(labels, ","))
}
