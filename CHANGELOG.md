## 1.8.0 (2023-06-26)

### New Features

* **common**: add new helm configurations ([8c2e6c65](https://github.com/postfinance/kubenurse/commit/8c2e6c65))
* **helm**: new configuration options (#57) ([13484e61](https://github.com/postfinance/kubenurse/commit/13484e61))



## 1.7.0 (2022-11-01)


### Bug Fixes

* **common**: Use current main branch naming for the helm releaser ([4dd5eded](https://github.com/postfinance/kubenurse/commit/4dd5eded))
* **common**: use new ingress spefification (#52) ([8b896f4c](https://github.com/postfinance/kubenurse/commit/8b896f4c))
* **helm**: chart should respect `-n <namespace>` flag (#53) ([a5a3a792](https://github.com/postfinance/kubenurse/commit/a5a3a792))
* **helm**: parse error when using extraEnvs (#48) ([3a56edbb](https://github.com/postfinance/kubenurse/commit/3a56edbb))


### New Features

* **common**: Implement helm chart releaser (#47) ([7f52b474](https://github.com/postfinance/kubenurse/commit/7f52b474))
* **helm**: add dnsConfig option (#50) ([3fed2690](https://github.com/postfinance/kubenurse/commit/3fed2690))
* **helm**: add support for volumes and volumeMounts (#49) ([986d3dc9](https://github.com/postfinance/kubenurse/commit/986d3dc9))
* **helm**: make KUBENURSE_INSECURE configurable (#51) ([4d4dc397](https://github.com/postfinance/kubenurse/commit/4d4dc397))


## 1.5.1 (2022-01-21)


### Bug Fixes

* **common**: enforce timeouts in the kubenurse http.Server to avoid possible goroutine/memory leaks ([d07df3bc](https://github.com/postfinance/kubenurse/commit/d07df3bc))


### New Features

* **common**: expose metrics from the kubenurse httpclient (#31) ([ebb07646](https://github.com/postfinance/kubenurse/commit/ebb07646))
  > The following new metrics were added:
  > * kubenurse_httpclient_requests_total - Total issued requests by kubenurse, partitioned by http code/method.
  > * kubenurse_httpclient_trace_request_duration_seconds - Latency histogram for requests from the kubenurse httpclient, partitioned by event.
  > * httpclient_request_duration_seconds - Latency histogram of request latencies from the kubenurse httpclient.



## 1.5.0 (2022-01-17)


### Breaking Changes

* **common**
  * **[7beac307](https://github.com/postfinance/kubenurse/commit/7beac307)**:
    rewrite and cleanup kubenurse server code  (#29)
    > * refactor!: rewrite and cleanup kubenurse server code
    > 
    > By using a package and multiple separate files the code is easier to
    > understand and test. A new /ready handler was added so we can configure
    > a readiness probe to allow seamless updates of kubenurse.
    > 
    > * build: update golangci-lint version
    > 
    > * build: update golangci-lint timeout, default is too short
    > 
    > * build: extract lint step and use go version 1.17
    > 
    > * feat: configure new readinessprobe in kustomize and helm templates
    > 
    > * fix: linter errors
    > 
    > * chore: cleanup, remove not needed WaitGroup
    > 
    > * refactor!: move pkg/kubediscovery to internal/kubediscovery
    > 
    > * refactor!: move pkg/checker to internal/servicecheck
    > 
    > * refactor!: incorporate pkg/metrics in internal/servicecheck
    > 
    > * refactor!: more refactorings to allow easier unit testing
    > 
    > * feat: more unit tests and coverage calculation in workflows
    > 
    > * docs: include ci and coverage badges in readme
    > 
    > * docs: fix coverage status URL


### New Features

* **common**: expose metrics from the kubenurse httpclient (#31) ([ebb07646](https://github.com/postfinance/kubenurse/commit/ebb07646))
  > The following new metrics were added:
  > * kubenurse_httpclient_requests_total - Total issued requests by kubenurse, partitioned by http code/method.
  > * kubenurse_httpclient_trace_request_duration_seconds - Latency histogram for requests from the kubenurse httpclient, partitioned by event.
  > * httpclient_request_duration_seconds - Latency histogram of request latencies from the kubenurse httpclient.



## 1.4.1 (2021-09-30)


### Bug Fixes

* **examples**: Bump kubenurse version to v1.4.0 ([6f1228c0](https://github.com/postfinance/kubenurse/commit/6f1228c0))



## 1.4.0 (2021-05-25)


### Bug Fixes

* **examples**: Bump kubenurse version to v1.3.4 ([4e0a4c33](https://github.com/postfinance/kubenurse/commit/4e0a4c33))



## 1.3.4 (2021-04-20)


### Bug Fixes

* **discovery**: Prevent panic when checking for schedulable nodes only ([2243226b](https://github.com/postfinance/kubenurse/commit/2243226b))
* **examples**: Bump kubenurse version to v1.3.3 ([c13ebc11](https://github.com/postfinance/kubenurse/commit/c13ebc11))



## 1.3.3 (2021-04-20)


### New Features

* **common**: CI improvements and RBAC fixes ([394daf19](https://github.com/postfinance/kubenurse/commit/394daf19))
* **common**: Flag to consider kubenurses on unschedulable nodes ([cd9ac29b](https://github.com/postfinance/kubenurse/commit/cd9ac29b))



## 1.3.1 (2020-12-09)


### Bug Fixes

* **common**: remove unwanted linter configuration ([d9284394](https://github.com/postfinance/kubenurse/commit/d9284394))



## 1.3.0 (2020-12-09)


### New Features

* **common**: exclude nodes which are not schedulable from neighbour checks ([b6acb939](https://github.com/postfinance/kubenurse/commit/b6acb939))
