module github.com/giantswarm/architect

go 1.15

require (
	github.com/giantswarm/app/v4 v4.13.0
	github.com/giantswarm/argoapp v0.1.4
	github.com/giantswarm/gitrepo v0.2.2
	github.com/giantswarm/microerror v0.3.0
	github.com/google/go-cmp v0.5.6
	github.com/spf13/afero v1.7.0
	github.com/spf13/cobra v1.1.3
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/coreos/etcd v3.3.10+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/etcd v3.3.13+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
)
