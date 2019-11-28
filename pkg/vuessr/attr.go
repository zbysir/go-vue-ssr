package vuessr

import (
	"fmt"
	"github.com/bysir-zl/vue-ssr/pkg/vuessr/ast_from_api"
	"strings"
)

// 生成attr, 包括class style和其他
// isRoot: 当时root节点的时候才会从options里读取上一层传递而来的数据 用来组装
func genAttr(e *VueElement) string {
	var a = ""

	// go代码
	var classCode = ""
	var styleCode = ""

	// 查找props中的class 与 style, 将处理为动态class
	classProps := ""
	styleProps := ""
	if e.Props != nil {
		classProps = e.Props["class"]
		styleProps = e.Props["style"]
	}

	// 如果是root组件, 则始终应该使用动态class/style(因为上级传递的class/style不是固定的), 应该合并上级和本级的class/style
	if e.IsRoot {
		// class
		{
			staticClassCode := sliceStringToGoCode(e.Class)

			classPropsCode := "nil"
			if classProps != "" {
				var err error
				classPropsCode, err = ast_from_api.JsCode2Go(classProps, DataKey)
				if err != nil {
					panic(err)
				}
			}

			classCode = fmt.Sprintf(`"class=\""+mixinClass(options, %s, %s)+"\""`, staticClassCode, classPropsCode)
		}

		// style
		{
			staticStyleCode := mapStringToGoCode(e.Style)

			stylePropsCode := "nil"
			if styleProps != "" {
				var err error
				stylePropsCode, err = ast_from_api.JsCode2Go(styleProps, DataKey)
				if err != nil {
					panic(err)
				}
			}

			styleCode = fmt.Sprintf(`"style=\""+mixinStyle(options, %s, %s)+"\""`, staticStyleCode, stylePropsCode)
		}
	} else {
		// class
		{
			staticClassCode := sliceStringToGoCode(e.Class)

			classPropsCode := "nil"
			if classProps != "" {
				var err error
				classPropsCode, err = ast_from_api.JsCode2Go(classProps, DataKey)
				if err != nil {
					panic(err)
				}
			}
			if classPropsCode != "nil" {
				classCode = fmt.Sprintf(`"class=\""+mixinClass(nil, %s, %s)+"\""`, staticClassCode, classPropsCode)
			} else if staticClassCode == "nil" {
				classCode = `""`
			} else {
				classCode = fmt.Sprintf(`"class=\"%s\""`, strings.Join(e.Class, " "))
			}
		}

		// style
		{
			staticStyleCode := mapStringToGoCode(e.Style)

			stylePropsCode := "nil"
			if styleProps != "" {
				var err error
				stylePropsCode, err = ast_from_api.JsCode2Go(styleProps, DataKey)
				if err != nil {
					panic(err)
				}
			}
			if stylePropsCode != "nil" {
				styleCode = fmt.Sprintf(`"style=\""+mixinStyle(nil, %s, %s)+"\""`, staticStyleCode, stylePropsCode)
			} else if staticStyleCode == "nil" {
				styleCode = `""`
			} else {
				styleCode = fmt.Sprintf(`"style=\"%s\""`, genStyle(e.Style, e.StyleKeys))
			}
		}
	}

	a += classCode

	// props类
	// 需要传递给子级的变量, 全部需要显示写v-bind:, 不支持像vue一样的传递字符串变量可以这样写 <el-input v-bind:size="mini" >

	// 样式
	if len(styleCode) != 0 {
		if a != "" {
			a += `+" "+`
		}
		a += styleCode
	}

	// 其他属性
	if len(e.Attrs) != 0 {
		if a != "" {
			a += `+" "+`
		}

		at := ""
		for k, v := range e.Attrs {
			at += fmt.Sprintf(`%s=%v `, k, v)
		}

		a += `"` + at + `"`
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
