package check

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/kubernetes"
	"sync"
	"time"
)

// InspectResult 最终输出结果
type InspectResult struct {
	Diagnostics []Diagnostic       `json:"diagnostics"`
	Summary     map[string]Summary `json:"summary"`
}

// Run 应用过滤器并并行运行生成的检查列表
func Run(ctx context.Context, client *kubernetes.Clientset, checkFilter CheckFilter, levelFilter LevelFilter, namespace []string) (*InspectResult, error) {
	objects, err := FetchObjects(client, namespace, ctx)
	if err != nil {
		return nil, err
	}

	all, err := checkFilter.FilterChecks()
	if err != nil {
		return nil, err
	}
	if len(all) == 0 {
		return nil, errors.New("no checks to run. please provided the right check name")
	}
	var diagnostics []Diagnostic
	var mu sync.Mutex
	var g errgroup.Group
	summary := make(map[string]Summary)
	for _, check := range all {
		check := check
		g.Go(func() error {
			// 记录开始时间
			start := time.Now()
			d, s, err := check.Run(objects)
			// 记录自开始经过的时间
			elapsed := time.Since(start)
			if err != nil {
				return err
			}
			// 加锁确保diagnostics不会应为多协程导致数据不正确
			mu.Lock()
			// 填写诊断的检查名称和组名。在这里这样做可以免除检查的需要，并确保它们是一致的
			for i := 0; i < len(d); i++ {
				d[i].Check = check.Name()
				d[i].Group = check.Groups()
			}
			diagnostics = append(diagnostics, d...)
			s.Duration = elapsed
			s.Group = check.Groups()
			summary[check.Name()] = s
			mu.Unlock()
			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	// 过滤namespace
	//diagnostics = filterNamespace(diagnostics, namespace)
	// 根据注释过滤诊断信息
	// To Do 过滤诊断信息在检测的时候就过滤而不是过滤最终的诊断信息
	//diagnostics, summary = filterEnabled(diagnostics, summary)
	// 根据级别过滤诊断信息
	diagnostics = filterSeverity(levelFilter.Severity, diagnostics)
	CheckResult := &InspectResult{
		Diagnostics: diagnostics,
		Summary:     summary,
	}
	return CheckResult, err
}

func filterEnabled(diagnostics []Diagnostic, summary map[string]Summary) ([]Diagnostic, map[string]Summary) {
	var ret []Diagnostic
	sum := summary
	for _, d := range diagnostics {
		if IsEnabled(d.Check, d.Object) {
			ret = append(ret, d)
		} else {
			// 减去总数信息
			if s, ok := summary[d.Check]; ok {
				if s.Total != 0 {
					s.Total--
				}
				if s.Issue != 0 {
					s.Issue--
				}
				if s.Error != 0 {
					s.Error--
				}
				if s.Warning != 0 {
					s.Warning--
				}
				if s.Suggestion != 0 {
					s.Suggestion--
				}
				sum[d.Check] = s
			}
		}
	}
	return ret, sum
}

func filterNamespace(diagnostics []Diagnostic, namespace []string) []Diagnostic {
	var ret []Diagnostic
	for _, d := range diagnostics {
		for _, n := range namespace {
			if d.Object.Namespace == n {
				ret = append(ret, d)
			}
		}
	}
	return ret
}

func filterSeverity(level []Severity, diagnostics []Diagnostic) []Diagnostic {
	if len(level) == 0 {
		return diagnostics
	}
	var ret []Diagnostic
	for _, d := range diagnostics {
		for _, l := range level {
			if d.Severity == l {
				ret = append(ret, d)
			}
		}
	}
	return ret
}
