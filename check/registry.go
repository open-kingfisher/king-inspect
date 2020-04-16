package check

import (
	"errors"
	"fmt"
	"sync"
)

type checkRegistry struct {
	mu     sync.RWMutex
	checks map[string]Check
	groups map[string][]Check
}

var (
	registry *checkRegistry
	initOnce sync.Once
)

// 注册检查项，每个检查使用init()进行注册
func Register(check Check) error {
	initOnce.Do(func() {
		registry = &checkRegistry{
			checks: make(map[string]Check),
			groups: make(map[string][]Check),
		}
	})

	registry.mu.Lock()
	defer registry.mu.Unlock()

	name := check.Name()
	if name == "" {
		return errors.New("checks must have non-empty names")
	}
	if _, ok := registry.checks[name]; ok {
		return fmt.Errorf("check named %q already exists", name)
	}
	registry.checks[name] = check
	for _, group := range check.Groups() {
		registry.groups[group] = append(registry.groups[group], check)
	}

	return nil
}

// 列出所有注册的检查
func List() []Check {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	ret := make([]Check, 0, len(registry.checks))
	for _, check := range registry.checks {
		ret = append(ret, check)
	}

	return ret
}

// 列出所有注册的组
func ListGroups() []string {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	ret := make([]string, 0, len(registry.groups))
	for group := range registry.groups {
		ret = append(ret, group)
	}

	return ret
}

// GetGroup 返回特定组中的检查
func GetGroup(name string) []Check {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	ret := make([]Check, 0, len(registry.groups[name]))
	for _, check := range registry.groups[name] {
		ret = append(ret, check)
	}

	return ret
}

// GetGroups 返回属于返回多个组的检查
func GetGroups(groups []string) ([]Check, error) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	var ret []Check
	for _, group := range groups {
		if checks, ok := registry.groups[group]; ok {
			ret = append(ret, checks...)
		} else {
			return nil, fmt.Errorf("group %s not found", group)
		}

	}

	return ret, nil
}

// Get从registry中检索指定的检查名如果为空抛出错误
func Get(name string) (Check, error) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	if registry.checks[name] != nil {
		return registry.checks[name], nil
	}
	return nil, fmt.Errorf("check not found: %s", name)
}
