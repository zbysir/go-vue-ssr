// cd internal/test/bench_buffer
// go-vue-ssr -src=./ -to=./ -pkg=bench_buffer

package bench_buffer

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"
)

type data struct {
	C   []*data `json:"c"`
	Msg string  `json:"msg"`
}

func TestBuffer(t *testing.T) {
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

	w := r.NewWriter()
	r.Component_bench(w, &Options{
		Props: map[string]interface{}{
			"data": ii,
		},
	})
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	kb := 1024.0
	//Alloc = 1048.96875	TotalAlloc = 3898.6484375	Sys = 9160.25	 NumGC = 1
	logstr := fmt.Sprintf("\nAlloc = %v\tTotalAlloc = %v\tSys = %v\t NumGC = %v\n\n", float64(m.Alloc)/kb, float64(m.TotalAlloc)/kb, float64(m.Sys)/kb, m.NumGC)
	t.Log(logstr)
	t.Log(w.Result())
}

// 10000 97,250,017 ns/op
// 1000 9,474,573 ns/op
// 100 839,999 ns/op
func BenchmarkBuffer(b *testing.B) {
	var ii interface{}
	// 生成100000个数据
	index := 0
	var ds []*data
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

	for i := 0; i < b.N; i++ {
		w := r.NewWriter()
		r.Component_bench(w, &Options{
			Props: map[string]interface{}{
				"data": ii,
			},
		})
		w.Result()
	}
}

func TestBufferNest(t *testing.T) {
	var ii interface{}
	// 生成100000个数据
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

	w := r.NewWriter()
	r.Component_bench(w, &Options{
		Props: map[string]interface{}{
			"data": ii,
		},
	})

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	kb := 1024.0
	// Alloc = 2138.0078125	TotalAlloc = 3934.1640625	Sys = 10822.4921875	 NumGC = 1
	logstr := fmt.Sprintf("\nAlloc = %v\tTotalAlloc = %v\tSys = %v\t NumGC = %v\n\n", float64(m.Alloc)/kb, float64(m.TotalAlloc)/kb, float64(m.Sys)/kb, m.NumGC)
	t.Log(logstr)

	t.Log(w.Result())
}

// 测试递归嵌套的数据
// 测试主机: 公司
// 10000 152,000,229 ns/op
// 1000  13,238,065 ns/op
// 100   905,501 ns/op
func BenchmarkBufferNest(b *testing.B) {
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
		w := r.NewWriter()
		r.Component_bench(w, &Options{
			Props: map[string]interface{}{
				"data": ii,
			},
		})
		_ = w.Result()
	}
}
