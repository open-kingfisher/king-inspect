package check

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// Diagnostic 封装了每次检查返回的信息
type Diagnostic struct {
	Check    string                  `json:"check"`
	Group    []string                `json:"group"`
	Severity Severity                `json:"severity"`
	Message  string                  `json:"message"`
	Kind     Kind                    `json:"kind"`
	Object   *metav1.ObjectMeta      `json:"object"`
	Owners   []metav1.OwnerReference `json:"owners"`
}

// 检查项摘要信息
type Summary struct {
	Total      int           `json:"total"`
	Issue      int           `json:"issue"`
	Error      int           `json:"error"`
	Warning    int           `json:"warning"`
	Suggestion int           `json:"suggestion"`
	Duration   time.Duration `json:"duration"`
	Group      []string      `json:"group"`
}

func (d Diagnostic) String() string {
	return fmt.Sprintf("[%s] %s/%s/%s: %s", d.Severity, d.Object.Namespace,
		d.Kind, d.Object.Name, d.Message)
}

// LevelFilter 过滤诊断级别
type LevelFilter struct {
	Severity []Severity
}

// 严重程度确定每个诊断的优先级别
type Severity string

// Kind 表示诊断所针对的k8s对象的类型
type Kind string

const (
	Error                          Severity = "error"
	Warning                        Severity = "warning"
	Suggestion                     Severity = "suggestion"
	Pod                            Kind     = "pod"
	PodTemplate                    Kind     = "pod template"
	PersistentVolumeClaim          Kind     = "persistent volume claim"
	ConfigMap                      Kind     = "config map"
	Secret                         Kind     = "secret"
	ServiceAccount                 Kind     = "service account"
	ClusterRole                    Kind     = "cluster role"
	Role                           Kind     = "role"
	PersistentVolume               Kind     = "persistent volume"
	PodDisruptionBudget            Kind     = "pod disruption budget"
	PodPreset                      Kind     = "pod preset"
	Service                        Kind     = "service"
	ValidatingWebhookConfiguration Kind     = "validating webhook configuration"
	MutatingWebhookConfiguration   Kind     = "mutating webhook configuration"
	Node                           Kind     = "node"
	MetricServer                   Kind     = "metric server"
	HPA                            Kind     = "hpa"
	ReplicaSet                     Kind     = "replica set"
	APIServer                      Kind     = "api server"
	ControllerManager              Kind     = "controller manager"
	Scheduler                      Kind     = "scheduler"
	ETCD                           Kind     = "etcd"
)

var Message = map[int]string{
	100: "避免在集群中使用裸Pod",
	101: "%s%s的镜像%s是畸形的镜像名",
	102: "%s%s的镜像%s没有使用用完全合格的镜像名",
	103: "避免使用主机路径作为挂载卷",
	104: "容器%s镜像名不能解析",
	105: "避免容器%s镜像使用latest标签",
	106: "容器%s建议添加存活探针",
	107: "避免使用默认命名空间",
	108: "如果节点被替换或升级，自定义节点标签将丢失",
	109: "容器%s建议添加就绪探针",
	110: "设置%s%s资源限制及资源要求以防止资源争用",
	111: "设置%s%s容器资源限制以防止资源争用",
	112: "%s%sCPU资源要求过高在资源紧张的情况下可能无法被调度",
	113: "%s%s内存资源要求过高在资源紧张的情况下可能无法被调度",
	114: "%s%sCPU和内存资源要求过高在资源紧张的情况下可能无法被调度",
	115: "已配置的Validating Webhook针对的服务的命名空间不存在",
	116: "已配置的Validating Webhook针对的service不存在",
	117: "已配置的Mutating Webhook针对的服务的命名空间不存在",
	118: "已配置的Mutating Webhook针对的service不存在",
	119: "避免使用节点的kubernetes.io/hostname标签作为Pod的节点选择标签，以防节点主机名修改而无法调度Pod",
	120: "集群Metric Server没有安装",
	121: "已配置的Validating Webhook针对的Namespace为kubernetes系统Namespace",
	122: "已配置的Mutating Webhook针对的Namespace为kubernetes系统Namespace",
	200: "没有使用的PV",
	201: "没有使用的PVC",
	202: "没有使用的Secret",
	203: "没有使用的ConfigMap",
	204: "没有使用的HPA",
	205: "副本集的拥有者不存在",
	206: "没有使用的服务帐户",
	207: "该服务帐户开启了自动载入API Server Token功能(使用此服务账户的Pod将自动载入API Server Token)",
	208: "引用了不存在的Secret, %s命名空间下%s不存在",
	209: "没有使用的集群角色",
	210: "没有使用的角色",
	211: "没有使用的Pod中断预算",
	212: "没有使用的Pod预设",
	213: "没有使用的服务",
	300: "节点处于未知状态",
	301: "节点未处于就绪状态",
	302: "节点内存不足",
	303: "节点硬盘空间不足",
	304: "节点PID不足",
	305: "节点节点的网络不可达",
	306: "命名空间未处于就绪状态，可以尝试使命名空间spec.finalizers:[]",
	307: "Pod状态为%s，Pod状态应该是Running或者Succeeded",
	308: "Pod中容器%s重启次数为%d，大于%d次",
	400: "禁用对API Server的匿名请求,建议添加--anonymous-auth=false参数",
	401: "不要使用基本身份验证，建议删除--basic-auth-file参数",
	402: "不要使用基于Token的基本身份验证，建议删除--token-auth-file参数",
	403: "使用https进行kubelet连接，建议添加--kubelet-https=true或者删除--kubelet-https参数",
	404: "启用基于证书的kubelet身份验证，建议添加--kubelet-client-certificate参数",
	405: "启用基于证书的kubelet身份验证，建议添加--kubelet-client-key参数",
	406: "在建立连接之前验证kubelet的证书，建议添加--kubelet-certificate-authority参数",
	407: "不要总是授权所有请求，建议--authorization-mode参数不被设置为AlwaysAllow",
	408: "限制kubelet节点只读取与其关联的对象，建议--authorization-mode参数包含Node",
	409: "基于角色的访问控制(RBAC)允许对不同实体可以在集群中的不同对象上执行的操作进行细粒度控制，建议--authorization-mode参数包含RBAC",
	410: "限制API Server接受请求的速度，使用EventRateLimit允许控制对API Server在给定时间片内接受的" +
		"事件数量施加限制（1.15才有的特性），建议--enable-admission-plugins参数包含EventRateLimit",
	411: "允许API Server接受所有请求，不过滤任何请求（1.13中已经弃用了AlwaysAdmit），它的行为相当于关闭所有的入口控制器，" +
		"建议--enable-admission-plugins参数不包含AlwaysAdmit",
	412: "强制每个新pod每次拉取所需的image。在多租户集群中，可以确保用户的私有image只能由具有凭据来提取它们的人使用。" +
		"没有这个允许控制策略，一旦一个image被拉到一个节点，来自任何用户的任何pod都可以通过已知的image名称来使用它，" +
		"而不需要对image所有者进行任何授权检查。当启用此插件时，总是在启动容器之前拉取image，这意味着需要有效的凭据。" +
		"建议--enable-admission-plugins参数包含AlwaysPullImages",
	413: "SecurityContextDeny值用于准入控制器可以拒绝使用了SecurityContext字段的Pod，此字段可以允许集群中的特权升级。" +
		"当集群中没有使用PodSecurityPolicy时，应该使用此值。假如PodSecurityPolicy值不存在建议--enable-admission-plugins参数包含ServiceAccount",
	414: "当您创建一个pod时，如果您没有指定一个服务帐户，它将自动分配相同名称空间中的默认服务帐户。" +
		"建议--disable-admission-plugins参数不包含ServiceAccount",
	415: "拒绝在正在终止的Namespace中创建对象。将准入控制策略设置为NamespaceLifecycle可以确保不能在不存在的Namespace中创建对象，" +
		"并且在Namespace终止时不会用于创建新对象。建议这样做，以加强名称空间终止过程的完整性，并确保新对象的可用性。" +
		"建议--disable-admission-plugins参数不包含NamespaceLifecycle",
	416: "拒绝创建与PodSecurityPolicy不匹配的Pod。PodSecurityPolicy是集群级别的资源，它控制Pod可以执行的操作和它能够访问的内容。" +
		"PodSecurityPolicy对象定义了一组pod必须在哪些条件下运行才能被系统接受。" +
		"PodSecurityPolicy由控制Pod可以访问的安全特性的设置和策略组成，因此必须使用这些设置和策略来控制Pod访问权限。" +
		"建议--enable-admission-plugins包含PodSecurityPolicy参数",
	417: "限制kubelet可以修改的节点和Pod对象。使用NodeRestriction插件可以确保kubelet被限制在它可以修改的节点和Pod对象中。" +
		"建议--enable-admission-plugins参数包含NodeRestriction",
	418: "不要绑定不安全的API Server。建议删除--insecure-bind-address参数",
	419: "不要绑定到不安全的端口。建议添加--insecure-port=0参数",
	420: "不要禁用安全端口。建议添加--insecure-port=6443参数",
	421: "如果不需要，禁用分析。概要分析允许识别特定的性能瓶颈。它生成大量的程序数据，这些数据可能被用来揭示系统和程序细节。" +
		"如果您没有遇到任何瓶颈，并且不需要使用分析器进行故障排除。建议添加--profiling=false参数",
	422: "在API Server上启用审计，并设置所需的审计日志路径。建议添加--audit-log-path=/var/log/apiserver/audit.log参数",
	423: "设置保留审计日志的天数。建议添加--audit-log-maxage=30参数",
	424: "保留10个或适当数量的旧审计日志文件。建议添加--audit-log-maxbackup=10参数",
	425: "在达到100 MB或更大时轮转审计日志文件。建议添加--audit-log-maxsize=100参数",
	500: "在pod终止时激活垃圾收集器。垃圾收集对于确保足够的资源可用性和避免性能和可用性下降非常重要。" +
		"在最坏的情况下，系统可能会崩溃或在很长一段时间内无法使用。当前的垃圾收集设置是12500个终止的pod，这可能太高了，您的系统无法承受。" +
		"根据系统资源和测试，选择适当的阈值来激活垃圾收集。建议添加--terminated-pod-gc-threshold=10参数",
	501: "如果不需要，禁用分析。概要分析允许识别特定的性能瓶颈。它生成大量的程序数据，这些数据可能被用来揭示系统和程序细节。" +
		"如果您没有遇到任何瓶颈，并且不需要使用分析器进行故障排除。建议添加--profiling=false参数",
	502: "为每个控制器使用单独的服务帐户凭据。控制器管理器在kube-system名称空间中为每个控制器创建一个服务帐户，为其生成凭据，" +
		"并使用该服务帐户凭据为每个控制器循环构建专用API客户端。将--use-service-account-credentials设置为true，" +
		"将使用单独的服务帐户凭证在控制器管理器中运行每个控制循环。" +
		"当与RBAC结合使用时，这将确保控制循环以执行预期任务所需的最低权限运行。建议添加--use-service-account-credentials=true参数",
	503: "显式地为控制器管理器上的服务帐户设置一个服务帐户私钥文件。建议添加--service-account-private-key-file=<filename>参数",
	504: "允许pod在建立连接之前验证API服务器的服务证书。在需要联系API服务器的pod中运行的进程必须验证API服务器的服务证书。" +
		"如果做不到这一点，就可能成为中间人攻击的目标。使用--root-ca-file参数为控制器管理器提供API服务器的服务证书的根证书，" +
		"允许控制器管理器将受信任的bundle注入pods，以便它们可以验证到API服务器的TLS连接。建议添加--root-ca-file=<path/to/file>参数",
	505: "RotateKubeletServerCertificate导致kubelet在引导其客户端凭据后请求服务证书，并在其现有凭据过期时轮换证书。 " +
		"这种自动定期轮换可确保不会因证书过期而造成停机，从而解决了CIA安全三合会中的可用性。注意：仅当您允许kubelet从API服务器获取其证书时，此建议才适用。 " +
		"如果您的kubelet证书来自外部授权/工具，那么您需要自己进行轮换。建议添加--feature-gates=RotateKubeletServerCertificate=true参数",
	506: "不要将控制器管理器服务绑定到非环回的不安全地址。Controller Manager API服务默认运行在端口10252/TCP上，用于健康信息和度量信息，不需要身份验证或加密。" +
		"因此，它应该只绑定到本地主机接口，以最小化集群的攻击。建议添加--bind-address=127.0.0.1参数",
	600: "如果不需要，禁用分析。概要分析允许识别特定的性能瓶颈。它生成大量的程序数据，这些数据可能被用来揭示系统和程序细节。" +
		"如果您没有遇到任何瓶颈，并且不需要使用分析器进行故障排除。建议添加--profiling=false参数",
	601: "不要将调度器服务绑定到非环回的不安全地址。Scheduler API服务默认运行在端口10251/TCP上，用于健康信息和度量信息，不需要身份验证或加密。" +
		"因此，它应该只绑定到本地主机接口，以最小化集群的攻击。建议添加--bind-address=127.0.0.1参数",
	700: "为etcd服务配置TLS加密。建议添加--cert-file=</path/to/ca-file>参数",
	701: "为etcd服务配置TLS加密。建议添加--key-file=</path/to/key-file>参数",
	702: "在etcd服务上启用客户端身份验证。建议添加--client-cert-auth=true参数",
	703: "不要为TLS使用自签名证书。建议--auto-tls不被设置为true",
	704: "etcd应该配置为使用TLS加密进行对等连接。建议添加--peer-client-file=</path/to/peer-cert-file>参数",
	705: "etcd应该配置为使用TLS加密进行对等连接。建议添加--peer-key-file=</path/to/peer-key-file>参数",
	706: "etcd应该配置为对等身份验证。建议添加--peer-client-cert-auth=true参数",
	707: "不要为TLS使用自签名证书。建议--peer-auto-tls不被设置为true",
	708: "etcd使用与Kubernetes不同的证书颁发机构。建议添加--trusted-ca-file=</path/to/ca-file>参数",
	800: "容器%s使用了特权模式，请求确保镜像的可靠性",
	801: "容器%s使用了root运行，请配置runAsNonRoot为true，或者指定runAsUser和runAsGroup",
}
