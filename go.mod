module github.com/open-kingfisher/king-inspect

go 1.14

require (
	github.com/docker/distribution v2.7.1+incompatible
	github.com/gin-gonic/gin v1.6.2
	github.com/open-kingfisher/king-utils v0.0.0-20200422073733-6505a8c88560
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/metrics v0.18.2
)

replace (
	k8s.io/api => k8s.io/api v0.17.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.3
	k8s.io/client-go => k8s.io/client-go v0.17.3
	k8s.io/metrics => k8s.io/metrics v0.17.3
)