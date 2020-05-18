// cd internal/test/async
// go-vue-ssr -src=./ -to=./ -pkg=async

package async

import (
	"encoding/json"
	"fmt"
	"github.com/zbysir/go-vue-ssr/internal/pkg/log"
	"testing"
)

type data struct {
	C   []*data `json:"c"`
	Msg string  `json:"msg"`
}

// 家
// 10000	74,666,960 ns/op
// 1000		7,141,393 ns/op
// 10		53,782 ns/op

// 公司
// 10000    386,333,300 ns/op v1
//          157,857,743 ns/op // v2异步
//          134,125,250 ns/op // v2不异步
// 1000     19,700,057 ns/op v1
//          15,750,026 ns/op v2
// 10		53,782 ns/op
func BenchmarkString(b *testing.B) {
	var ii interface{}
	index := 0
	var ds []*data
	// 生成1000个数据
	for i := 0; i < 10000; i++ {
		ds = append(ds, &data{
			C:   nil,
			Msg: fmt.Sprintf("%d", index),
		})
		index++
	}

	d := data{
		C:   ds,
		Msg: "1",
	}
	bs, _ := json.Marshal(d)
	json.Unmarshal(bs, &ii)

	r := NewRender()

	for i := 0; i < b.N; i++ {
		g := r.Component_bench(&Options{
			Props: map[string]interface{}{
				"data": ii,
			},
		})
		g.Join()
		b.Log(g.Len())
	}
}

// 公司
// 1000		71,133,400 ns/op
// 100 		3,182,321 ns/op
func BenchmarkString2(b *testing.B) {
	var ii interface{}
	// 生成10000个嵌套数据
	index := 0
	d := &data{
		C:   nil,
		Msg: "1",
	}

	head := d
	for i := 0; i < 1000; i++ {
		n := &data{
			C:   nil,
			Msg: fmt.Sprintf("%d", index),
		}
		head.C = append(head.C, n)
		head = n

		index++
	}
	bs, _ := json.Marshal(d)
	json.Unmarshal(bs, &ii)

	r := NewRender()

	for i := 0; i < b.N; i++ {
		g := r.Component_bench(&Options{
			Props: map[string]interface{}{
				"data": ii,
			},
		})
		g.Join()
		b.Log(g.Len())
	}
}

func TestAsync(t *testing.T) {
	var ii interface{}
	index := 0

	var ds []*data
	for i := 0; i < 10; i++ {
		ds = append(ds, &data{
			C:   nil,
			Msg: fmt.Sprintf("%d", index),
		})
		index++
	}

	d := data{
		C:   ds,
		Msg: "1",
	}
	bs, _ := json.Marshal(d)
	json.Unmarshal(bs, &ii)

	r := NewRender()
	g := r.Component_bench(&Options{
		Props: map[string]interface{}{
			"data": ii,
		},
	})

	log.Infof("%+v", g.Len())
	t.Log(g.Join())
}
