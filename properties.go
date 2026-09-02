package gbcobserve

const (
	// StarterID 是可观测自动配置的稳定标识。
	StarterID = "goark.boot.observe"
	// BeanNameProvider 是统一 observe.Provider Bean 名称。
	BeanNameProvider = "goark.observe.provider"
	// BeanNameLifecycle 是 Provider 生命周期适配器 Bean 名称。
	BeanNameLifecycle = "goark.observe.lifecycle"
)

const (
	// PropertyEnabled 控制是否启用 SDK Provider。
	PropertyEnabled = "goark.observe.enabled"
	// PropertyServiceName 设置服务名称。
	PropertyServiceName = "goark.observe.service.name"
	// PropertyServiceVersion 设置服务版本。
	PropertyServiceVersion = "goark.observe.service.version"
	// PropertyServiceNamespace 设置服务命名空间。
	PropertyServiceNamespace = "goark.observe.service.namespace"
	// PropertyServiceInstanceID 设置服务实例标识。
	PropertyServiceInstanceID = "goark.observe.service.instance-id"
	// PropertyEnvironment 设置部署环境。
	PropertyEnvironment = "goark.observe.service.environment"
	// PropertyMetricCardinalityLimit 设置单个指标最大 series 数量。
	PropertyMetricCardinalityLimit = "goark.observe.metrics.cardinality-limit"
)

const (
	DefaultEnabled                = true
	DefaultServiceName            = "goark"
	DefaultMetricCardinalityLimit = 2000
)
