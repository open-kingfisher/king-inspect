# Kingfisher Inspect
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-kingfisher/king-inspect)](https://goreportcard.com/report/github.com/open-kingfisher/king-inspect)

对Kubernetes集群进行健康扫描，以图表的方式进行展示

![image](screenshots/inspect.gif)

## 现有检查项目
基本检查 | 描述 
------------ | -------------
裸Pod | test
完全合格的镜像名(FQIN) | test 
镜像Latest标签 | test 
存活探针 | test 
就绪探针 | test
默认命名空间 | test 
资源配额(资源要求检测大于5核5G) | test 
卷挂载(主机路径) | test
节点自定义标签 | test 
Metric Server | test 
Pod节点选择标签(节点名作为节点选择标签) | test 
准入控制Webhook(Validating Webhook 和 Mutating Webhook) | test 

## 依赖

- Golang： `Go >= 1.13`

## 说明

- 安全审查基于[CIS](https://www.cisecurity.org/cis-benchmarks/) Kubernetes_Benchmark_v1.5.0
- 借鉴项目 [clusterlint](https://github.com/digitalocean/clusterlint)

## Makefile的使用

- 根据需求修改对应的REGISTRY变量，即可修改推送的仓库地址
- 编译成二进制文件： make build
- 生成镜像推送到镜像仓库： make push

