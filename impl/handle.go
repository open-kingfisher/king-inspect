package impl

import (
	"github.com/gin-gonic/gin"
	"kingfisher/kf/common"
	"kingfisher/kf/common/handle"
)

type HandleGetInterface interface {
	Get() (interface{}, error)
}

type HandleListInterface interface {
	List() (interface{}, error)
}

type HandleDeleteInterface interface {
	Delete() error
}

type HandleUpdateInterface interface {
	Update(c *gin.Context) error
}

type HandleCreateInterface interface {
	Create(c *gin.Context) error
}

type HandleActionInterface interface {
	Action() (interface{}, error)
}

type HandleEventInterface interface {
	Event() (interface{}, error)
}

type HandleTimeInterface interface {
	Time() (interface{}, error)
}

func HandleGet(r HandleGetInterface) *common.ResponseData {
	responseData := handle.HandlerResponse(r.Get())
	return responseData
}

func HandleList(r HandleListInterface) *common.ResponseData {
	responseData := handle.HandlerResponse(r.List())
	return responseData
}

func HandleDelete(r HandleDeleteInterface) *common.ResponseData {
	responseData := handle.HandlerResponse(nil, r.Delete())
	return responseData
}

func HandleUpdate(r HandleUpdateInterface, c *gin.Context) *common.ResponseData {
	responseData := handle.HandlerResponse(nil, r.Update(c))
	return responseData
}

func HandleCreate(r HandleCreateInterface, c *gin.Context) *common.ResponseData {
	responseData := handle.HandlerResponse(nil, r.Create(c))
	return responseData
}

func HandleAction(r HandleActionInterface) *common.ResponseData {
	responseData := handle.HandlerResponse(r.Action())
	return responseData
}

func HandleEvent(r HandleEventInterface) *common.ResponseData {
	responseData := handle.HandlerResponse(r.Event())
	return responseData
}

func HandleTime(r HandleTimeInterface) *common.ResponseData {
	responseData := handle.HandlerResponse(r.Time())
	return responseData
}
