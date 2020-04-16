package check

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kingfisher/king-inspect/util"
	"strings"
)

const checkAnnotation = "kingfisher.inspect.com/disabled-inspect"
const separator = ","

// Check is a check that can run on Kubernetes objects.
type Check interface {
	// Name 返回此检查的唯一名称
	Name() string
	// Groups 返回一个组名列表，此检查应是该列表的一部分。如果检查不属于任何组，则返回nil或空列表是有效的
	Groups() []string
	// Description 返回关于此检查功能的详细的可读的描述
	Description() string
	// 对一组Kubernetes对象运行此检查。它可以返回警告(低优先级问题)和错误(高优先级问题)，以及指示检查未能运行的错误值
	Run(*Objects) ([]Diagnostic, Summary, error)
}

// IsEnabled 检查对象注释，查看是否禁用了检查
func IsEnabled(name string, item *metav1.ObjectMeta) bool {
	annotations := item.GetAnnotations()
	if value, ok := annotations[checkAnnotation]; ok {
		disabledChecks := strings.Split(value, separator)
		if util.Contains(disabledChecks, name) {
			return false
		}
	}
	return true
}
