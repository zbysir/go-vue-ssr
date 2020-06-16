package vuessr

import (
	"fmt"
	"github.com/zbysir/go-vue-ssr/internal/pkg/log"
	"github.com/zbysir/go-vue-ssr/pkg/vuessr/ast"
	"strings"
)

func genPropsClassCode(classJs string) string {
	if classJs == "" {
		return "nil"
	}

	code, err := ast.Js2Go(classJs, ScopeKey)
	if err != nil {
		panic(err)
	}

	return code
}

func genProps(props Props) string {
	if len(props) == 0 {
		return "Props{}"
	}

	// orderKeyCode
	orderKeyCode := `[]string{`
	for _, p := range props {
		orderKeyCode += fmt.Sprintf(`"%s",`, p.Key)
	}
	orderKeyCode += "}"

	// dataCode
	dataCode := "map[string]interface{}{"
	for _, p := range props {
		k := p.Key
		v := p.Val
		valueCode, err := ast.Js2Go(v, ScopeKey)
		if err != nil {
			log.Panicf("%v, %s", err, v)
		}
		dataCode += fmt.Sprintf(`"%s": %s,`, k, valueCode)
	}
	dataCode += "}"

	return fmt.Sprintf(`Props{orderKey: %s, data: %s}`, orderKeyCode, dataCode)
}

func genPropsStyleCode(styleJs string) string {
	if styleJs == "" {
		return "nil"
	}

	code, err := ast.Js2Go(styleJs, ScopeKey)
	if err != nil {
		panic(err)
	}

	return code
}

// 生成!动态节点的!attr, 包括class style和其他
func genAllAttrCode(e *VueElement) string {
	var a = ""

	// go代码
	var classCode = ""
	var styleCode = ""
	var attrCode = ""

	// 查找props中的class 与 style, 将处理为动态class
	classProps, _ := e.Props.Get("class")
	styleProps, _ := e.Props.Get("style")

	// 额外处理class/style

	// class
	{
		// 静态Class GoCode
		staticClassCode := sliceStringToGoCode(e.Class)

		// 动态class GoCode
		classPropsCode := "nil"
		if classProps != "" {
			var err error
			classPropsCode, err = ast.Js2Go(classProps, ScopeKey)
			if err != nil {
				panic(err)
			}
		}

		if classPropsCode != "nil" {
			classCode = fmt.Sprintf(`mixinClass(nil, %s, %s)`, staticClassCode, classPropsCode)
		} else if staticClassCode == "nil" {
			classCode = ``
		} else {
			classCode = safeStringCode(fmt.Sprintf(` class="%s"`, strings.Join(e.Class, " ")))
		}
	}
	// style
	{
		staticStyleCode := mapStringToGoCode(e.Style)

		stylePropsCode := "nil"
		if styleProps != "" {
			var err error
			stylePropsCode, err = ast.Js2Go(styleProps, ScopeKey)
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
			styleCode = safeStringCode(fmt.Sprintf(` style="%s"`, genStyle(e.Style, e.StyleKeys)))
		}
	}
	// attr
	{
		// 静态attr GoCode
		staticAttrCode := genAttrsCode(e.Attrs)
		// 动态attr GoCode
		attrProps := e.Props.Omit("class", "style")

		// todo 可以预先判断static与Props是否有key冲突, 如果key不冲突, 则可以直接把static生成为go代码
		if len(attrProps) != 0 {
			attrPropsCode := genProps(attrProps)
			attrCode = fmt.Sprintf(`mixinAttr(nil, %s, %s)`, staticAttrCode, attrPropsCode)
		} else if staticAttrCode == "nil" {
			attrCode = ``
		} else {
			// 静态attrs 字符串
			attrCode = safeStringCode(fmt.Sprintf(` %s`, genAttr(e.Attrs)))
		}
	}

	if classCode != `` {
		a = classCode
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

func genAttrsCode(a []Attribute) string {
	if len(a) == 0 {
		return "nil"
	}
	st := "[]Attribute{\n"
	for _, v := range a {
		st += fmt.Sprintf(`{Key: %s, Val: %s},`, safeStringCode(v.Key), safeStringCode(v.Val))
	}
	st += "\n}"
	return st
}

// 生成静态style
func genStyle(style map[string]string, styleKeys []string) string {
	st := ""
	// 为了每次编译的代码都一样, style的顺序也应一样
	for _, k := range styleKeys {
		v := style[k]
		st += fmt.Sprintf("%s: %s; ", k, v)
	}
	return st
}

// 生成静态attr
func genAttr(attr []Attribute) string {
	c := strings.Builder{}
	for _, a := range attr {
		v := a.Val
		k := a.Key

		if c.Len() != 0 {
			c.WriteString(" ")
		}
		if v != "" {
			c.WriteString(fmt.Sprintf(`%s="%s"`, k, v))
		} else {
			c.WriteString(fmt.Sprintf(`%s`, k))
		}
	}

	return c.String()
}
