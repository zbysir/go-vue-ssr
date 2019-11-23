package genera

import "testing"

func TestXComponent_main(t *testing.T) {
	html := XComponent_helloworld(map[string]interface{}{
		"name": "bysir",
		"sex":  "男",
		"age":  "18",
	}, "")
	t.Log(html)
}
