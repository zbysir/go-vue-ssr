package vuessr

import (
	"fmt"
	"strings"
)

// 处理静态的attr, 如class/style
func attr(attrs map[string]string, class []string) string {
	var a = ""
	if len(class) != 0 {
		a += fmt.Sprintf(`class=\"%s\"`, strings.Join(class, " "))
	}
	for k, v := range attrs {
		a += fmt.Sprintf(` %s=%v`, k, v)
	}

	return a
}
