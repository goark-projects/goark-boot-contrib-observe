# Goark Boot Observe

[简体中文](README.zh-CN.md)

`goark.dev/gbc-observe` integrates `goark.dev/observe` with Goark Boot. It creates or accepts an observability provider, registers stable provider and lifecycle beans, and flushes telemetry before shutdown.

The module is an assembly layer. It does not implement protocol instrumentation or concrete telemetry exporters.

## Usage

```go
gbcobserve.AutoConfigure(
    gbcobserve.WithExporters(exporter),
    gbcobserve.WithProcessors(processor),
)
```

The starter registers `goark.observe.provider` as the primary `observe.Provider` bean and `goark.observe.lifecycle` as its lifecycle adapter. Shutdown force-flushes accepted telemetry and then invokes the provider's idempotent shutdown contract.

## Properties

- `goark.observe.enabled`: defaults to `true`; disabled mode registers the no-op provider.
- `goark.observe.service.name`: defaults to `goark`.
- `goark.observe.service.version`, `.namespace`, `.instance-id`, `.environment`: resource identity.
- `goark.observe.metrics.cardinality-limit`: defaults to `2000` and must be positive.

Use `WithProvider` for an externally owned provider or `WithProviderFactory` for a custom SDK.

Licensed under Apache License 2.0.
