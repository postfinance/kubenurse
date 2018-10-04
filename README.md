# Kubenurse
kubenurse is a little service that monitors all network connections in a kubernetes
cluster and exports the taken metrics as prometheus endpoint.

## Project state
This project was written in only a few hours without receiving a polish but worked well.
Documentation and polish will come.

## Deployment
TODO

## Configuration
TODO

### SSL
The http client appends the certificate `/var/run/secrets/kubernetes.io/serviceaccount/ca.crt` if found. You
can disable certificate validation with `KUBENURSE_INSECURE=true`.

## Alive Endpoint
The following json will be returned when accessing `http://0.0.0.0:8080/alive`:

```json
{
  "api_server_direct": "ok",
  "api_server_dns": "ok",
  "me_ingress": "ok",
  "me_service": "ok",
  "hostname": "example.com",
  "neighbourhood_state": "ok",
  "neighbourhood" : [neighbours],
  "headers": {http_request_headers}
}
```

if everything is alright it returns status code 200, else an 500.

## Health Checks
The checks are described in the follwing subsections

### api_server_direct
Checks if the `/version` of the Kubernetes API Server is available through
the direct link provided by the kubelet.

### api_server_dns
Checks if the `/version` of the Kubernetes API Server is available through
the Cluster DNS URL `https://kubernetes.default.svc:PORT`.

### me_ingress
Checks if itself is reachable at the `/alwayshappy` endpoint behind the ingress.
The address is provided by the env var `KUBENURSE_INGRESS_URL` which
could look like `https://kubenurse.example.com`

### me_service
Checks if it isself reachable at the `/alwayshappy` endpoint over the kubernetes service.
The address is provided by the env var `KUBENURSE_SERVICE_URL` which
could look like `http://kubenurse.kube-system.default.svc:8080`
