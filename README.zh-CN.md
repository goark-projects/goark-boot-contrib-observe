# Goark Boot Observe

[English](README.md)

`goark.dev/gbc-observe` 负责把 `goark.dev/observe` 接入 Goark Boot。它创建或接收 Provider，注册稳定的 Provider Bean 和生命周期 Bean，并在应用关闭前刷新观测数据。

本模块只负责装配，不实现协议埋点，也不实现具体 telemetry exporter。

本项目采用 Apache License 2.0。
