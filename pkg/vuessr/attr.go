package vuessr

import (
	"fmt"
	"go.zhuzi.me/go/log"
	"strings"
)

// 处理静态的attr, 如class/style
func genAttr(e *VueElement) string {
	var a = ""
	// 类
	if len(e.Class) != 0 {
		a += fmt.Sprintf(`class="%s"`, strings.Join(e.Class, " "))
	}

	// props类
	// 需要传递给子级的变量, 全部需要显示写v-bind:, 不支持像vue一样的传递字符串变量可以这样写 <el-input v-bind:size="mini" >


	// 样式
	if len(e.Style) != 0 {
		log.Infof("%+v", e.Style)
		if a != "" {
			a += " "
		}
		st := ""
		// 为了每次编译的代码都一样, style的顺序也应一样
		for _, k := range e.StyleKeys {
			v := e.Style[k]
			st += fmt.Sprintf("%s: %s; ", k, v)
		}
		a += fmt.Sprintf(`style="%s"`, st)
	}

	// 其他属性
	if len(e.Attrs) != 0 {
		if a != "" {
			a += " "
		}

		for k, v := range e.Attrs {
			a += fmt.Sprintf(`%s=%v `, k, v)
		}
	}

	return a
}
