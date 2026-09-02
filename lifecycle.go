package gbcobserve

import (
	"context"
	"errors"

	"goark.dev/observe"
)

type providerLifecycle struct{ provider observe.Provider }

// Stop 先刷新已接收数据，再关闭 Provider。
func (l *providerLifecycle) Stop(ctx context.Context) error {
	if l == nil || l.provider == nil {
		return nil
	}
	return errors.Join(l.provider.ForceFlush(ctx), l.provider.Shutdown(ctx))
}

// Order 让 Provider 在其他普通组件之后停止。
func (*providerLifecycle) Order() int { return -10000 }
