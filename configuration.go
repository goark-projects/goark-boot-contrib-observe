package gbcobserve

import (
	"context"

	"goark.dev/boot"
	goarkcontainer "goark.dev/goark/container"
	appcontext "goark.dev/goark/context"
	"goark.dev/observe"
	observesdk "goark.dev/observe-sdk"
)

// AutoConfigure 创建可观测自动配置。
func AutoConfigure(options ...Option) boot.AutoConfiguration {
	copied := append([]Option(nil), options...)
	return boot.NewAutoConfiguration(StarterID, func(_ context.Context, app *appcontext.ApplicationContext) error {
		return app.RegisterConfiguration(configuration{options: copied})
	})
}

type configuration struct{ options []Option }

func (configuration) Name() string { return StarterID + ".configuration" }
func (configuration) Order() int   { return -10000 }
func (c configuration) Register(ctx context.Context, registry *goarkcontainer.Registry) error {
	return c.RegisterWithContext(ctx, appcontext.NewConfigurationContext(nil, registry))
}
func (c configuration) RegisterWithContext(_ context.Context, config appcontext.ConfigurationContext) error {
	resolved, err := newSettings(config.Environment(), c.options)
	if err != nil {
		return err
	}
	if err := registerProvider(config.Registry(), resolved); err != nil {
		return err
	}
	return goarkcontainer.Register[*providerLifecycle](config.Registry(), BeanNameLifecycle, func(ctx context.Context, resolver goarkcontainer.Resolver) (*providerLifecycle, error) {
		provider, err := goarkcontainer.Get[observe.Provider](ctx, resolver, BeanNameProvider)
		if err != nil {
			return nil, err
		}
		return &providerLifecycle{provider: provider}, nil
	}, goarkcontainer.WithDependsOn(BeanNameProvider))
}

func registerProvider(registry *goarkcontainer.Registry, resolved settings) error {
	return goarkcontainer.Register[observe.Provider](registry, BeanNameProvider, func(ctx context.Context, resolver goarkcontainer.Resolver) (observe.Provider, error) {
		if !*resolved.enabled {
			return observe.NoopProvider(), nil
		}
		if resolved.provider != nil {
			return resolved.provider, nil
		}
		options := append([]observesdk.Option(nil), resolved.sdkOptions...)
		exporters, err := goarkcontainer.GetAllByType[observe.Exporter](ctx, resolver)
		if err != nil {
			return nil, err
		}
		if len(exporters) > 0 {
			options = append(options, observesdk.WithExporters(exporters...))
		}
		options = append(options, observesdk.WithMetricCardinalityLimit(resolved.cardinalityLimit))
		return resolved.factory(resolved.resource(), options)
	}, goarkcontainer.WithPrimary())
}
