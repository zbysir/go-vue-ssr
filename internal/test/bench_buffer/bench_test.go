// cd internal/test/bench_string
// go-vue-ssr -src=./ -to=./ -pkg=bench_string

package bench_buffer

import (
	"encoding/json"
	"fmt"
	"testing"
)

type data struct {
	C   []*data `json:"c"`
	Msg string  `json:"msg"`
}

// 10000 103900070 ns/op
// 1000 9,174,621 ns/op
func BenchmarkBuffer(b *testing.B) {
	var ii interface{}
	// 生成100000个数据
	index := 0
	var ds []*data
	for i := 0; i < 1000; i++ {
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
		_ = r.Component_bench(&Options{
			Props: map[string]interface{}{
				"data": ii,
			},
		})
	}
}

func TestBuffer(b *testing.T) {
	var ii interface{}
	// 生成100000个数据
	index := 0
	var ds []*data
	for i := 0; i < 1000000; i++ {
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

	_ = r.Component_bench(&Options{
		Props: map[string]interface{}{
			"data": ii,
		},
	})

}

// 7946 ns/op
func BenchmarkBuffer2(b *testing.B) {
	var ii interface{}
	// 生成10000个嵌套数据
	index := 0
	d := &data{
		C:   nil,
		Msg: "1",
	}

	head := d
	for i := 0; i < 10000; i++ {
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
		_ = r.Component_bench(&Options{
			Props: map[string]interface{}{
				"data": ii,
			},
		})
	}
}
