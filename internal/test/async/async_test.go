// cd internal/test/async
// go-vue-ssr -src=./ -to=./ -pkg=async

package async

import (
	"encoding/json"
	"fmt"
	"testing"
)

type data struct {
	C   []*data `json:"c"`
	Msg string  `json:"msg"`
}

// 10000	74,666,960 ns/op
// 1000		7,141,393 ns/op
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
	}
}

// 1000		35,269,300 ns/op
// 100 		1,623,805 ns/op
func BenchmarkString2(b *testing.B) {
	var ii interface{}
	// 生成10000个嵌套数据
	index := 0
	d := &data{
		C:   nil,
		Msg: "1",
	}

	head := d
	for i := 0; i < 100; i++ {
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
	}
}

func TestAsync(t *testing.T) {
	var ii interface{}
	index := 0

	var ds []*data
	// 生成100个数据
	for i := 0; i < 100; i++ {
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
	t.Log(g.Join())
}
