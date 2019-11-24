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

func VForDemo() {
	res:="div"

	for index, item:=range lookInterfaceToStr(data,"list"){
		data = map[string]interface{}{
			"item": item,
			"index": index,
		}

		-> 子级
	}
}
