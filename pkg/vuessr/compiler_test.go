package vuessr

import (
	"strconv"
	"strings"
	"testing"
)

func TestParseVueVif(t *testing.T) {
	e, err := ParseVue(`Z:\go_project\go-vue-ssr\internal\test\vue\page.vue`)
	if err != nil {
		t.Fatal(err)
	}
	c := NewCompiler()
	code, _ := c.GenEleCode(e)

	code = minifyCode(code)

	t.Log(code)
	return
}

func TestQuote(t *testing.T) {
	want := `"\"\"{{title +""}}"`
	x := safeStringCode(`""{{title +""}}`)
	if x != want {
		t.Fatalf("%v; want:%v", x, want)
	}
}

func TestInjectVal(t *testing.T) {
	want := `interfaceToStr(scope.Get("total"), true)`
	x := injectVal(`{{total}}`)
	if x != want {
		t.Fatalf("%s; want: %s", x, want)
	}
}

func TestTextNode(t *testing.T) {
	text := `123 {{title}}`
	code := safeStringCode(text)

	t.Log(code)
	// 处理变量
	code = injectVal(code)

	want := `interfaceToStr(scope.Get("title"), true)`
	if code != want {
		t.Fatalf("code = %v; want:%v", code, want)
	}
}

// 使用结果表明, 能用+号的地方就用+号
func BenchmarkAppend(b *testing.B) {
	// 171 ns/op
	b.Run("builder", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var x strings.Builder
			x.WriteString("123123123123" + "123123123123" + strconv.Itoa(i))
			x.String()
		}
	})

	// 191 ns/op
	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var x strings.Builder
			x.WriteString("123123123123")
			x.WriteString("123123123123")
			x.WriteString(strconv.Itoa(i))
			x.String()
		}
	})
}

func TestMini(t *testing.T) {
	src := `w.WriteString("<head>")
w.WriteString("<link rel=\"stylesheet\"href=\"//static.f.cdn-static.cn/3.7.0/animate.min.css\"type=\"text/css\"></link>")
w.WriteString("<link rel=\"stylesheet\"href=\"//static.f.cdn-static.cn/swiper/swiper.min.css\"></link>")`

	t.Log(minifyCode(src))
}
