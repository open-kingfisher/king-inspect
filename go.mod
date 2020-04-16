module kingfisher/king-inspect

go 1.12

require (
	cloud.google.com/go v0.50.0 // indirect
	github.com/digitalocean/clusterlint v0.1.3
	github.com/docker/distribution v2.7.1+incompatible
	github.com/gin-gonic/gin v1.4.0
	github.com/gogo/protobuf v1.3.0 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20191227163750-53104e6ec876 // indirect
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20191228213918-04cbcbbfeed8 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
	k8s.io/api v0.0.0-20191206001707-7edad22604e1
	k8s.io/apimachinery v0.0.0-20191203211716-adc6f4cd9e7d
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/metrics v0.0.0-20190822063148-e60d8d0865eb
	kingfisher/kf v0.0.0-00010101000000-000000000000
)

replace (
	github.com/docker/docker => github.com/docker/docker v0.7.3-0.20190924004649-91870ed38213
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
	kingfisher/kf => ../kf
	kingfisher/king-k8s => ../king-k8s
)
