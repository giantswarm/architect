module github.com/giantswarm/architect

go 1.19

require (
	github.com/giantswarm/app/v4 v4.13.0
	github.com/giantswarm/gitrepo v0.2.2
	github.com/giantswarm/microerror v0.4.0
	github.com/google/go-cmp v0.6.0
	github.com/spf13/afero v1.10.0
	github.com/spf13/cobra v1.7.0
	sigs.k8s.io/yaml v1.4.0
)

require (
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/giantswarm/apiextensions/v3 v3.22.0 // indirect
	github.com/giantswarm/k8smetadata v0.3.0 // indirect
	github.com/go-errors/errors v1.0.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20190725054713-01f96b0aa0cd // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/sergi/go-diff v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/src-d/gcfg v1.4.0 // indirect
	github.com/xanzy/ssh-agent v0.2.1 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.15.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/src-d/go-billy.v4 v4.3.2 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/apiextensions-apiserver v0.27.0 // indirect
	k8s.io/apimachinery v0.27.0 // indirect
	k8s.io/klog/v2 v2.90.1 // indirect
	k8s.io/utils v0.0.0-20230209194617-a36077c30491 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)

replace (
	github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.46.6
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	github.com/miekg/dns => github.com/miekg/dns v1.1.56
	github.com/nats-io/jwt => github.com/nats-io/jwt/v2 v2.2.0
	github.com/nats-io/nats-server/v2 => github.com/nats-io/nats-server/v2 v2.10.4
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v1.17.0
	go.mongodb.org/mongo-driver => go.mongodb.org/mongo-driver v1.12.1

	golang.org/x/net => golang.org/x/net v0.17.0
)
