package util

import "fmt"

const MetricsNamespace = "kubenurse"

func GenMetricsName(name string, kvs ...string) string {
	n := len(kvs)
	labels := ""
	if n > 0 {
		if n%2 != 0 {
			panic("odd number or label tags, cannot construct the metric name")
		}
		for i := 0; i < n; i += 2 {
			labels += fmt.Sprintf("%s=%q,", kvs[i], kvs[i+1])
		}
		labels = labels[:len(labels)-1]
	}

	return fmt.Sprintf("%s_%s{%s}", MetricsNamespace, name, labels)
}
