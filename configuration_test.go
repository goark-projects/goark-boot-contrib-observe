package gbcobserve_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"goark.dev/boot"
	gbcobserve "goark.dev/gbc-observe"
	"goark.dev/goark"
	goarkcontainer "goark.dev/goark/container"
	appcontext "goark.dev/goark/context"
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

func TestAutoConfigureCollectsExporterBeans(t *testing.T) {
	t.Parallel()
	exporter := &captureExporter{}
	app, err := boot.Run(t.Context(),
		boot.WithAutoConfiguration(gbcobserve.AutoConfigure()),
		boot.WithConfiguration(exporterConfiguration{exporter: exporter}),
	)
	if err != nil {
		t.Fatalf("boot.Run: %v", err)
	}
	appContext, _ := app.Context()
	provider := goark.MustGet[observe.Provider](t.Context(), appContext, gbcobserve.BeanNameProvider)
	_, span := provider.Tracer("test").Start(t.Context(), "registered-exporter")
	span.End()
	if exporter.count() != 1 {
		t.Fatalf("exported spans = %d, want 1", exporter.count())
	}
	if err := app.Close(t.Context()); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

type exporterConfiguration struct{ exporter observe.Exporter }

func (exporterConfiguration) Name() string { return "test.exporter.configuration" }
func (exporterConfiguration) Order() int   { return -20000 }
func (c exporterConfiguration) Register(ctx context.Context, registry *goarkcontainer.Registry) error {
	return c.RegisterWithContext(ctx, appcontext.NewConfigurationContext(nil, registry))
}
func (c exporterConfiguration) RegisterWithContext(_ context.Context, config appcontext.ConfigurationContext) error {
	return goarkcontainer.RegisterInstance[observe.Exporter](config.Registry(), "test.exporter", c.exporter, goarkcontainer.WithLazy())
}

type captureExporter struct {
	mu    sync.Mutex
	spans []observe.SpanSnapshot
}

func (*captureExporter) Descriptor() observe.ExporterDescriptor {
	return observe.ExporterDescriptor{
		Name:         "test.capture",
		Signals:      observe.SignalTraces,
		Stability:    observe.StabilityStable,
		Capabilities: observe.ExporterCapabilities{Push: true},
	}
}
func (*captureExporter) ForceFlush(context.Context) error { return nil }
func (*captureExporter) Shutdown(context.Context) error   { return nil }
func (e *captureExporter) ExportSpans(_ context.Context, spans []observe.SpanSnapshot) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.spans = append(e.spans, spans...)
	return nil
}
func (e *captureExporter) count() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.spans)
}

type trackingProvider struct {
	observe.Provider
	flushes   atomic.Int32
	shutdowns atomic.Int32
}

func (p *trackingProvider) ForceFlush(context.Context) error { p.flushes.Add(1); return nil }
func (p *trackingProvider) Shutdown(context.Context) error   { p.shutdowns.Add(1); return nil }
