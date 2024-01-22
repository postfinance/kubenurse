## 1.9.1 (2024-01-22)


### Dependencies

* **common**: make helm chart version equal to tag ([f248d2ad](https://github.com/postfinance/kubenurse/commit/f248d2ad))
  > [skip ci]



## 1.9.0 (2024-01-22)


### Bug Fixes

* **common**: added missing condition in ingress.yaml chart file ([21165026](https://github.com/postfinance/kubenurse/commit/21165026))
* **common**: create empty tls.Config when loading extraCA fails ([41130653](https://github.com/postfinance/kubenurse/commit/41130653))
* **common**: do not reuse connections per default ([4f1f5b80](https://github.com/postfinance/kubenurse/commit/4f1f5b80))
  > stop reusing connections per default, by modifying the http.Transport
  > DisableKeepAlive field to true per default. Can be overriden through
  > setting KUBENURSE_REUSE_CONNECTIONS env var to "true"
* **common**: linting and error handling ([1057536e](https://github.com/postfinance/kubenurse/commit/1057536e))
* **common**: use same histogram buckets everywhere ([03505e9c](https://github.com/postfinance/kubenurse/commit/03505e9c))
* **helm-lint**: place separator at correct location ([0fa8b064](https://github.com/postfinance/kubenurse/commit/0fa8b064))
* **helm-lint**: place separator at correct location ([c7724bba](https://github.com/postfinance/kubenurse/commit/c7724bba))


### Dependencies

* **common**: bump go version in gh-actions ([fec132df](https://github.com/postfinance/kubenurse/commit/fec132df))
* **deps**: bump github.com/prometheus/client_golang ([533a4ec5](https://github.com/postfinance/kubenurse/commit/533a4ec5))
  > Bumps [github.com/prometheus/client_golang](https://github.com/prometheus/client_golang) from 1.16.0 to 1.18.0.
  > - [Release notes](https://github.com/prometheus/client_golang/releases)
  > - [Changelog](https://github.com/prometheus/client_golang/blob/main/CHANGELOG.md)
  > - [Commits](https://github.com/prometheus/client_golang/compare/v1.16.0...v1.18.0)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/prometheus/client_golang
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump k8s.io/client-go from 0.28.4 to 0.29.0 ([15d67150](https://github.com/postfinance/kubenurse/commit/15d67150))
  > Bumps [k8s.io/client-go](https://github.com/kubernetes/client-go) from 0.28.4 to 0.29.0.
  > - [Changelog](https://github.com/kubernetes/client-go/blob/master/CHANGELOG.md)
  > - [Commits](https://github.com/kubernetes/client-go/compare/v0.28.4...v0.29.0)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/client-go
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...


### New Features

* **httptrace**: add back total and duration instrumentation ([330d2d45](https://github.com/postfinance/kubenurse/commit/330d2d45))
  > as this is used by our Grafana dashboards, and quite useful
  > also prepare for adding the `type` label, once the N^2 issue
  > is solved



## 1.8.1 (2023-12-14)


### Bug Fixes

* **graceful-shutdown**: implement configurable sutdown delay ([e5c13c8d](https://github.com/postfinance/kubenurse/commit/e5c13c8d))
  > prevents in-flight requests to land on a pod that has already stopped.
* **shutdown**: implement 5 seconds shutdown period ([cef5f2ef](https://github.com/postfinance/kubenurse/commit/cef5f2ef))
  > will already help prevent false positives from other kubenurse pods when
  > trying to reach me_ingress through the ingress controller during
  > teardown. without this 5sec wait, in-flight requests from e.g. the
  > ingress controller will reach a pod that is already terminated.
  > Might not be sufficient for similar for "path" errors, as there is no
  > filter for terminating pods.
* **shutdown**: make shutdown duration configurable ([a9d101a4](https://github.com/postfinance/kubenurse/commit/a9d101a4))
* **shutdown**: stop querying pending/terminating neighbors ([3d6050c6](https://github.com/postfinance/kubenurse/commit/3d6050c6))
  > prevents false positive path_error when checks are made to pending or
  > terminating pods


### Dependencies

* **ci**: rollout restart the daemonset to "erase" bootstrap errors ([e96ed6f7](https://github.com/postfinance/kubenurse/commit/e96ed6f7))
  > during the first start of kubenurse, if the ingress isn't
  > ready yet or if kubenurse makes a check before a kubenurse
  > pod is actually ready, this will result in errors in the logs
  > and this will prevent the pipeline from working correctly.
* **common**: helmChart improvements ([6e82de23](https://github.com/postfinance/kubenurse/commit/6e82de23))
  > the image tag is now .Chart.AppVersion except if
  > .values.daemonset.image.tag is set
  > 
  > also, the .Chart.AppVersion field is automatically set
  > to the tag, with another chart-releaser plugin
* **deps**: bump actions/checkout from 3 to 4 ([21c103d0](https://github.com/postfinance/kubenurse/commit/21c103d0))
  > Bumps [actions/checkout](https://github.com/actions/checkout) from 3 to 4.
  > - [Release notes](https://github.com/actions/checkout/releases)
  > - [Changelog](https://github.com/actions/checkout/blob/main/CHANGELOG.md)
  > - [Commits](https://github.com/actions/checkout/compare/v3...v4)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: actions/checkout
  >   dependency-type: direct:production
  >   update-type: version-update:semver-major
  > ...
* **deps**: bump actions/setup-go from 4 to 5 ([b395623c](https://github.com/postfinance/kubenurse/commit/b395623c))
  > Bumps [actions/setup-go](https://github.com/actions/setup-go) from 4 to 5.
  > - [Release notes](https://github.com/actions/setup-go/releases)
  > - [Commits](https://github.com/actions/setup-go/compare/v4...v5)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: actions/setup-go
  >   dependency-type: direct:production
  >   update-type: version-update:semver-major
  > ...
* **deps**: bump docker/login-action from 2 to 3 ([87f61117](https://github.com/postfinance/kubenurse/commit/87f61117))
  > Bumps [docker/login-action](https://github.com/docker/login-action) from 2 to 3.
  > - [Release notes](https://github.com/docker/login-action/releases)
  > - [Commits](https://github.com/docker/login-action/compare/v2...v3)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: docker/login-action
  >   dependency-type: direct:production
  >   update-type: version-update:semver-major
  > ...
* **deps**: bump helm/chart-releaser-action from 1.5.0 to 1.6.0 ([efc98fa1](https://github.com/postfinance/kubenurse/commit/efc98fa1))
  > Bumps [helm/chart-releaser-action](https://github.com/helm/chart-releaser-action) from 1.5.0 to 1.6.0.
  > - [Release notes](https://github.com/helm/chart-releaser-action/releases)
  > - [Commits](https://github.com/helm/chart-releaser-action/compare/v1.5.0...v1.6.0)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: helm/chart-releaser-action
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump k8s.io/api from 0.27.3 to 0.27.4 ([7ad9eb28](https://github.com/postfinance/kubenurse/commit/7ad9eb28))
  > Bumps [k8s.io/api](https://github.com/kubernetes/api) from 0.27.3 to 0.27.4.
  > - [Commits](https://github.com/kubernetes/api/compare/v0.27.3...v0.27.4)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/api
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump k8s.io/api from 0.28.0 to 0.28.1 ([ca5a74c3](https://github.com/postfinance/kubenurse/commit/ca5a74c3))
  > Bumps [k8s.io/api](https://github.com/kubernetes/api) from 0.28.0 to 0.28.1.
  > - [Commits](https://github.com/kubernetes/api/compare/v0.28.0...v0.28.1)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/api
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump k8s.io/client-go from 0.27.3 to 0.28.0 ([7791489d](https://github.com/postfinance/kubenurse/commit/7791489d))
  > Bumps [k8s.io/client-go](https://github.com/kubernetes/client-go) from 0.27.3 to 0.28.0.
  > - [Changelog](https://github.com/kubernetes/client-go/blob/master/CHANGELOG.md)
  > - [Commits](https://github.com/kubernetes/client-go/compare/v0.27.3...v0.28.0)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/client-go
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump k8s.io/client-go from 0.28.0 to 0.28.1 ([52bfac3f](https://github.com/postfinance/kubenurse/commit/52bfac3f))
  > Bumps [k8s.io/client-go](https://github.com/kubernetes/client-go) from 0.28.0 to 0.28.1.
  > - [Changelog](https://github.com/kubernetes/client-go/blob/master/CHANGELOG.md)
  > - [Commits](https://github.com/kubernetes/client-go/compare/v0.28.0...v0.28.1)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/client-go
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump k8s.io/client-go from 0.28.1 to 0.28.4 ([eb3c96c1](https://github.com/postfinance/kubenurse/commit/eb3c96c1))
  > Bumps [k8s.io/client-go](https://github.com/kubernetes/client-go) from 0.28.1 to 0.28.4.
  > - [Changelog](https://github.com/kubernetes/client-go/blob/master/CHANGELOG.md)
  > - [Commits](https://github.com/kubernetes/client-go/compare/v0.28.1...v0.28.4)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/client-go
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **dockerfile**: update misconfigured maintainer label ([461bda53](https://github.com/postfinance/kubenurse/commit/461bda53))


### New Features

* **helm**: make shutdown duration configurable ([a518f562](https://github.com/postfinance/kubenurse/commit/a518f562))



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
