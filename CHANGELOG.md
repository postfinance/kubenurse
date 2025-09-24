
## [1.15.1] - 2025-09-24

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.15.0...v1.15.1

### 🚀 Features

- Extra objects in helm (#219) - ([56cde01](https://github.com/postfinance/kubenurse/commit/56cde01c7289086057d6d169ac7da289a4c8c919))
- Allow ipv6 services (#218) - ([c5aa96c](https://github.com/postfinance/kubenurse/commit/c5aa96cdee73481488fc1d2562a7c0777be23f12))

### 🐛 Bug Fixes

- Run collectMetrics func in a dedicated goroutine - ([a68f4a2](https://github.com/postfinance/kubenurse/commit/a68f4a2672985bb2399e3726581837edf533236d))

### ⚙️ Miscellaneous Tasks

- Update changelog - ([44f9c9d](https://github.com/postfinance/kubenurse/commit/44f9c9de8349e70016b83b1ca01ec765c523a7d0))

### Build

- *(ci)* Update goveralls to v0.0.12 - ([5851a2f](https://github.com/postfinance/kubenurse/commit/5851a2fa2ab7921e9b4a0b97682b66633d14ef7c))
- *(deps)* Bump actions/setup-go from 5 to 6 (#232) - ([28c5ee4](https://github.com/postfinance/kubenurse/commit/28c5ee4fd06ca8a5dfe7d697f32b74f55cd41e0a))
- *(deps)* Bump actions/checkout from 4 to 5 (#224) - ([a6dbd88](https://github.com/postfinance/kubenurse/commit/a6dbd888241f07ef7262ec3a7e9cd069f5b37678))
- *(helm)* Add ndots:1 suggestion in values.yaml - ([2672af2](https://github.com/postfinance/kubenurse/commit/2672af2f18ef27aeae8b213c5a454d34640491b7))
- Cleanup github-action and fix ci - ([2589358](https://github.com/postfinance/kubenurse/commit/258935804f6aceb7dd976fd0ef05595d6cb2dde8))
- Update all go deps - ([ff36a41](https://github.com/postfinance/kubenurse/commit/ff36a41f90fa1de2d9f930076c74ec47810bc41d))




## New Contributors
* @borislitv made their first contribution in [#218](https://github.com/postfinance/kubenurse/pull/218)
* @github-actions[bot] made their first contribution## [1.15.0] - 2025-06-16

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.14.6...v1.15.0

### 🚀 Features

- Introduce VictoriaMetrics/metrics histograms (#214) - ([b640a2b](https://github.com/postfinance/kubenurse/commit/b640a2ba58fa95dd1ef0e970fae4ee6064120100))

### Build

- Automatically generate CHANGELOG.md - ([7d43e8a](https://github.com/postfinance/kubenurse/commit/7d43e8a4ceb47af4c03cb8a5b5a0d402e965a1c6))


## [1.14.6] - 2025-06-02

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.14.5...v1.14.6

### 🐛 Bug Fixes

- *(controller-runtime)* Remove logger error - ([49087be](https://github.com/postfinance/kubenurse/commit/49087be65bf674a74dc185ecac62f12e12d7b433))
- *(helm)* Use %v to print the port - ([ec38a80](https://github.com/postfinance/kubenurse/commit/ec38a807e3c480ffbb681afeb070a0e33773328c))


## [1.14.5] - 2025-06-02

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.14.4...v1.14.5

### 🐛 Bug Fixes

- *(helm)* Change service_url port template to %s - ([83251dc](https://github.com/postfinance/kubenurse/commit/83251dc8ac6025f5fa37b2b8aef044b5be8a20f0))

### Build

- *(deps)* Bump sigs.k8s.io/controller-runtime from 0.20.4 to 0.21.0 - ([b64a8bd](https://github.com/postfinance/kubenurse/commit/b64a8bdb2dd05191f5438b37f171b77e966a1e33))
- *(deps)* Bump k8s.io/client-go from 0.33.0 to 0.33.1 - ([515c081](https://github.com/postfinance/kubenurse/commit/515c081751b959b475a8845cc29157ec5628f7ed))
- Bump to go v1.24 - ([f48166a](https://github.com/postfinance/kubenurse/commit/f48166ad6086be5d87fb1765c507a258c34ca27b))
- Update gh-action to v1.24 - ([0a6923a](https://github.com/postfinance/kubenurse/commit/0a6923a27c211d2887ccd2a8c43d7a2b00ca29ec))


## [1.14.4] - 2025-04-14

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.14.3...v1.14.4

### 🐛 Bug Fixes

- Match dialer timeout with http request timeout - ([5e64bac](https://github.com/postfinance/kubenurse/commit/5e64bacd3a42e18f46830ce2c09b2a4b780a7b8e))

### Build

- *(deps)* Bump github.com/prometheus/client_golang - ([df248d6](https://github.com/postfinance/kubenurse/commit/df248d6958dd0217e79dc067e70ff1d1ce084341))
- *(deps)* Bump k8s.io/client-go from 0.32.1 to 0.32.2 - ([31e782e](https://github.com/postfinance/kubenurse/commit/31e782e943bed64144dac474aa1b12dccd21767f))
- *(deps)* Bump sigs.k8s.io/controller-runtime from 0.20.1 to 0.20.2 - ([7a92bbf](https://github.com/postfinance/kubenurse/commit/7a92bbf39f44a7713a84f7538f4a3d89e8f58b14))
- *(deps)* Bump sigs.k8s.io/controller-runtime from 0.20.2 to 0.20.3 - ([2e7439b](https://github.com/postfinance/kubenurse/commit/2e7439be3f90598802d39fbd9386e4e72e50b051))
- *(deps)* Bump github.com/prometheus/client_golang - ([b7776ab](https://github.com/postfinance/kubenurse/commit/b7776ab679b9e02664a623f75dc4b42519c32474))
- *(deps)* Bump golangci/golangci-lint-action from 6 to 7 - ([96e87d1](https://github.com/postfinance/kubenurse/commit/96e87d1befd9cef523c5dbe69b0db212f3a45099))
- *(deps)* Bump sigs.k8s.io/controller-runtime from 0.20.3 to 0.20.4 - ([869eb85](https://github.com/postfinance/kubenurse/commit/869eb859ed88a9875e7f5c094f1cd0ea88636c77))
- *(deps)* Bump k8s.io/client-go from 0.32.2 to 0.32.3 - ([af23821](https://github.com/postfinance/kubenurse/commit/af23821a09ee76772c790fb66163aca89dd31ce6))
- Migrate linter to v2 - ([d6c3535](https://github.com/postfinance/kubenurse/commit/d6c3535d740523e35fcd85b7eed174469c364d6a))
- Upgrade to golangci-lint v2 - ([553de54](https://github.com/postfinance/kubenurse/commit/553de546df1662ad7dca0dfd52626300688ec7e2))




## New Contributors
* @marsom made their first contribution in [#196](https://github.com/postfinance/kubenurse/pull/196)## [1.14.3] - 2025-02-12

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.14.2...v1.14.3

### ⚙️ Miscellaneous Tasks

- Add arm64 go builds and docker images - ([7cee33a](https://github.com/postfinance/kubenurse/commit/7cee33a4607b32dba096b63c047550e48a8fae8c))
- Fix linting errors and update to go 1.23 - ([2a0767b](https://github.com/postfinance/kubenurse/commit/2a0767be7087387532bcf2ad2f97647cd49be2c0))

### Build

- *(deps)* Bump github.com/prometheus/client_golang - ([e69c609](https://github.com/postfinance/kubenurse/commit/e69c609885229cde02b8c5fe41188355367e9c48))
- *(deps)* Bump sigs.k8s.io/controller-runtime from 0.19.0 to 0.19.3 - ([6b19121](https://github.com/postfinance/kubenurse/commit/6b19121a5a0642ed22befb285acc37087c166184))
- *(deps)* Bump sigs.k8s.io/controller-runtime from 0.19.3 to 0.20.0 - ([53e466f](https://github.com/postfinance/kubenurse/commit/53e466f1003b91c4a5425ba51931432cc84b24d6))
- *(deps)* Update golangci-lint and go in gh-action - ([9ec0c57](https://github.com/postfinance/kubenurse/commit/9ec0c57170a82ec771ec741482d3055f5a13e0ea))
- *(deps)* Bump github.com/stretchr/testify from 1.9.0 to 1.10.0 - ([53e6e93](https://github.com/postfinance/kubenurse/commit/53e6e938bede46268b35e9e5ab91ebf33fea5c02))
- *(deps)* Bump engineerd/setup-kind from 0.5.0 to 0.6.2 - ([c8c7762](https://github.com/postfinance/kubenurse/commit/c8c77620a03274a476a97ba0bb3617c02b378e67))
- *(deps)* Bump k8s.io/client-go from 0.31.1 to 0.32.1 - ([607e294](https://github.com/postfinance/kubenurse/commit/607e294ad20e15621ec9c9df4e8f15c915b248ee))
- *(deps)* Bump sigs.k8s.io/controller-runtime from 0.20.0 to 0.20.1 - ([30802a7](https://github.com/postfinance/kubenurse/commit/30802a7aa52d7ddde215ba5181ab12dc8d9a848f))
- Update kind version to v0.24 - ([4a504d1](https://github.com/postfinance/kubenurse/commit/4a504d14c716495a489c99856fe71336a3e43fad))




## New Contributors
* @fmarani made their first contribution## [1.14.2] - 2024-10-02

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.14.1...v1.14.2

### 🚀 Features

- Support setting .endpoints[].relabelings in ServiceMonitor - ([d369419](https://github.com/postfinance/kubenurse/commit/d3694193c0386735920acf0fdd562b3c7a656d04))

### 🐛 Bug Fixes

- Don't set spec.updateStrategy.maxSurge of DaemonSet by default - ([3facb9c](https://github.com/postfinance/kubenurse/commit/3facb9ca29504fc372e52d15421d1cc2929843c8))

### 📚 Documentation

- Add the serviceMonitor.relabelings Helm parameter - ([2a107b0](https://github.com/postfinance/kubenurse/commit/2a107b05d9431a226f37eb6fdbf4c7ae42cd623d))
- Add the daemonset.updateStrategy.maxUnavailable and daemonset.updateStrategy.maxSurge Helm parameters - ([143a0e1](https://github.com/postfinance/kubenurse/commit/143a0e1dbdbd89b22273b460445350eff1038f91))

### Build

- *(goreleaser)* Update deprecated field name - ([4a5a679](https://github.com/postfinance/kubenurse/commit/4a5a679fcc346ecbe426d706aedaa32d73d8e9c4))




## New Contributors
* @justin0u0 made their first contribution## [1.14.1] - 2024-09-30

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.14.0...v1.14.1

### 🚀 Features

- Permit specifying API server DNS name - ([545d7b6](https://github.com/postfinance/kubenurse/commit/545d7b6df8763fc99b378cb04e54ce445a5deda3))

### Build

- *(deps)* Bump github.com/prometheus/client_golang - ([a4fb79b](https://github.com/postfinance/kubenurse/commit/a4fb79b607b4e0793294180fe28fd3da161261d6))
- *(deps)* Bump orhun/git-cliff-action from 3 to 4 - ([6bc186f](https://github.com/postfinance/kubenurse/commit/6bc186fbcd3a7e29367a714c14beddabf7522855))
- *(deps)* Bump github.com/prometheus/client_golang - ([38e9083](https://github.com/postfinance/kubenurse/commit/38e90832c79c27a819709b18bec72779cc1a5e09))
- *(deps)* Bump sigs.k8s.io/controller-runtime from 0.18.4 to 0.19.0 - ([206ecd9](https://github.com/postfinance/kubenurse/commit/206ecd951d4a3d972cb4803b67b374d2d74cff94))
- *(deps)* Bump k8s.io/client-go from 0.31.0 to 0.31.1 - ([7b42df7](https://github.com/postfinance/kubenurse/commit/7b42df71118e80cf4f6f6308bb8e6d09bb237a7f))
- *(deps)* Bump github.com/prometheus/client_golang - ([b87a8ee](https://github.com/postfinance/kubenurse/commit/b87a8ee892983577f3a7698fb705b497e894673e))
- *(helm)* Fix helm variable name used to disable ingress/service - ([754b407](https://github.com/postfinance/kubenurse/commit/754b407ce768dc220984b8cb8afcba2a4f1d5c11))
- Add bin and git-cliff to gitignore - ([eec4466](https://github.com/postfinance/kubenurse/commit/eec4466d4ff73a71525cb11926d57eed4f5e2b46))


## [1.14.0] - 2024-08-07

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.13.1...v1.14.0

### 🚀 Features

- Run checks in parallel - ([04caaa4](https://github.com/postfinance/kubenurse/commit/04caaa4ed82eb14cca74060a5f00e1861da55fd7))
- Implement additional checks functionality - ([4d4b962](https://github.com/postfinance/kubenurse/commit/4d4b9626bb9bc352b3aa8f3aef165afe9f5478d4))

### 🐛 Bug Fixes

- Include failed requests in total count - ([e5e0c55](https://github.com/postfinance/kubenurse/commit/e5e0c55967ae484170a132994a8ce40d87d3f560))

### 🚜 Refactor

- Turn the result struct to a map - ([be906bb](https://github.com/postfinance/kubenurse/commit/be906bbd375403f3961f8b5bbce0c9ba4d2b1845))
- Simplify check functions - ([b37c8d0](https://github.com/postfinance/kubenurse/commit/b37c8d06d808ee8665cea59702af963fef5085b0))
- Further simplify Check function - ([a2d392f](https://github.com/postfinance/kubenurse/commit/a2d392ff0855fd28b77fb4c48a7f7c5c7fb6efaa))

### 📚 Documentation

- Introduce extra_checks feature and corresponding env - ([656a4a6](https://github.com/postfinance/kubenurse/commit/656a4a6f56bac3712ae0422dacc4e7180290d1c9))

### ⚙️ Miscellaneous Tasks

- *(helm)* Make TLS secret name configurable and optional - ([943bf66](https://github.com/postfinance/kubenurse/commit/943bf66e6c7fa8f6bca2e9f9adbeea8a6a66c9ea))
- Use atomic.Bool for the ready flag - ([f64d442](https://github.com/postfinance/kubenurse/commit/f64d442fabe39f7e5c1118e26b1fa9d988d3954d))
- Proper Context handling and simplifications - ([cafb7ae](https://github.com/postfinance/kubenurse/commit/cafb7aeb4d92b5191aa14c9eed6e57f3b1c2236e))
- Upgrade dependencies - ([ba240e9](https://github.com/postfinance/kubenurse/commit/ba240e95b98e8755bd5b0da8a2c8b5b2262060b4))
- Add v1.14.0 changelog - ([b369b6e](https://github.com/postfinance/kubenurse/commit/b369b6e80b8d77ba2e8410930ad32eb001d4e857))

### Build

- *(deps)* Bump goreleaser/goreleaser-action from 4 to 6 - ([beb1c4b](https://github.com/postfinance/kubenurse/commit/beb1c4bcf52d718a7bf2ba8056e9db34db04322e))
- Fix goreleaser command - ([b07ada8](https://github.com/postfinance/kubenurse/commit/b07ada83eee9a72fce3f30017ddaa726783a3505))
- Add goreleaser config version - ([5fc9e99](https://github.com/postfinance/kubenurse/commit/5fc9e990678d3f7bccd4f6b8fea1b0584b658800))




## New Contributors
* @karnzx made their first contribution in [#141](https://github.com/postfinance/kubenurse/pull/141)## [1.13.1] - 2024-06-03

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.13.0...v1.13.1

### 🐛 Bug Fixes

- StatusCode != 200 trigger an error and increase error counter - ([b403ddd](https://github.com/postfinance/kubenurse/commit/b403dddabbcd5ee30a532c05a735754cb7c95044))
- Only increase error counter once when catching roundTripper errors - ([740b8d4](https://github.com/postfinance/kubenurse/commit/740b8d46f557039e924200fd87976192867c50a3))

### 📚 Documentation

- *(drawings)* Change request_type to request type - ([175c17c](https://github.com/postfinance/kubenurse/commit/175c17cec93f373166a4df042d34085659df67c2))
- *(grafana)* Add updated grafana dashboard - ([2760b9f](https://github.com/postfinance/kubenurse/commit/2760b9f41ef7f2d44cec416067c67dfb4bee4df0))
- *(readme)* Simplify documentation and improve readbility - ([983a2bc](https://github.com/postfinance/kubenurse/commit/983a2bca4fc0011ceb000c46b9f501b514865315))

### ⚙️ Miscellaneous Tasks

- Update changelog - ([154f931](https://github.com/postfinance/kubenurse/commit/154f93114a4d586581f8bffe3a816ab4008db65e))

### Build

- *(deps)* Bump sigs.k8s.io/controller-runtime from 0.17.2 to 0.17.3 - ([7966f1c](https://github.com/postfinance/kubenurse/commit/7966f1c39b6fc386ba50359472ea61225b073864))
- *(deps)* Bump golangci/golangci-lint-action from 4 to 5 - ([9ac2183](https://github.com/postfinance/kubenurse/commit/9ac2183c35e46fc96700119d7c6f27bcbc20bcdd))
- *(deps)* Bump sigs.k8s.io/controller-runtime from 0.18.0 to 0.18.1 - ([504c082](https://github.com/postfinance/kubenurse/commit/504c082245dd89a5a32d8a6a5784f0615f339cc2))
- *(deps)* Bump golangci/golangci-lint-action from 5 to 6 - ([b8e6781](https://github.com/postfinance/kubenurse/commit/b8e678178493358f75f3dc73b7c2d823833be718))
- *(deps)* Bump sigs.k8s.io/controller-runtime from 0.18.1 to 0.18.2 - ([7acb1db](https://github.com/postfinance/kubenurse/commit/7acb1db6c81344cefd665c015dbf0cd205b91f24))
- *(deps)* Bump github.com/prometheus/client_golang - ([5067e6c](https://github.com/postfinance/kubenurse/commit/5067e6cdb99245ac45bd96c54ce29e2db2e1c8fc))
- *(deps)* Bump sigs.k8s.io/controller-runtime from 0.18.2 to 0.18.3 - ([0fcddee](https://github.com/postfinance/kubenurse/commit/0fcddeeccb56f69c83a76853d4ce6e2e297b27e8))
- Controller-runtime v1.18 with k8s 1.30 client - ([c5534da](https://github.com/postfinance/kubenurse/commit/c5534da37499d7471f471f8a470249db2f0feaa9))


## [1.13.0] - 2024-04-15

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.12.3...v1.13.0

### 🚀 Features

- [**breaking**] Use slog and add event type to error counter - ([e9d51a2](https://github.com/postfinance/kubenurse/commit/e9d51a2f781eaa2cb7e415fd175d158fdd681bb3))

### 🐛 Bug Fixes

- *(ci)* Change kubenurse startup string in test condition - ([bd91d7d](https://github.com/postfinance/kubenurse/commit/bd91d7da68b3408d70dfe846562557396b4522c4))
- *(ci)* Update metric name for sanity check - ([c3e7287](https://github.com/postfinance/kubenurse/commit/c3e72874e509cff3fa72381d8aa882878ef32385))
- *(server)* Improve error handling for histogram buckets parsing - ([e324325](https://github.com/postfinance/kubenurse/commit/e3243259626f3e4f2fd4ad68307205298b340276))

### 🚜 Refactor

- [**breaking**] Remove kubenurse_request_duration histogram - ([3061b64](https://github.com/postfinance/kubenurse/commit/3061b6490d26e39bf096d050e54e15f7d43c0f8e))

### 📚 Documentation

- *(drawings)* Add excalidraw file - ([852ab85](https://github.com/postfinance/kubenurse/commit/852ab85fc2ec00804a290730b5276fb9689116f4))
- *(readme.md)* Format metrics as a table - ([6a35209](https://github.com/postfinance/kubenurse/commit/6a3520955e2e91e3cccb203c528e23cd2032e127))

### ⚙️ Miscellaneous Tasks

- Add v1.13.0 changelog - ([ee0ad5e](https://github.com/postfinance/kubenurse/commit/ee0ad5e685a10571ec322eebd279ff2e3dacbc8b))
- Update changelog - ([55bd79f](https://github.com/postfinance/kubenurse/commit/55bd79f6888910b60eac9e7c498b2b0ae0dbc1f8))

### Build

- *(artifacthub)* Add validation file - ([e2bf079](https://github.com/postfinance/kubenurse/commit/e2bf079504223a3293ed1be73d061bed49e00e02))


## [1.12.3] - 2024-04-09

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.12.2...v1.12.3

### Build

- *(helm)* Add artifacthub metadata - ([6466e17](https://github.com/postfinance/kubenurse/commit/6466e17a25352baa8f82c6d4a0420d35e9f4c702))


## [1.12.2] - 2024-04-09

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.12.1...v1.12.2

### 🚀 Features

- *(helm)* Set rollingUpdate.maxUnavailable to 34% per default - ([9c44a69](https://github.com/postfinance/kubenurse/commit/9c44a69b697f4f389d888e3ee64c551cf3503751))

### 📚 Documentation

- Add drawings as png again - ([2fc3d08](https://github.com/postfinance/kubenurse/commit/2fc3d088dfa01ad3d0f0f2a6cc08a136c5e84e3b))

### Build

- *(ci)* Use git-cliff release notes - ([8f8c64f](https://github.com/postfinance/kubenurse/commit/8f8c64f067d6402319145e6217e53fb01e477fb4))


## [1.12.1] - 2024-04-05

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.12.0...v1.12.1

### 📚 Documentation

- *(README)* Add background to drawings - ([bb1a24d](https://github.com/postfinance/kubenurse/commit/bb1a24d70cfd1004d4ac59640a4d8a706ab53cba))

### ⚙️ Miscellaneous Tasks

- *(changelog)* Add v1.12.1 changelog - ([5cfd422](https://github.com/postfinance/kubenurse/commit/5cfd42237b30b64a91ce38f9b58a49d8c2c38940))

### Build

- *(ci)* Switch branch name to main - ([7d06dea](https://github.com/postfinance/kubenurse/commit/7d06deac3f4d4170e8d15650d03f1dfb19ecc773))
- *(helm)* Add values.histogram_buckets option - ([9824bc7](https://github.com/postfinance/kubenurse/commit/9824bc72de66ebc6474a149cd19491f1175de335))


## [1.12.0] - 2024-04-05

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.11.0...v1.12.0

### 🚀 Features

- Recent neighbours checks metrics - ([145f6b5](https://github.com/postfinance/kubenurse/commit/145f6b51e1b7bfbc4d0ac0e7e8c74ef26f7c01f7))
- Add incoming neighbouring checks gauge - ([b27de93](https://github.com/postfinance/kubenurse/commit/b27de934e8ec4f39b70f1ecc844304b25d841aaf))

### 🐛 Bug Fixes

- Use full URL for changelog commit ids - ([ee3951f](https://github.com/postfinance/kubenurse/commit/ee3951fa4d526be837f7b98c94001553afcbd4a4))

### 📚 Documentation

- *(README)* Add ToC and a drawing - ([99b52d8](https://github.com/postfinance/kubenurse/commit/99b52d8cb0520a6d2c496843036dc67f287fb730))
- *(README)* Reduce doctoc maxlevel and reorder badges - ([cd175cd](https://github.com/postfinance/kubenurse/commit/cd175cd65455ade1937ab161f89854415fd70a03))
- *(node filtering)* Add neighbourhood incoming checks metric and drawing - ([f8b17eb](https://github.com/postfinance/kubenurse/commit/f8b17eb74f36594143aaf799a7978dd4bc9bc2a9))
- *(ttl cache)* Explain utility and future improvements - ([9943b92](https://github.com/postfinance/kubenurse/commit/9943b9235558aca8f989b532a88cf58eda8ef9b3))

### 🧪 Testing

- *(ttl_cache)* Add basic test - ([2efa0d8](https://github.com/postfinance/kubenurse/commit/2efa0d84c1fe724017affc50b9ceaae456dc4392))
- Add basic TTLCache test - ([fa24d14](https://github.com/postfinance/kubenurse/commit/fa24d149375345c771794909954c3f0aa6fe7eec))

### Build

- *(deps)* Bump k8s.io/api from 0.29.2 to 0.29.3 - ([4cd0aa4](https://github.com/postfinance/kubenurse/commit/4cd0aa49eeea2c1458ed1a237f8d28ef93d59129))
- *(deps)* Bump k8s.io/client-go from 0.29.2 to 0.29.3 - ([c883fb1](https://github.com/postfinance/kubenurse/commit/c883fb1b3a6063c9fed3bfdb661dc735df4d1d3b))
- Tailor git-cliff to add scope to changelog + release v1.12.0 - ([aa8c556](https://github.com/postfinance/kubenurse/commit/aa8c5566bf9e40421bd1579b57efa15b05828281))


## [1.11.0] - 2024-03-15

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.10.0...v1.11.0

### 🚀 Features

- Use hashing to distribute node checks - ([a0b49bb](https://github.com/postfinance/kubenurse/commit/a0b49bbbac77022cefe28c852a8e33c06a764431))
- Use uint64 hashes and store neighbours in a heap - ([270e208](https://github.com/postfinance/kubenurse/commit/270e20837a8967f7054859c93d66fa36e7e6365f))
- Add request type to httptrace and request duration metrics - ([cdcc063](https://github.com/postfinance/kubenurse/commit/cdcc0633a3a084ed17f29dd9d37d6d8109c87911))

### 🐛 Bug Fixes

- Current node hash can never be in the map - ([5753890](https://github.com/postfinance/kubenurse/commit/5753890a8fd224b49c5022c4f109b7605f4a56c6))

### 🚜 Refactor

- Put Uint64Heap at the end of servicecheck.go - ([4950dd6](https://github.com/postfinance/kubenurse/commit/4950dd631b71215121b7cf6e5d755902b7cb29b7))

### 📚 Documentation

- Neighbour filtering - ([bd1ee9f](https://github.com/postfinance/kubenurse/commit/bd1ee9f472aee9686664304e582f39326dfb3970))

### ⚙️ Miscellaneous Tasks

- Linting - ([7aac7e0](https://github.com/postfinance/kubenurse/commit/7aac7e0ec18f496488f9234147614064199ecb29))
- Switch to Go 1.22 - ([1689c1c](https://github.com/postfinance/kubenurse/commit/1689c1c6bfd252f5dc6fe5560f49bcaf1224c578))

### Build

- *(deps)* Bump k8s.io/api from 0.29.0 to 0.29.2 - ([48af8fc](https://github.com/postfinance/kubenurse/commit/48af8fceab5269007c600f5a444e7b3a3569872b))
- *(deps)* Bump k8s.io/client-go from 0.29.0 to 0.29.2 - ([e4734c8](https://github.com/postfinance/kubenurse/commit/e4734c88eb0a64d924ae3c03f62ed077dbc368f2))
- *(deps)* Bump github.com/stretchr/testify from 1.8.4 to 1.9.0 - ([a06bffa](https://github.com/postfinance/kubenurse/commit/a06bffaf5a0b81682f193f50411dd01d51d569aa))
- *(deps)* Bump azure/setup-helm from 3 to 4 - ([688d08b](https://github.com/postfinance/kubenurse/commit/688d08bfcc07ac2a01cecd0b8f1eb9fe052f79b9))
- *(deps)* Bump sigs.k8s.io/controller-runtime from 0.17.0 to 0.17.2 - ([8837f46](https://github.com/postfinance/kubenurse/commit/8837f46b557e0259cfe88a2887de28901ace2439))
- *(deps)* Bump github.com/prometheus/client_golang - ([fa80824](https://github.com/postfinance/kubenurse/commit/fa80824f202dce96b9e9ebccb41e27493bd187f7))
- Switch changelog tool to cliff + release 1.10.0 - ([1cd6d6b](https://github.com/postfinance/kubenurse/commit/1cd6d6b295e7dd6375811b49c00a777585b316d6))


## [1.10.0] - 2024-02-20

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.9.1...v1.10.0

### 🚀 Features

- Use controller-runtime's client with integrated caching - ([7b1edea](https://github.com/postfinance/kubenurse/commit/7b1edea5524f01fe46e717c08615ccd483346a5b))

### 🐛 Bug Fixes

- *(neighbours)* Only check other kubenurse pods - ([62e737c](https://github.com/postfinance/kubenurse/commit/62e737c9d4cb73df8ea6ea9f0a32c2551d020ac7))
- Don't log nil error returned when the cache terminates - ([8d891b6](https://github.com/postfinance/kubenurse/commit/8d891b6a209654a579d495da2754564abfcd6373))

### ⚙️ Miscellaneous Tasks

- Remove "caching" of results and simplify code - ([92b4922](https://github.com/postfinance/kubenurse/commit/92b4922d1814a8b65375a0bafc59d92465d59a62))
- Update changelog with 1.10.0 release - ([0426258](https://github.com/postfinance/kubenurse/commit/0426258c1509346998a1144ba7d9e19525687439))

### Build

- *(deps)* Bump golangci/golangci-lint-action from 3 to 4 - ([8efc905](https://github.com/postfinance/kubenurse/commit/8efc90502f0488c08cad811dab217875e175d7e8))


## [1.9.1] - 2024-01-22

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.9.0...v1.9.1

### ⚙️ Miscellaneous Tasks

- Update changelog with 1.9.1 release - ([96a1713](https://github.com/postfinance/kubenurse/commit/96a1713a54ec24fc8ca8f51dd0c1fc2ef1155fca))

### Build

- Make helm chart version equal to tag - ([f248d2a](https://github.com/postfinance/kubenurse/commit/f248d2adc176a049ce4a88b801b53dc8ed1412a1))


## [1.9.0] - 2024-01-22

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.8.1...v1.9.0

### 🚀 Features

- *(httptrace)* Add back total and duration instrumentation - ([330d2d4](https://github.com/postfinance/kubenurse/commit/330d2d4598f9a4add4ac6aa7a3e43219bf60ad29))

### 🐛 Bug Fixes

- *(helm-lint)* Place separator at correct location - ([c7724bb](https://github.com/postfinance/kubenurse/commit/c7724bba8e85747ac7b33c43fec7bd0477876435))
- *(helm-lint)* Place separator at correct location - ([0fa8b06](https://github.com/postfinance/kubenurse/commit/0fa8b064278750b346b2979cf415774e777c4c8a))
- Added missing condition in ingress.yaml chart file - ([2116502](https://github.com/postfinance/kubenurse/commit/211650261b61cc976a72f38b39c6cca157c2acb6))
- Linting and error handling - ([1057536](https://github.com/postfinance/kubenurse/commit/1057536eb620a83742a94d05a31c6a9239492276))
- Do not reuse connections per default - ([4f1f5b8](https://github.com/postfinance/kubenurse/commit/4f1f5b809aae3aedfd77045cf115e204136b583f))
- Create empty tls.Config when loading extraCA fails - ([4113065](https://github.com/postfinance/kubenurse/commit/411306537de72bbdf255fb749ddba24143893af2))
- Use same histogram buckets everywhere - ([03505e9](https://github.com/postfinance/kubenurse/commit/03505e9ca01229c4c60d7e24b8c6ba2a87a10a16))

### 📚 Documentation

- Customizable histogram buckets with env var - ([dd7ce2d](https://github.com/postfinance/kubenurse/commit/dd7ce2db649258e6a8cd6846e8fca876cffcb7b4))
- Reuse_connections option/env variable - ([9cb33d7](https://github.com/postfinance/kubenurse/commit/9cb33d7e18ce8efcb02a52dbf105dd3bed23bbb3))

### ⚙️ Miscellaneous Tasks

- *(linting)* Set tls.Config.MinVersion per default - ([f32c37b](https://github.com/postfinance/kubenurse/commit/f32c37ba06192004cef4a63d9a0bcf50f5fee8b8))
- Update changelog with 1.9.0 release - ([7c03ef1](https://github.com/postfinance/kubenurse/commit/7c03ef1b3ef2e738ada49ea16624f0059493f32c))

### WIP

- Feat: Replacing promhttp with own httptrace and logging - ([ff0e1b0](https://github.com/postfinance/kubenurse/commit/ff0e1b06721328c47d1dc7edbcc3ac022c17386b))

### Build

- *(deps)* Bump k8s.io/client-go from 0.28.4 to 0.29.0 - ([15d6715](https://github.com/postfinance/kubenurse/commit/15d671503ada336440b599800d2fc726b14c2236))
- *(deps)* Bump github.com/prometheus/client_golang - ([533a4ec](https://github.com/postfinance/kubenurse/commit/533a4ec5a945b023d64c6afc7ef62c55bbf88a2f))
- Bump go version in gh-actions - ([fec132d](https://github.com/postfinance/kubenurse/commit/fec132df47f498d634b136069621ce9be1339a22))




## New Contributors
* @matthisholleville made their first contribution## [1.8.1] - 2023-12-14

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.7.1...v1.8.1

### 🚀 Features

- *(helm)* Make shutdown duration configurable - ([a518f56](https://github.com/postfinance/kubenurse/commit/a518f56288836c025346d63ef154239bcfafdc3e))

### 🐛 Bug Fixes

- *(graceful-shutdown)* Implement configurable sutdown delay - ([e5c13c8](https://github.com/postfinance/kubenurse/commit/e5c13c8df0144c30e308cc833da8f6f0caac38ef))
- *(shutdown)* Implement 5 seconds shutdown period - ([cef5f2e](https://github.com/postfinance/kubenurse/commit/cef5f2effeb1de8ee41602c74a20d8d57e5f450d))
- *(shutdown)* Stop querying pending/terminating neighbors - ([3d6050c](https://github.com/postfinance/kubenurse/commit/3d6050c652bd6115a8be7f7879fa316c4bbad8d5))
- *(shutdown)* Make shutdown duration configurable - ([a9d101a](https://github.com/postfinance/kubenurse/commit/a9d101a4fc3607ebd914b462776b43462d354e0d))

### ⚙️ Miscellaneous Tasks

- Fix linting errors and update golangci-lint config - ([65ee3ec](https://github.com/postfinance/kubenurse/commit/65ee3ece3f306eac32cc34a456dc1464f1f3d44b))
- Update changelog with 1.8.1 release - ([e54d02d](https://github.com/postfinance/kubenurse/commit/e54d02d2e9618c2a4d9e4299cf61532bc7f8d74e))

### Build

- *(ci)* Rollout restart the daemonset to "erase" bootstrap errors - ([e96ed6f](https://github.com/postfinance/kubenurse/commit/e96ed6f7d3b25216121b5f89af9960576e6e47f8))
- *(deps)* Bump k8s.io/api from 0.27.3 to 0.27.4 - ([7ad9eb2](https://github.com/postfinance/kubenurse/commit/7ad9eb284f6707aa86119c1b1a3673599babb545))
- *(deps)* Bump k8s.io/client-go from 0.27.3 to 0.28.0 - ([7791489](https://github.com/postfinance/kubenurse/commit/7791489d427e94880dc7f8d51050e22a38b7a9fe))
- *(deps)* Bump k8s.io/api from 0.28.0 to 0.28.1 - ([ca5a74c](https://github.com/postfinance/kubenurse/commit/ca5a74c3826f644d8d4022bb0fd8e17ce07484fc))
- *(deps)* Bump k8s.io/client-go from 0.28.0 to 0.28.1 - ([52bfac3](https://github.com/postfinance/kubenurse/commit/52bfac3fa98bb0bdc476bef7f3668656d062cebb))
- *(deps)* Bump actions/checkout from 3 to 4 - ([21c103d](https://github.com/postfinance/kubenurse/commit/21c103d0f8f81dbf094a8bb645e3cc0051ac5d66))
- *(deps)* Bump k8s.io/client-go from 0.28.1 to 0.28.4 - ([eb3c96c](https://github.com/postfinance/kubenurse/commit/eb3c96c13817603edb69af60701b02b78cf896de))
- *(deps)* Bump actions/setup-go from 4 to 5 - ([b395623](https://github.com/postfinance/kubenurse/commit/b395623c305024412aec23cad91a6266f4436f1b))
- *(deps)* Bump helm/chart-releaser-action from 1.5.0 to 1.6.0 - ([efc98fa](https://github.com/postfinance/kubenurse/commit/efc98fa1d2501bb1d78255831aebb264add61b16))
- *(deps)* Bump docker/login-action from 2 to 3 - ([87f6111](https://github.com/postfinance/kubenurse/commit/87f611174c5ece346b464045da853ce1b0b54e99))
- *(dockerfile)* Update misconfigured maintainer label - ([461bda5](https://github.com/postfinance/kubenurse/commit/461bda530444775647846390338005adc1c96d43))
- HelmChart improvements - ([6e82de2](https://github.com/postfinance/kubenurse/commit/6e82de238d3181462b0be02843c29729277f10e2))




## New Contributors
* @clementnuss made their first contribution
* @luanaBanana made their first contribution in [#84](https://github.com/postfinance/kubenurse/pull/84)
* @ckittelmann made their first contribution## [1.7.1] - 2023-06-26

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/kubenurse-0.3.1...v1.7.1

### 📚 Documentation

- Add changlog - ([928cb07](https://github.com/postfinance/kubenurse/commit/928cb077a04e3da3690eda673b829055b15e5f85))
- Udpate changelog - ([515a7c4](https://github.com/postfinance/kubenurse/commit/515a7c41a85b301c36122f46fa5fec0758c0418e))

### ⚙️ Miscellaneous Tasks

- Update packages, CI actions and Go version. Fix linting. - ([88a900b](https://github.com/postfinance/kubenurse/commit/88a900bcf5ca8ce7fc3d6923231c8de28ba530a2))
- Enable dependabot - ([304b996](https://github.com/postfinance/kubenurse/commit/304b996b9d9adaa81ff7160d05e331c6d4ba037f))
- Update dependabot - ([85b19a6](https://github.com/postfinance/kubenurse/commit/85b19a66a66e155d9957c2a2903a71834fe22b90))
- Update dependabot commit message - ([1d445fe](https://github.com/postfinance/kubenurse/commit/1d445fe4c97bab318b09829e799503825bc131bf))
- Update .cc.yml - ([99c490a](https://github.com/postfinance/kubenurse/commit/99c490a361b8cd9f67ce08a4adc8becacfd610d4))
- Set dependabot interval to weekly - ([a9e53ae](https://github.com/postfinance/kubenurse/commit/a9e53ae3ea94342cea19f7981254d099ff704540))




## New Contributors
* @dependabot[bot] made their first contribution## [kubenurse-0.3.1] - 2022-12-07

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/kubenurse-0.3.0...kubenurse-0.3.1

### 🚀 Features

- Add new helm configurations - ([8c2e6c6](https://github.com/postfinance/kubenurse/commit/8c2e6c6526b737b79134c3b6b9235f363a1fce47))

### ⚙️ Miscellaneous Tasks

- *(helm)* Bump chart to 0.3.1 - ([ae27984](https://github.com/postfinance/kubenurse/commit/ae27984db9c5241ead166c1bd964da6c263f8a65))


## [kubenurse-0.3.0] - 2022-12-06

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.7.0...kubenurse-0.3.0

### 🚀 Features

- *(helm)* New configuration options (#57) - ([13484e6](https://github.com/postfinance/kubenurse/commit/13484e613a9e44b4d73f95b2e30ba54fa6cda7a1))

### ⚙️ Miscellaneous Tasks

- *(helm)* Bump chart to 0.3.0 - ([d5985e2](https://github.com/postfinance/kubenurse/commit/d5985e2c336e963b9f1d6f7bd99632e0e868ac61))


## [1.7.0] - 2022-11-01

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/kubenurse-0.2.1...v1.7.0

### 🚀 Features

- *(helm)* Make KUBENURSE_INSECURE configurable (#51) - ([4d4dc39](https://github.com/postfinance/kubenurse/commit/4d4dc3979bdd344f6f1f3ad746ea0249315b4b36))

### 🐛 Bug Fixes

- *(helm)* Chart should respect `-n <namespace>` flag (#53) - ([a5a3a79](https://github.com/postfinance/kubenurse/commit/a5a3a792a228302e30ea60845069f23faaeafb67))
- Use new ingress spefification (#52) - ([8b896f4](https://github.com/postfinance/kubenurse/commit/8b896f4c60339f503680443922211dfdf970d5d2))

### ⚙️ Miscellaneous Tasks

- *(helm)* Bump chart to 2.2 - ([c0c1db5](https://github.com/postfinance/kubenurse/commit/c0c1db523f6c5c9ea5049da4a825517d6227f5cb))




## New Contributors
* @myaser made their first contribution## [kubenurse-0.2.1] - 2022-10-25

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/kubenurse-0.2.0...kubenurse-0.2.1

### 🚀 Features

- *(helm)* Add support for volumes and volumeMounts (#49) - ([986d3dc](https://github.com/postfinance/kubenurse/commit/986d3dc96faa05458802f5cc1c8f118e19ba2fca))
- *(helm)* Add dnsConfig option (#50) - ([3fed269](https://github.com/postfinance/kubenurse/commit/3fed269073d45e6175fa1ef0be6977cb14cc575e))

### 🐛 Bug Fixes

- *(helm)* Parse error when using extraEnvs (#48) - ([3a56edb](https://github.com/postfinance/kubenurse/commit/3a56edbb5a6a080e5a589b9b0812c7bda14c94a4))

### 📚 Documentation

- Add reference to online helm repository - ([f04a6f7](https://github.com/postfinance/kubenurse/commit/f04a6f7a62c109b16081a7a26b8cc2836233085e))

### ⚙️ Miscellaneous Tasks

- *(helm)* Update chart version - ([5383aa7](https://github.com/postfinance/kubenurse/commit/5383aa776a83fb2c2482a1a44f178fbb332e9ea3))


## [kubenurse-0.2.0] - 2022-07-21

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.6.0-rc1...kubenurse-0.2.0

### 🚀 Features

- Implement helm chart releaser (#47) - ([7f52b47](https://github.com/postfinance/kubenurse/commit/7f52b474a6531cdd1b63c5358e6f51e67a02a317))

### 🐛 Bug Fixes

- Use current main branch naming for the helm releaser - ([4dd5ede](https://github.com/postfinance/kubenurse/commit/4dd5eded72dd83798f622da716cd28ddf4404b0c))

### ⚙️ Miscellaneous Tasks

- Update helm package version to 1.6.0 - ([e261007](https://github.com/postfinance/kubenurse/commit/e26100791c3e824867218b150d12b83162e115c4))


## [1.6.0-rc1] - 2022-06-03

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.5.2...v1.6.0-rc1

### 📚 Documentation

- Add example grafana dashboard (#40) - ([aa94ef8](https://github.com/postfinance/kubenurse/commit/aa94ef8e431995e6017c6bb9e3eee6ab47d1c92b))
- Add link to example Grafana dashboard - ([95ad678](https://github.com/postfinance/kubenurse/commit/95ad678f3078c0113c2d721b63333a442cb38bb5))

### ⚙️ Miscellaneous Tasks

- Split workflows and create initial CI setup with traefik (#39) - ([806e7c7](https://github.com/postfinance/kubenurse/commit/806e7c712e0869dc57c921e119054f1a67d4d62d))
- Update golangci-lint to v1.46 (#41) - ([797f3fb](https://github.com/postfinance/kubenurse/commit/797f3fba6edd6cb5e441e1efe23b1f2bcf63e1ab))
- Use example domains instead of assignable ones - ([94e7075](https://github.com/postfinance/kubenurse/commit/94e70751f20c4639a25fbc94b1dcf22a3c53cd01))
- Update dependencies (#43) - ([6b0761c](https://github.com/postfinance/kubenurse/commit/6b0761cadd27957043d121491a8d10aef8e430cc))




## New Contributors
* @domi2120 made their first contribution in [#42](https://github.com/postfinance/kubenurse/pull/42)## [1.5.2] - 2022-02-17

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.5.1...v1.5.2

### ⚙️ Miscellaneous Tasks

- Update go dependencies to use latest available stable versions (#36) - ([ca04845](https://github.com/postfinance/kubenurse/commit/ca048452c244730620302432d834b6739308cd9f))




## New Contributors
* @eli-halych made their first contribution in [#34](https://github.com/postfinance/kubenurse/pull/34)## [1.5.1] - 2022-01-21

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.5.0...v1.5.1

### 🐛 Bug Fixes

- Enforce timeouts in the kubenurse http.Server to avoid possible goroutine/memory leaks - ([d07df3b](https://github.com/postfinance/kubenurse/commit/d07df3bc86fdc12275e2a73ca8866cf81c2c5a29))


## [1.5.0] - 2022-01-17

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.5.0-beta1...v1.5.0

### 🚀 Features

- Expose metrics from the kubenurse httpclient (#31) - ([ebb0764](https://github.com/postfinance/kubenurse/commit/ebb076466723e71f50c0e721025b676923a7889d))

### 📚 Documentation

- Update README and fix some spelling/grammar mistakes (#30) - ([9f02d56](https://github.com/postfinance/kubenurse/commit/9f02d56be5031e6c8395f6b009a8f2a580d44010))


## [1.5.0-beta1] - 2022-01-05

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.4.1...v1.5.0-beta1

### 🚜 Refactor

- [**breaking**] Rewrite and cleanup kubenurse server code  (#29) - ([7beac30](https://github.com/postfinance/kubenurse/commit/7beac30751e8fd96093486fad9ada30f314e7dc4))


## [1.4.1] - 2021-09-30

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.4.0...v1.4.1

### 🐛 Bug Fixes

- *(examples)* Bump kubenurse version to v1.4.0 - ([6f1228c](https://github.com/postfinance/kubenurse/commit/6f1228c0975ad136ca51b75ac921f9f86cea3575))

### 📚 Documentation

- Update changelog (reference commits and PR) - ([ac79bfb](https://github.com/postfinance/kubenurse/commit/ac79bfbdcc38688c042074cbec0e142525e02ea7))

### ⚙️ Miscellaneous Tasks

- Update goreleaser config to newest version (0.178.0) - ([2f8cb96](https://github.com/postfinance/kubenurse/commit/2f8cb963fce6b6847d1a620a26ddb6a826a60c32))
- Fix ingress deployment in kind cluster - ([1d819ad](https://github.com/postfinance/kubenurse/commit/1d819adf329d84879789bc7c32459789ee7f1cbe))
- Updates for k8s v1.21.2 (#28) - ([a792cd8](https://github.com/postfinance/kubenurse/commit/a792cd8f7aaae4741d8922982c48f3aae580b938))
- Update changelog with 1.4.1 release - ([50fb9eb](https://github.com/postfinance/kubenurse/commit/50fb9eb8d05f544ccbe374b968a37e02c69e1849))




## New Contributors
* @bramvdklinkenberg made their first contribution## [1.4.0] - 2021-05-25

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.3.4...v1.4.0

### 🐛 Bug Fixes

- *(examples)* Bump kubenurse version to v1.3.4 - ([4e0a4c3](https://github.com/postfinance/kubenurse/commit/4e0a4c3359988f4baf2133a33e0c0f8567fe3e5d))

### 📚 Documentation

- Add changelog - ([ede0034](https://github.com/postfinance/kubenurse/commit/ede0034da45e5447964f4f9e172141d1e5175c32))
- Update changelog - ([044d105](https://github.com/postfinance/kubenurse/commit/044d10594e57523ab07f0de92aa5c908e922981a))

### ⚙️ Miscellaneous Tasks

- Update all go dependencies - ([c0df790](https://github.com/postfinance/kubenurse/commit/c0df7900e872d448f7aaa710a664502719432620))
- Update changelog with 1.4.0 release - ([47b6cfc](https://github.com/postfinance/kubenurse/commit/47b6cfc3a0393f40b92a0e91cb44036e94e3c8aa))




## New Contributors
* @zbindenren made their first contribution
* @matthyx made their first contribution in [#19](https://github.com/postfinance/kubenurse/pull/19)## [1.3.4] - 2021-04-20

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.3.3...v1.3.4

### 🐛 Bug Fixes

- *(discovery)* Prevent panic when checking for schedulable nodes only - ([2243226](https://github.com/postfinance/kubenurse/commit/2243226bd4031239c4d5f89afa6336e7bfd3c9fd))
- *(examples)* Bump kubenurse version to v1.3.3 - ([c13ebc1](https://github.com/postfinance/kubenurse/commit/c13ebc11f9a4612b13a81c2fef7dde1b71567c2e))

### CI

- Use latest-ci image for CI testing - ([eb11afb](https://github.com/postfinance/kubenurse/commit/eb11afb82af5ea5e3a455813a67f9834ae35c070))
- Use latest-ci image for CI testing - ([caa2105](https://github.com/postfinance/kubenurse/commit/caa21051da0f1a67eb1fcf1e0065b2f0a87888d1))


## [1.3.3] - 2021-04-20

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.3.2...v1.3.3

### 🚀 Features

- Flag to consider kubenurses on unschedulable nodes - ([cd9ac29](https://github.com/postfinance/kubenurse/commit/cd9ac29bfdec070a374ca44f41d3a03f466c8607))
- CI improvements and RBAC fixes - ([394daf1](https://github.com/postfinance/kubenurse/commit/394daf190c0813a38c0849c29aed63ea09ec4199))

### ⚙️ Miscellaneous Tasks

- Liniting - ([b99d08d](https://github.com/postfinance/kubenurse/commit/b99d08d32ededf9e5ddffbc56acf302fc3571d39))


## [1.3.2] - 2021-03-01

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.3.1...v1.3.2

### 📚 Documentation

- Add toleration example for master/control-plane - ([c5bfacb](https://github.com/postfinance/kubenurse/commit/c5bfacbf1d13493c313136186509f10dd8e16eb5))

### ⚙️ Miscellaneous Tasks

- Update dependencies - ([b1200a9](https://github.com/postfinance/kubenurse/commit/b1200a956560bdbf7e9158804800834daa1a8c92))


## [1.3.1] - 2020-12-09

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.3.0...v1.3.1

### 🐛 Bug Fixes

- Remove unwanted linter configuration - ([d928439](https://github.com/postfinance/kubenurse/commit/d92843948c2ab79bf4fb175181ab4e72b23b5a10))

### ⚙️ Miscellaneous Tasks

- Setup github actions, configure golangci-lint and fix lint errors - ([a4deaf8](https://github.com/postfinance/kubenurse/commit/a4deaf8c566cf9045a61a6df6fc25eb417dca5cf))


## [1.3.0] - 2020-12-09

**Full Changelog**: https://github.com/postfinance/kubenurse/compare/v1.2.0...v1.3.0

### 🚀 Features

- Exclude nodes which are not schedulable from neighbour checks - ([b6acb93](https://github.com/postfinance/kubenurse/commit/b6acb939004710141b08170d07dcfbe3db923347))

### ⚙️ Miscellaneous Tasks

- Update go dependencies - ([163433c](https://github.com/postfinance/kubenurse/commit/163433c214571ebcc859b3a06e4415b179f5325d))




## New Contributors
* @ made their first contribution## [1.0.0] - 2018-12-06



### ApiServerDNS

- Change name to fqdn - ([c127f98](https://github.com/postfinance/kubenurse/commit/c127f986e5762baa3fec345e0c20bab3dc480928))




## New Contributors
* @djboris9 made their first contribution
* @ghouscht made their first contribution
* @opensourcepf made their first contribution
