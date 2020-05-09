package async

import (
	"testing"
	"time"
)

func x() *PromiseGroup {
	p := &PromiseGroup{}

	p.Append("100")
	return p
}

func TestPromise(t *testing.T) {
	var p = &PromiseGroup{
	}

	p.Append("1")
	p.Append("2")
	p.Append(PromiseString("3"))
	{

		s := make(PromiseChan, 1)
		go func() {
			time.Sleep(2 * time.Second)
			s <- "4"
		}()
		p.Append(s)
	}

	p3 := &PromiseGroup{}
	p3.Append("xxxx")

	p.Append(p3)

	p2 := &PromiseGroup{}
	p2.Append("6")

	p.AppendGroup(p2)

	t.Log(p.Join())
}
