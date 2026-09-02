package gbcobserve

import (
	"fmt"
	"strings"

	coreenv "goark.dev/goark/core/env"
	"goark.dev/observe"
	observesdk "goark.dev/observe-sdk"
)

type settings struct {
	enabled           *bool
	serviceName       *string
	serviceVersion    string
	serviceNamespace  string
	serviceInstanceID string
	environment       string
	cardinalityLimit  int
	provider          observe.Provider
	factory           ProviderFactory
	sdkOptions        []observesdk.Option
}

func newSettings(environment coreenv.Environment, options []Option) (settings, error) {
	resolved := settings{cardinalityLimit: DefaultMetricCardinalityLimit, factory: defaultProviderFactory}
	if environment != nil {
		enabled, err := coreenv.ResolveValueAs[bool](environment, "${"+PropertyEnabled+":true}")
		if err != nil {
			return settings{}, err
		}
		resolved.enabled = &enabled
		name, err := coreenv.ResolveValueAs[string](environment, "${"+PropertyServiceName+":"+DefaultServiceName+"}")
		if err != nil {
			return settings{}, err
		}
		resolved.serviceName = &name
		resolved.serviceVersion = property(environment, PropertyServiceVersion)
		resolved.serviceNamespace = property(environment, PropertyServiceNamespace)
		resolved.serviceInstanceID = property(environment, PropertyServiceInstanceID)
		resolved.environment = property(environment, PropertyEnvironment)
		limit, err := coreenv.ResolveValueAs[int](environment, "${"+PropertyMetricCardinalityLimit+":2000}")
		if err != nil {
			return settings{}, err
		}
		resolved.cardinalityLimit = limit
	}
	for _, option := range options {
		if option != nil {
			if err := option(&resolved); err != nil {
				return settings{}, err
			}
		}
	}
	if resolved.enabled == nil {
		enabled := DefaultEnabled
		resolved.enabled = &enabled
	}
	if resolved.serviceName == nil {
		name := DefaultServiceName
		resolved.serviceName = &name
	}
	*resolved.serviceName = strings.TrimSpace(*resolved.serviceName)
	if *resolved.serviceName == "" {
		return settings{}, fmt.Errorf("gbc-observe: service name is empty")
	}
	if resolved.cardinalityLimit <= 0 {
		return settings{}, fmt.Errorf("gbc-observe: metric cardinality limit must be positive")
	}
	return resolved, nil
}

func (s settings) resource() observe.Resource {
	return observe.NewResource(*s.serviceName, observe.WithResourceVersion(s.serviceVersion), observe.WithResourceNamespace(s.serviceNamespace), observe.WithResourceInstanceID(s.serviceInstanceID), observe.WithResourceEnv(s.environment))
}
func property(environment coreenv.Environment, key string) string {
	value, _ := environment.GetProperty(key)
	return strings.TrimSpace(value)
}
func defaultProviderFactory(resource observe.Resource, options []observesdk.Option) (observe.Provider, error) {
	combined := make([]observesdk.Option, 0, len(options)+1)
	combined = append(combined, observesdk.WithResource(resource))
	combined = append(combined, options...)
	return observesdk.NewProvider(combined...)
}
