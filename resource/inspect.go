package resource

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/open-kingfisher/king-inspect/check"
	"github.com/open-kingfisher/king-utils/common"
	"github.com/open-kingfisher/king-utils/common/handle"
	"github.com/open-kingfisher/king-utils/common/log"
	"github.com/open-kingfisher/king-utils/db"
	"github.com/open-kingfisher/king-utils/kit"
	"strconv"
	"time"
)

type InspectResource struct {
	Params   *handle.Resources
	PostData *common.InspectDB
}

func (r *InspectResource) Get() (interface{}, error) {
	inspect := common.InspectDB{}
	if err := db.GetById(common.InspectTable, r.Params.Name, &inspect); err != nil {
		return nil, err
	}
	return inspect, nil
}

func (r *InspectResource) List() (interface{}, error) {
	inspect := make([]common.InspectDB, 0)
	if err := db.List(common.DataField, common.InspectTable, &inspect, "WHERE data-> '$.cluster'=?", r.Params.Cluster); err != nil {
		return nil, err
	}
	return inspect, nil
}

func (r *InspectResource) Create(c *gin.Context) (err error) {
	inspect := common.InspectDB{}
	if err = c.BindJSON(&inspect); err != nil {
		return err
	}
	r.PostData = &inspect
	// 对提交的数据进行校验
	if err = c.ShouldBindWith(r.PostData, binding.Query); err != nil {
		return err
	}
	inspectList := make([]*common.InspectDB, 0)
	if err = db.List(common.DataField, common.InspectTable, &inspectList, "WHERE data-> '$.name'=? and data-> '$.cluster'=?", r.PostData.Name, r.Params.Cluster); err == nil {
		if len(inspectList) > 0 {
			return errors.New("the Inspect name already exists")
		}
	} else {
		return
	}
	r.PostData.Id = kit.UUID("i")
	r.PostData.CreateTime = time.Now().Unix()
	r.PostData.ModifyTime = time.Now().Unix()
	r.PostData.User = r.Params.User.Name
	r.PostData.Cluster = r.Params.Cluster
	if err = db.Insert(common.InspectTable, r.PostData); err != nil {
		log.Errorf("Inspect create error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return err
	}
	auditLog := handle.AuditLog{
		Kind:       common.Inspect,
		ActionType: common.Create,
		Resources:  r.Params,
		PostData:   r.PostData,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *InspectResource) Delete() (err error) {
	if err = db.Delete(common.InspectTable, r.Params.Name); err != nil {
		return
	}
	if err = db.Delete(common.InspectInfoTable, r.Params.Name); err != nil {
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Inspect,
		ActionType: common.Delete,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *InspectResource) Update(c *gin.Context) (err error) {
	inspect := common.InspectDB{}
	if err = c.BindJSON(&inspect); err != nil {
		return err
	}
	r.PostData = &inspect
	// 对提交的数据进行校验
	if err = c.ShouldBindWith(r.PostData, binding.Query); err != nil {
		return err
	}
	inspectList := make([]*common.InspectDB, 0)
	if err = db.List(common.DataField, common.InspectTable, &inspectList, "WHERE data-> '$.name'=? and data-> '$.cluster'=?", r.PostData.Name, r.Params.Cluster); err == nil {
		if len(inspectList) > 0 {
			for _, v := range inspectList {
				if v.Id != r.PostData.Id {
					return errors.New("the Inspect name already exists")
				}
			}
		}
	} else {
		return
	}
	inspects := common.InspectDB{}
	if err = db.GetById(common.InspectTable, r.PostData.Id, &inspects); err != nil {
		log.Errorf("Inspect update error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return
	}
	r.PostData.CreateTime = inspects.CreateTime
	r.PostData.ModifyTime = time.Now().Unix()
	r.PostData.Cluster = r.Params.Cluster
	if err = db.Update(common.InspectTable, r.PostData.Id, r.PostData); err != nil {
		log.Errorf("Inspect update error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Inspect,
		ActionType: common.Update,
		Resources:  r.Params,
		PostData:   r.PostData,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *InspectResource) Action() (interface{}, error) {
	inspect := common.InspectDB{}
	if err := db.GetById(common.InspectTable, r.Params.Name, &inspect); err != nil {
		log.Errorf("get inspect by db error:%s", err)
		return nil, err
	}
	allCheck := make([]string, 0)
	allCheck = append(allCheck, inspect.Basic...)
	allCheck = append(allCheck, inspect.Unused...)
	allCheck = append(allCheck, inspect.State...)
	allCheck = append(allCheck, inspect.Security...)
	filter, err := check.NewCheckFilter([]string{}, allCheck)
	if err != nil {
		return nil, err
	}
	level := make([]check.Severity, 0)
	for _, l := range inspect.Level {
		level = append(level, check.Severity(l))
	}
	levelFilter := check.LevelFilter{Severity: level}
	output, err := check.Run(context.Background(), r.Params.ClientSet, filter, levelFilter, inspect.Namespace)
	if err != nil {
		return nil, err
	}
	inspectInfo := common.InspectInfoDB{}
	inspectInfo.Id = r.Params.Name
	inspectInfo.CreateTime = time.Now().Unix()
	inspectInfo.Info = output
	if err = db.Insert(common.InspectInfoTable, inspectInfo); err != nil {
		log.Errorf("Inspect info insert error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return nil, err
	}
	return nil, nil
}

func (r *InspectResource) Event() (interface{}, error) {
	inspectInfoList := make([]common.InspectInfoDB, 0)
	if r.Params.Time == "" {
		if err := db.List(common.DataField, common.InspectInfoTable, &inspectInfoList, "WHERE data-> '$.id'=? order by data -> '$.createTime' desc limit 1", r.Params.Name); err != nil {
			return nil, err
		}
	} else {
		if createTime, err := strconv.ParseInt(r.Params.Time, 10, 64); err != nil {
			return nil, err
		} else {
			if err := db.List(common.DataField, common.InspectInfoTable, &inspectInfoList, "WHERE data-> '$.id'=? and data-> '$.createTime'=?", r.Params.Name, createTime); err != nil {
				return nil, err
			}
		}
	}
	return inspectInfoList, nil
}

func (r *InspectResource) Time() (interface{}, error) {
	createTime := make([]int64, 0)
	if err := db.List("data -> '$.createTime'", common.InspectInfoTable, &createTime, "WHERE data-> '$.id'=? order by data -> '$.createTime' desc limit 10", r.Params.Name); err != nil {
		return nil, err
	}
	return createTime, nil
}
