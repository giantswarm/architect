module github.com/giantswarm/architect

go 1.15

require (
	github.com/giantswarm/app/v4 v4.13.0
	github.com/giantswarm/gitrepo v0.2.2
	github.com/giantswarm/microerror v0.4.0
	github.com/google/go-cmp v0.5.7
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/spf13/afero v1.8.1
	github.com/spf13/cobra v1.4.0
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	sigs.k8s.io/yaml v1.3.0
)

replace (
	github.com/coreos/etcd v3.3.10+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/etcd v3.3.13+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
)
