# Goark Boot Observe

[简体中文](README.zh-CN.md)

`goark.dev/gbc-observe` integrates `goark.dev/observe` with Goark Boot. It creates or accepts an observability provider, registers stable provider and lifecycle beans, and flushes telemetry before shutdown.

The module is an assembly layer. It does not implement protocol instrumentation or concrete telemetry exporters.

Licensed under Apache License 2.0.
