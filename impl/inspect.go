package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-inspect/resource"
	"github.com/open-kingfisher/king-utils/common"
	"github.com/open-kingfisher/king-utils/common/access"
	"github.com/open-kingfisher/king-utils/common/handle"
	"github.com/open-kingfisher/king-utils/common/log"
	"net/http"
)

func GetInspect(c *gin.Context) {
	r := resource.InspectResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleGet(&r)
	c.JSON(responseData.Code, responseData)
}

func ListInspect(c *gin.Context) {
	r := resource.InspectResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleList(&r)
	c.JSON(responseData.Code, responseData)
}

func CreateInspect(c *gin.Context) {
	r := resource.InspectResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleCreate(&r, c)
	c.JSON(responseData.Code, responseData)
}

func DeleteInspect(c *gin.Context) {
	r := resource.InspectResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleDelete(&r)
	c.JSON(responseData.Code, responseData)
}

func UpdateInspect(c *gin.Context) {
	r := resource.InspectResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleUpdate(&r, c)
	c.JSON(responseData.Code, responseData)
}

func ActionInspect(c *gin.Context) {
	// 获取clientSet，如果失败直接返回错误
	clientSet, err := access.Access(c.Query("cluster"))
	responseData := handle.HandlerResponse(nil, err)
	if responseData.Code != http.StatusOK {
		log.Errorf("%s%s", common.K8SClientSetError, err)
		return
	}
	r := resource.InspectResource{
		Params: handle.GenerateCommonParams(c, clientSet),
	}
	responseData = HandleAction(&r)
	c.JSON(responseData.Code, responseData)
}

func EventInspect(c *gin.Context) {
	r := resource.InspectResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleEvent(&r)
	c.JSON(responseData.Code, responseData)
}

func TimeInspect(c *gin.Context) {
	r := resource.InspectResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleTime(&r)
	c.JSON(responseData.Code, responseData)
}
