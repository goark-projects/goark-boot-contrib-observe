package gbcobserve

import (
	"fmt"
	"strings"

	"goark.dev/observe"
	observesdk "goark.dev/observe-sdk"
)

// ProviderFactory 延迟创建应用级 Provider。
type ProviderFactory func(resource observe.Resource, sdkOptions []observesdk.Option) (observe.Provider, error)

// Option 定制可观测自动配置。
type Option func(*settings) error

// WithEnabled 显式覆盖环境中的启用状态。
func WithEnabled(enabled bool) Option {
	return func(settings *settings) error { settings.enabled = &enabled; return nil }
}

// WithProvider 使用外部创建的 Provider。
func WithProvider(provider observe.Provider) Option {
	return func(settings *settings) error {
		if provider == nil {
			return fmt.Errorf("gbc-observe: provider is nil")
		}
		settings.provider = provider
		return nil
	}
}

// WithProviderFactory 设置 Provider 工厂。
func WithProviderFactory(factory ProviderFactory) Option {
	return func(settings *settings) error {
		if factory == nil {
			return fmt.Errorf("gbc-observe: provider factory is nil")
		}
		settings.factory = factory
		return nil
	}
}

// WithServiceName 显式设置服务名称。
func WithServiceName(name string) Option {
	return func(settings *settings) error {
		name = strings.TrimSpace(name)
		if name == "" {
			return fmt.Errorf("gbc-observe: service name is empty")
		}
		settings.serviceName = &name
		return nil
	}
}

// WithSDKOptions 追加 SDK 运行时选项。
func WithSDKOptions(options ...observesdk.Option) Option {
	copied := append([]observesdk.Option(nil), options...)
	return func(settings *settings) error {
		settings.sdkOptions = append(settings.sdkOptions, copied...)
		return nil
	}
}

// WithExporters 追加 SDK 导出器。
func WithExporters(exporters ...observe.Exporter) Option {
	return WithSDKOptions(observesdk.WithExporters(exporters...))
}

// WithProcessors 追加 SDK 处理器。
func WithProcessors(processors ...observe.Processor) Option {
	return WithSDKOptions(observesdk.WithProcessors(processors...))
}
