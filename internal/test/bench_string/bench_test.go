// cd internal/test/bench_string
// go-vue-ssr -src=./ -to=./ -pkg=bench_string

package bench_string

import (
	"fmt"
	"go.zhuzi.me/go/util"
	"testing"
)

type data struct {
	C   []*data `json:"c"`
	Msg string  `json:"msg"`
}

// 7415 ns/op
func BenchmarkString(b *testing.B) {
	var i interface{}
	// 生成10000个数据
	index := 0
	var ds []*data
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
	util.CopyObj(d, &i)

	r := NewRender()

	for i := 0; i < b.N; i++ {
		_ = r.Component_bench(&Options{
			Props: map[string]interface{}{
				"data": i,
			},
		})
	}
}

// 7946 ns/op
func BenchmarkString2(b *testing.B) {
	var i interface{}
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

	util.CopyObj(d, &i)

	r := NewRender()

	for i := 0; i < b.N; i++ {
		_ = r.Component_bench(&Options{
			Props: map[string]interface{}{
				"data": i,
			},
		})
	}
}

