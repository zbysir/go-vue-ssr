// cd internal/test/bench_string
// go-vue-ssr -src=./ -to=./ -pkg=bench_string

package bench_string

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"testing"
)

type data struct {
	C   []*data `json:"c"`
	Msg string  `json:"msg"`
}

func TestString(t *testing.T) {
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

	h := r.Component_bench(&Options{
		Props: map[string]interface{}{
			"data": ii,
		},
	})

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	kb := 1024.0
	//Alloc = 1849.40625	TotalAlloc = 4711.25	Sys = 10822.4921875	 NumGC = 1
	logstr := fmt.Sprintf("\nAlloc = %v\tTotalAlloc = %v\tSys = %v\t NumGC = %v\n\n", float64(m.Alloc)/kb, float64(m.TotalAlloc)/kb, float64(m.Sys)/kb, m.NumGC)
	t.Log(logstr)

	t.Log(h)
}

// 10000	110,299,760 ns/op
// 1000		10,105,247 ns/op
// 100		870,556 ns/op
func BenchmarkString(b *testing.B) {
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
		_ = r.Component_bench(&Options{
			Props: map[string]interface{}{
				"data": ii,
			},
		})
	}
}


func TestStringNest(t *testing.T) {
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

	h := r.Component_bench(&Options{
		Props: map[string]interface{}{
			"data": ii,
		},
	})

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	kb := 1024.0
	//  Alloc = 8654.1171875	TotalAlloc = 235710.8984375	Sys = 24082.8671875	 NumGC = 51
	logstr := fmt.Sprintf("\nAlloc = %v\tTotalAlloc = %v\tSys = %v\t NumGC = %v\n\n", float64(m.Alloc)/kb, float64(m.TotalAlloc)/kb, float64(m.Sys)/kb, m.NumGC)
	t.Log(logstr)

	t.Log(h)
}


// 测试递归嵌套的数据
// 测试主机: 公司
// 10000    15,579,999,800 ns/op
// 1000		151,428,086 ns/op
// 100      2,268,346 ns/op
func BenchmarkStringNest(b *testing.B) {
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

// 14375 ns/op
func BenchmarkAppendBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buffer bytes.Buffer
		for i := 0; i < 1000; i++ {
			buffer.WriteString("a")
		}
	}
}

// 301735 ns/op
func BenchmarkAppendString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buffer string
		for i := 0; i < 1000; i++ {
			buffer += "a"
		}
	}
}
