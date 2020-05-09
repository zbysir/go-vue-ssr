package async

import (
	"testing"
	"time"
)

func TestPromise(t *testing.T) {
	var p = &PromiseGroup{
	}

	p.Append(PromiseString("1"))
	p.Append(PromiseString("2"))
	p.Append(PromiseString("3"))
	{

		s := make(chan string)
		go func() {
			time.Sleep(2 * time.Second)
			s <- "4"
		}()
		p.Append(PromiseFunc(func() string {
			return <-s
		}))
	}
	{

		s := make(chan string)
		go func() {
			time.Sleep(3 * time.Second)
			s <- "5"
		}()
		p.Append(PromiseFunc(func() string {
			return <-s
		}))

	}

	p2 := &PromiseGroup{}
	p2.Append(PromiseString("6"))
	p.AppendGroup(p2)

	t.Log(p.Join())
}
