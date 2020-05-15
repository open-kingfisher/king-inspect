# Kingfisher Inspect
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-kingfisher/king-inspect)](https://goreportcard.com/report/github.com/open-kingfisher/king-inspect)

对Kubernetes集群进行健康扫描，以图表的方式进行展示、建议及说明

![image](screenshots/inspect.gif)

## 现有检查项目
基本检查 | 描述 
------------ | -------------
裸Pod | 避免在集群中使用裸Pod
完全合格的镜像名(FQIN) | 避免使用不完全合格的镜像名 
镜像Latest标签 | 避免使用latest标签
存活探针 | 建议为Pod创建存活探针 
就绪探针 | 建议为Pod创建就绪探针
默认命名空间 | 避免在default命名空间创建各种资源
资源配额(资源要求检测大于5核5G) | 建议配置Pod的资源请求、Pod的资源限制同时检测Pod资源请求是否过大（大于5核5G视为过大）
卷挂载(主机路径) | 避免挂载主机路径
节点自定义标签 | 避免自定义节点标签 
Metric Server | 建议集群安装Metric Server服务
Pod节点选择标签(节点名作为节点选择标签) | 避免Pod节点选择标签为节点名
准入控制Webhook(Validating Webhook 和 Mutating Webhook) | 避免配置的Validating Webhook针对的服务的命名空间不存在；避免配置的Validating Webhook针对的service不存在；配置的Mutating Webhook针对的服务的命名空间不存在；避免配置的Mutating Webhook针对的service不存在；避免配置的Validating Webhook针对的Namespace为kubernetes系统Namespace；避免已配置的Mutating Webhook针对的Namespace为kubernetes系统Namespace
命名空间资源限制 | 建议为命名空间添加资源限制（LimitRange）（未完成）
命名空间资源配额 | 建议为命名空间添加资源配额（ResourceQuota）（未完成）

无用检查 | 描述 
------------ | -------------
服务账户 | 无用的ServiceAccount
ConfigMap | 无用的ConfigMap
Secret | 无用的Secret
PV | 无用的PersistentVolume
PVC | 无用的PersistentVolumeClaim
HPA | 无用的HorizontalPodAutoscaler
集群角色 | 无用的ClusterRole
角色 | 无用的Role
服务 | 无用的Service
副本集 | 无用的ReplicaSet
Pod中断预算 | 无用的PodDisruptionBudget
Pod预设 | 无用的PodPreset

状态检查 | 描述 
------------ | -------------
节点状态 | 节点处于未知状态;节点未处于就绪状态;节点内存不足；节点硬盘空间不足；节点PID不足；节点的网络不可达
命名空间状态 | 命名空间未处于就绪状态
Pod状态 | Pod状态应该是Running或者Succeeded；Pod中容器重启次数为大于10
组件状态 | 组件controller-manager、scheduler、etcd状态

安全检查 | 描述 
------------ | -------------
特权容器 | 避免使用特权模式的容器
API Server | 参考：[CIS](https://github.com/open-kingfisher/king-inspect/blob/master/doc/cis/CIS_Kubernetes_Benchmark_v1.5.0.pdf)
Controller Manager| 参考：[CIS](https://github.com/open-kingfisher/king-inspect/blob/master/doc/cis/CIS_Kubernetes_Benchmark_v1.5.0.pdf)
Scheduler | 参考：[CIS](https://github.com/open-kingfisher/king-inspect/blob/master/doc/cis/CIS_Kubernetes_Benchmark_v1.5.0.pdf)
Etcd | 参考：[CIS](https://github.com/open-kingfisher/king-inspect/blob/master/doc/cis/CIS_Kubernetes_Benchmark_v1.5.0.pdf)

## 关闭资源的指定检查项
- 在资源的Annotations中添加如下注释，例如：禁止检查Pod状态和存活检测
```yaml
annotations: 
      kingfisher.inspect.com/disabled-inspect: "pod-state,liveness"
```
## 依赖

- Golang： `Go >= 1.13`

## 说明

- 安全审查基于[CIS](https://www.cisecurity.org/cis-benchmarks/) Kubernetes_Benchmark_v1.5.0
- CIS中文说明
  - [CIS kubernetes benchmarks Master Node 配置文件](https://zhuanlan.zhihu.com/p/138445189)
  - [CIS kubernetes benchmarks API Server 配置](https://zhuanlan.zhihu.com/p/138995000)
- 借鉴项目 [clusterlint](https://github.com/digitalocean/clusterlint)

## Makefile的使用

- 根据需求修改对应的REGISTRY变量，即可修改推送的仓库地址
- 编译成二进制文件： make build
- 生成镜像推送到镜像仓库： make push

## 联系我们
- [交流群](https://github.com/open-kingfisher/community/blob/master/contact_us/README.md)