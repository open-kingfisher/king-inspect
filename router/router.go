package router

import (
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-inspect/impl"
	"github.com/open-kingfisher/king-utils/common"
	jwtAuth "github.com/open-kingfisher/king-utils/middleware/jwt"
	"net/http"
)

func SetupRouter(r *gin.Engine) *gin.Engine {
	// token
	//r.GET(common.KingfisherPath+"token", user.GetToken)

	//重新定义404
	r.NoRoute(NoRoute)

	authorize := r.Group("/", jwtAuth.JWTAuth())
	{
		// inspect
		authorize.GET(common.InspectPath+"inspect", impl.ListInspect)
		authorize.GET(common.InspectPath+"inspect/:name", impl.GetInspect)
		authorize.POST(common.InspectPath+"inspect", impl.CreateInspect)
		authorize.DELETE(common.InspectPath+"inspect/:name", impl.DeleteInspect)
		authorize.PUT(common.InspectPath+"inspect", impl.UpdateInspect)
		authorize.GET(common.InspectPath+"inspectAction/:name", impl.ActionInspect)
		authorize.GET(common.InspectPath+"inspectEvent/:name", impl.EventInspect)
		authorize.GET(common.InspectPath+"inspectTime/:name", impl.TimeInspect)
	}
	return r
}

// 重新定义404错误
func NoRoute(c *gin.Context) {
	responseData := common.ResponseData{Code: http.StatusNotFound, Msg: "404 Not Found"}
	c.JSON(http.StatusNotFound, responseData)
}
