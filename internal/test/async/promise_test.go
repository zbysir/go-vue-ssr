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

	p3 := &PromiseGroup{
		Note: "p3",
	}
	p3.AppendString("5")
	//p3.AppendString("5-1")
	t.Log(p3.Last == p3)

	//p.Append(p3)

	p2 := &PromiseGroup{
		Note: "p2",
	}
	p2.AppendGroup(p3)

	// 由于p2的长度为1, 所以last一定指向本身, 所以p2.Last = p2
	t.Log(p2.Last == p2, p2.Last == p3)

	p.AppendGroup(p2)

	p.Append("6")

	for cur := p; cur != nil; cur = cur.Next {
		t.Log(cur.Note)

	}

	t.Log(p.Join())
}
