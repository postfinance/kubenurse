# Kubenurse
kubenurse is a little service that monitors all network connections in a kubernetes
cluster and exports the taken metrics as prometheus endpoint.

## Deployment
You can get the Docker Image from [Docker Hub](https://hub.docker.com/r/postfinance/kubenurse/).
The [examples](https://github.com/postfinance/kubenurse/tree/master/examples) directory
contains yamls that deploy the kubenurse to the kube-system namespace.

### Helm deployment

You can also deploy kubenurse with helm, the chart can be found in ./helm/kubenurse.
`helm upgrade [RELEASE_NAME] --install -f values.yaml ./`

#### Configuration settings

| Setting                          | Description                                                                                         | Default               |
| ---------------------------------|---------------------------------------------------------------------------------------------------- | --------------------- |
| daemonset.image.repository       | The repository name                                                                                 | postfinance/kubenurse |
| daemonset.image.tag              | The tag/ version of the image                                                                       | v1.4.0                |
| namespace                        | The namespace where kubenurse will be deployed                                                      | kube-system           |
| serviceAccount.name              | The name of the service account which is used                                                       | kubenurse             |
| service.name                     | The name of service which exposes the kubenurse application                                         | 8080-8080             |
| service.port                     | The port number of the service                                                                      | 8080                  |
| ingress.enabled                  | Enable/ Disable the ingress                                                                         | true                  |
| ingress.className                | The classname of the ingress controller (e.g. the nginx ingress controller)                         | nginx                 |
| ingress.url                      | The url of the ingress; e.g. kubenurse.westeurope.cloudapp.azure.com                                | dummy.kubenurse.com   |
| rbac.allow_unschedulable.enabled | Configure a clusterrole and clusterrolebinding if env KUBENURSE_ALLOW_UNSCHEDULABLE is set to false | false                 |



After everything is set up and Prometheus scrapes the kubenurses, you can build
dashboards that show network latencies and errors or use the metrics for alarming.

![Grafana ingress view](doc/grafana_ingress.png "Grafana ingress view")
![Grafana path view](doc/grafana_path.png "Grafana path view")


## Configuration

kubenurse is configured with environment variables:

- `KUBENURSE_INGRESS_URL`: An URL to the kubenurse in order to check the ingress
- `KUBENURSE_SERVICE_URL`: An URL to the kubenurse in order to check the kubernetes service
- `KUBENURSE_INSECURE`: If "true", TLS connections will not validate the certificate
- `KUBENURSE_EXTRA_CA`: Additional CA cert path for TLS connections
- `KUBENURSE_NAMESPACE`: Namespace in which to look for the neighbour kubenurses
- `KUBENURSE_NEIGHBOUR_FILTER`: A label selector to filter neighbour kubenurses
- `KUBENURSE_ALLOW_UNSCHEDULABLE`: If this is `"true"`, path checks to neighbouring kubenurses are only made if they are running on schedulable nodes. This requires get/list/watch access to `api/v1 Node` resources
- `KUBENURSE_USE_TLS`: If this is `"true"`, enable TLS endpoint on port 8443
- `KUBENURSE_CERT_FILE`: Certificate to use with TLS endpoint
- `KUBENURSE_CERT_KEY`: Key to use with TLS endpoint

Following variables are injected to the Pod by Kubernetes and should not be defined manually:

- `KUBERNETES_SERVICE_HOST`: Host to communicate to the kube-apiserver
- `KUBERNETES_SERVICE_PORT`: Port to communicate to the kube-apiserver

The used http client appends the certificate `/var/run/secrets/kubernetes.io/serviceaccount/ca.crt` if found.

## http Endpoints

The kubenurse listens http on port 8080, optionally https on port 8443, and exposes endpoints:

- `/`: Redirects to `/alive`
- `/alive`: Returns a pretty printed JSON with the check results, described below
- `/alwayshappy`: Returns http-200 which is used for testing itself
- `/metrics`: Exposes [prometheus](https://prometheus.io/) metrics

The `/alive` endpoint retuns a JSON like this with status code 200 if everything is alright else 500:

```json
{
  "api_server_direct": "ok",
  "api_server_dns": "ok",
  "me_ingress": "ok",
  "me_service": "ok",
  "hostname": "kubenurse-1234-x2bwx",
  "neighbourhood_state": "ok",
  "neighbourhood": [
   {
    "PodName": "kubenurse-1234-8fh2x",
    "PodIP": "10.10.10.67",
    "HostIP": "10.12.12.66",
    "NodeName": "k8s-66.example.com",
    "Phase": "Running"
   },
   {
    "PodName": "kubenurse-1234-ffjbs",
    "PodIP": "10.10.10.138",
    "HostIP": "10.12.12.89",
    "NodeName": "k8s-89.example.com",
    "Phase": "Running"
   }
  ],
  "headers": {
   "Accept": [
    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
   ],
   "Accept-Encoding": [
    "gzip, deflate, br"
   ],
   ...
  }
}
```


## Health Checks
Every five seconds and on every access of `/alive`, the checks described below are run.
Check results are cached for 3 seconds in order to prevent excessive network traffic.

A little illustration of what communication occures, is here:

![Communication](doc/Communication.png "Communication")

### API Server Direct
Checks the `/version` endpoint of the Kubernetes API Server through
the direct link (`KUBERNETES_SERVICE_HOST`, `KUBERNETES_SERVICE_PORT`).

Metric type: `api_server_direct`

### API Server DNS
Checks the `/version` endpoint of the Kubernetes API Server through
the Cluster DNS URL `https://kubernetes.default.svc:$KUBERNETES_SERVICE_PORT`.
This also verifies a working `kube-dns` deployment.

Metric type: `api_server_dns`

### Me Ingress
Checks if the kubenurse is reachable at the `/alwayshappy` endpoint behind the ingress.
This address is provided by the environment variable `KUBENURSE_INGRESS_URL` that
could look like `https://kubenurse.example.com`.
This also verifies a correct upstream DNS resolution.

Metric type: `me_ingress`

### Me Service
Checks if the kubenurse is reachable at the `/alwayshappy` endpoint through the kubernetes service.
The address is provided by the environment variable `KUBENURSE_SERVICE_URL` that
could look like `http://kubenurse.mynamespace.default.svc:8080`.
This also verifies a working `kube-proxy` setup.

Metric type: `me_service`

### Neighbourhood
Checks if every neighbour kubenurse is reachable at the `/alwayshappy` endpoint.
Neighbours are discovered by querying the kube-apiserver for every Pod in the
`KUBENURSE_NAMESPACE` with label `KUBENURSE_NEIGHBOUR_FILTER`.
The request is done directly to the Pod-IP (port 8080, or 8443 if TLS is enabled) and the metric types contains the prefix
`path_` and the hostname of the kubelet on which the neighbour kubenurse should run.
Only kubenurses on nodes that are schedulable are considered as neighbours,
this can be changed by setting `KUBENURSE_ALLOW_UNSCHEDULABLE="true"`.

Metric type: `path_$KUBELET_HOSTNAME`

## Metrics
All checks create exposed metrics, that can be used to monitor:

- SDN network latencies and errors
- kubelet-to-kubelet network latencies and errors
- pod-to-apiserver communication
- Ingress roundtrip latencies and errors
- Service roundtrip latencies and errors (kube-proxy)
- Major kube-apiserver issues
- kube-dns (or CoreDNS) errors
- External DNS resolution errors (ingress URL resolution)

At `/metrics` you will find these:
- `kubenurse_errors_total`: Kubenurse error counter partitioned by error type
- `kubenurse_request_duration`: Kubenurse request duration partitioned by error type, summary over one minute
