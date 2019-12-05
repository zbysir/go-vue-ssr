package vuessr

import (
	"fmt"
	"github.com/bysir-zl/go-vue-ssr/pkg/vuessr/ast_from_api"
	"strings"
)

func getPropsClass(props Props) string {
	item := props.Get("class")
	if item == "" {
		return "nil"
	}

	code, err := ast_from_api.Js2Go(item, DataKey)
	if err != nil {
		panic(err)
	}

	return code
}

func getPropsStyle(props Props) string {
	item := props.Get("style")
	if item == "" {
		return "nil"
	}

	code, err := ast_from_api.Js2Go(item, DataKey)
	if err != nil {
		panic(err)
	}

	return code
}

// 生成!动态节点的!attr, 包括class style和其他
func genAttrCode(e *VueElement) string {
	var a = ""

	// go代码
	var classCode = ""
	var styleCode = ""
	var attrCode = ""

	// 查找props中的class 与 style, 将处理为动态class
	classProps := e.Props.Get("class")
	styleProps := e.Props.Get("style")

	// 额外处理class/style

	// class
	{
		staticClassCode := sliceStringToGoCode(e.Class)

		classPropsCode := "nil"
		if classProps != "" {
			var err error
			classPropsCode, err = ast_from_api.Js2Go(classProps, DataKey)
			if err != nil {
				panic(err)
			}
		}
		if classPropsCode != "nil" {
			classCode = fmt.Sprintf(`mixinClass(nil, %s, %s)`, staticClassCode, classPropsCode)
		} else if staticClassCode == "nil" {
			classCode = ``
		} else {
			classCode = fmt.Sprintf(`" class=\"%s\""`, strings.Join(e.Class, " "))
		}
	}
	// style
	{
		staticStyleCode := mapStringToGoCode(e.Style)

		stylePropsCode := "nil"
		if styleProps != "" {
			var err error
			stylePropsCode, err = ast_from_api.Js2Go(styleProps, DataKey)
			if err != nil {
				panic(err)
			}
		}
		if stylePropsCode != "nil" {
			// todo 可以预先判断static与Props是否有key冲突, 如果key不冲突, 则可以直接把static生成为go代码
			styleCode = fmt.Sprintf(`mixinStyle(nil, %s, %s)`, staticStyleCode, stylePropsCode)
		} else if staticStyleCode == "nil" {
			styleCode = ``
		} else {
			styleCode = fmt.Sprintf(`" style=\"%s\""`, genStyle(e.Style, e.StyleKeys))
		}
	}
	// attr
	{
		staticAttrCode := mapStringToGoCode(e.Attrs)
		attrPropsCode := mapJsCodeToCode(e.Props.Omit("class", "style"))

		// todo 可以预先判断static与Props是否有key冲突, 如果key不冲突, 则可以直接把static生成为go代码
		if attrPropsCode != "nil" {
			attrCode = fmt.Sprintf(`mixinAttr(nil, %s, %s)`, staticAttrCode, attrPropsCode)
		} else if staticAttrCode == "nil" {
			attrCode = ``
		} else {
			attrCode = fmt.Sprintf(`" %s"`, genAttr(e.Attrs, e.AttrsKeys))
		}
	}

	if classCode != `` {
		a += classCode
	}

	// 样式
	if styleCode != `` {
		if a != "" {
			a += `+`
		}
		a += styleCode
	}

	// attr
	if attrCode != `` {
		if a != "" {
			a += `+`
		}
		a += attrCode
	}

	if a == "" {
		a = `""`
	}

	return a
}

func genStyle(style map[string]string, styleKeys []string) string {
	st := ""
	// 为了每次编译的代码都一样, style的顺序也应一样
	for _, k := range styleKeys {
		v := style[k]
		st += fmt.Sprintf("%s: %s; ", k, v)
	}
	return st
}

func genAttr(attr map[string]string, keys []string) string {
	c := ""
	// 为了每次编译的代码都一样, style的顺序也应一样
	for _, k := range keys {
		v := attr[k]
		c += fmt.Sprintf(`%s=\"%s\"`, k, v)
	}
	return c

}
