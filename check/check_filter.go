package check

import (
	"kingfisher/kf/common/log"
	"kingfisher/king-inspect/util"
)

// CheckFilter 存储在运行检查时需要检查的名称和组名称
type CheckFilter struct {
	IncludeGroups []string
	IncludeChecks []string
}

// NewCheckFilter CheckFilter 的构造函数
func NewCheckFilter(includeGroups, includeChecks []string) (CheckFilter, error) {
	return CheckFilter{
		IncludeGroups: includeGroups,
		IncludeChecks: includeChecks,
	}, nil
}

// FilterChecks 根据检查名返回检查接口
func (c CheckFilter) FilterChecks() ([]Check, error) {
	all := make([]Check, 0)
	for _, checkName := range c.IncludeChecks {
		if check, err := Get(checkName); err != nil {
			log.Errorf("filter checks error:%s", err)
		} else {
			all = append(all, check)
		}
	}

	return all, nil
}

func (c CheckFilter) filterGroups() ([]Check, error) {
	if len(c.IncludeGroups) > 0 {
		groups, err := GetGroups(c.IncludeGroups)
		return groups, err
	}
	return List(), nil
}

func getChecksNotInGroups(groups []string) []Check {
	allGroups := ListGroups()
	var ret []Check
	for _, group := range allGroups {
		if !util.Contains(groups, group) {
			ret = append(ret, GetGroup(group)...)
		}
	}
	return ret
}
