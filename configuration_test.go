package gbcobserve_test

import (
	"context"
	"sync/atomic"
	"testing"

	"goark.dev/boot"
	gbcobserve "goark.dev/gbc-observe"
	"goark.dev/goark"
	"goark.dev/observe"
	observesdk "goark.dev/observe-sdk"
)

func TestAutoConfigureRegistersProviderAndClosesIt(t *testing.T) {
	t.Parallel()
	provider := &trackingProvider{Provider: observe.NoopProvider()}
	app, err := boot.Run(t.Context(), boot.WithAutoConfiguration(gbcobserve.AutoConfigure(gbcobserve.WithProvider(provider))))
	if err != nil {
		t.Fatalf("boot.Run: %v", err)
	}
	appContext, ok := app.Context()
	if !ok {
		t.Fatal("application context is unavailable")
	}
	resolved, err := goark.Get[observe.Provider](t.Context(), appContext, gbcobserve.BeanNameProvider)
	if err != nil {
		t.Fatalf("resolve provider: %v", err)
	}
	if resolved != provider {
		t.Fatalf("provider = %T, want supplied provider", resolved)
	}
	if err := app.Close(t.Context()); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if provider.flushes.Load() != 1 || provider.shutdowns.Load() != 1 {
		t.Fatalf("lifecycle calls = flush:%d shutdown:%d", provider.flushes.Load(), provider.shutdowns.Load())
	}
}

func TestAutoConfigureDisabledRegistersNoopProvider(t *testing.T) {
	t.Parallel()
	created := atomic.Bool{}
	app, err := boot.Run(t.Context(), boot.WithAutoConfiguration(gbcobserve.AutoConfigure(
		gbcobserve.WithEnabled(false),
		gbcobserve.WithProviderFactory(func(observe.Resource, []observesdk.Option) (observe.Provider, error) {
			created.Store(true)
			return observe.NoopProvider(), nil
		}),
	)))
	if err != nil {
		t.Fatalf("boot.Run: %v", err)
	}
	defer app.Close(context.Background())
	appContext, _ := app.Context()
	provider, err := goark.Get[observe.Provider](t.Context(), appContext, gbcobserve.BeanNameProvider)
	if err != nil {
		t.Fatalf("resolve provider: %v", err)
	}
	if created.Load() {
		t.Fatal("disabled configuration must not invoke provider factory")
	}
	_, span := provider.Tracer("test").Start(t.Context(), "noop")
	if span.IsRecording() {
		t.Fatal("disabled provider must return non-recording spans")
	}
}

type trackingProvider struct {
	observe.Provider
	flushes   atomic.Int32
	shutdowns atomic.Int32
}

func (p *trackingProvider) ForceFlush(context.Context) error { p.flushes.Add(1); return nil }
func (p *trackingProvider) Shutdown(context.Context) error   { p.shutdowns.Add(1); return nil }
