# Goark Boot Observe

[English](README.md)

`goark.dev/gbc-observe` 负责把 `goark.dev/observe` 接入 Goark Boot。它创建或接收 Provider，注册稳定的 Provider Bean 和生命周期 Bean，并在应用关闭前刷新观测数据。

本模块只负责装配，不实现协议埋点，也不实现具体 telemetry exporter。

## 使用

```go
gbcobserve.AutoConfigure(
    gbcobserve.WithExporters(exporter),
    gbcobserve.WithProcessors(processor),
)
```

starter 将 `goark.observe.provider` 注册为主 `observe.Provider` Bean，并将 `goark.observe.lifecycle` 注册为生命周期适配器。关闭时先刷新已接收数据，再调用 Provider 的幂等关闭契约。

各 exporter starter 以延迟 Bean 形式注册 `observe.Exporter`。Provider 创建时统一收集这些 Bean，因此本模块无需依赖任何具体 exporter；直接装配仍可使用 `WithExporters`。通过 `WithProvider` 传入外部 Provider 时，其 exporter 也由外部负责装配。

## 配置

- `goark.observe.enabled`：默认为 `true`；禁用时注册 no-op Provider。
- `goark.observe.service.name`：默认为 `goark`。
- `goark.observe.service.version`、`.namespace`、`.instance-id`、`.environment`：资源身份。
- `goark.observe.metrics.cardinality-limit`：默认为 `2000`，且必须为正数。

外部管理的 Provider 使用 `WithProvider`，自定义 SDK 使用 `WithProviderFactory`。

本项目采用 Apache License 2.0。
