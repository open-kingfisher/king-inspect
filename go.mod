module open-kingfisher/king-inspect

go 1.14

require (
	github.com/docker/distribution v2.7.1+incompatible
	github.com/gin-gonic/gin v1.6.2
	github.com/stretchr/testify v1.5.1
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
	k8s.io/api v0.18.1
	k8s.io/apimachinery v0.18.1
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/metrics v0.18.1
	kingfisher/kf v0.0.0-00010101000000-000000000000
	kingfisher/king-inspect v0.0.0-00010101000000-000000000000
)

replace (
	kingfisher/kf => ../kf
	kingfisher/king-inspect => ../king-inspect
)
