module github.com/realshuting/multi-tenancy/benchmarks

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Microsoft/go-winio v0.4.14
	github.com/beorn7/perks v1.0.1
	github.com/blang/semver v3.5.0+incompatible
	github.com/coreos/etcd v3.3.15+incompatible
	github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a
	github.com/coreos/pkg v0.0.0-20180108230652-97fdf19511ea
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c
	github.com/evanphx/json-patch v4.5.0+incompatible
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gogo/protobuf v1.3.1
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9
	github.com/golang/protobuf v1.3.2
	github.com/google/go-cmp v0.3.1
	github.com/google/gofuzz v1.0.0
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.3.1
	github.com/hashicorp/golang-lru v0.5.3
	github.com/hpcloud/tail v1.0.0
	github.com/imdario/mergo v0.3.8
	github.com/json-iterator/go v1.1.8
	github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 v1.0.1
	github.com/onsi/ginkgo v1.10.3
	github.com/onsi/gomega v1.7.1
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4
	github.com/prometheus/common v0.7.0
	github.com/prometheus/procfs v0.0.5
	github.com/realshuting/multi-tenancy-benchmarks v0.0.0-20191028041724-a27250830445
	github.com/spf13/afero v1.2.2
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	go.uber.org/atomic v1.5.0
	go.uber.org/multierr v1.3.0
	go.uber.org/tools v0.0.0-20190618225709-2cfd321de3ee
	go.uber.org/zap v1.12.0
	golang.org/x/crypto v0.0.0-20191029031824-8986dd9e96cf
	golang.org/x/lint v0.0.0-20190930215403-16217165b5de
	golang.org/x/net v0.0.0-20191028085509-fe3aa8a45271
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sys v0.0.0-20191029155521-f43be2a4598c
	golang.org/x/text v0.3.2
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/tools v0.0.0-20191031160344-02d0efc0fb61
	google.golang.org/appengine v1.6.5
	google.golang.org/genproto v0.0.0-20191028173616-919d9bdd9fe6
	google.golang.org/grpc v1.24.0
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/inf.v0 v0.9.1
	gopkg.in/square/go-jose.v2 v2.4.0
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
	gopkg.in/yaml.v2 v2.2.4
	honnef.co/go/tools v0.0.1-2019.2.3
	k8s.io/api v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/client-go v0.0.0
	k8s.io/component-base v0.0.0
	k8s.io/kubernetes v1.16.2
	sigs.k8s.io/yaml v1.1.0
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20191003000013-35e20aa79eb8
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20191028232452-c47e10e6d5a3
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191030190112-bb31b70367b7
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20191003001037-3c8b233e046c
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20191016114015-74ad18325ed5
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191028230319-1a481fb1e32d
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20191025232453-66dd06a864dd
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20191016115129-c07a134afb42
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20191004115455-8e001e5d1894
	k8s.io/component-base => k8s.io/component-base v0.0.0-20191029070825-5e0e35147053
	k8s.io/cri-api => k8s.io/cri-api v0.0.0-20191025232916-446748cffdda
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20191003003551-0eecdcdcc049
	k8s.io/klog => k8s.io/klog v1.0.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20191016112429-9587704a8ad4
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20191016114939-2b2b218dc1df
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190918143330-0270cf2f1c1d
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20191016114407-2e83b6f20229
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20191016114748-65049c67a58b
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20191031072635-2ba9448df4cc
	k8s.io/kubelet => k8s.io/kubelet v0.0.0-20191016114556-7841ed97f1b2
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20191016115753-cf0698c3a16b
	k8s.io/metrics => k8s.io/metrics v0.0.0-20191016113814-3b1a734dba6e
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20191016112829-06bb3c9d77c9
	k8s.io/utils => k8s.io/utils v0.0.0-20191030222137-2b95a09bc58d
)
