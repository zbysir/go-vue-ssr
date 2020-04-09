package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type Context struct {
	Context context.Context
	Cancel  func()
}

// 等待信号量
func NewContext(sig ...os.Signal) (ctx context.Context, cancel func()) {
	ctx, cancel = context.WithCancel(context.TODO())

	if len(sig) == 0 {
		sig = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, sig...)
		<-ch

		cancel()
	}()
	return
}

// 等待Term和Int信号量关闭
func NewTermContext() (ctx context.Context, cancel func()) {
	return NewContext(syscall.SIGINT, syscall.SIGTERM)
}
