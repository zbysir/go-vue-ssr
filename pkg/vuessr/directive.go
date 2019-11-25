package vuessr

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type Directive interface {
	Exec(e *VueElement, code string) string
}

type Directives map[string]Directive

func getDirectives(attrs []xml.Attr) (ds Directives) {
	ds = Directives{}
	for _, v := range attrs {
		if strings.HasPrefix(v.Name.Local, "v-") {
			name := v.Name.Local
			switch name {
			case "v-for":
				ds[name] = getVForDirective(v.Value)
			case "v-if":
				ds[name] = getVIfDirective(v.Value)
			}
		}
	}
	return
}

func (d Directives) Exec(e *VueElement, code string) string {
	for _, v := range d {
		code = v.Exec(e, code)
	}
	return code
}

type VForDirective struct {
	arrayKey string
	itemKey  string
	indexKey string
}

func (e VForDirective) Exec(el *VueElement, code string) string {
	vfArray := e.arrayKey
	vfItem := e.itemKey
	vfIndex := e.indexKey
	// 将自己for
	return fmt.Sprintf(`
func ()string{
  var c = ""

  for index, item := range lookInterfaceToSlice(data, "%s") {
    c += func(data map[string]interface{}) string{
        data = extendMap(map[string]interface{}{
          "%s": index,
          "%s": item,
        }, data)

        return %s
    }(data)
  }
return c
}()`, vfArray, vfIndex, vfItem, code)
}

type VIfDirective struct {
	condition string
}

func (e VIfDirective) Exec(el *VueElement, code string) string {
	// 将自己for
	return fmt.Sprintf(`
func ()string{
  if condition(data, "%s") {return %s}
  return ""
}()`, e.condition, code)
}

func getVForDirective(raw string) (d VForDirective) {
	val := raw

	ss := strings.Split(val, " in ")
	d.arrayKey = strings.Trim(ss[1], " ")

	left := strings.Trim(ss[0], " ")
	// (item, index) in list
	if strings.Contains(left, ",") {
		left = strings.Trim(left, "()")
		ss := strings.Split(left, ",")
		d.itemKey = strings.Trim(ss[0], " ")
		d.indexKey = strings.Trim(ss[1], " ")

	} else {
		d.itemKey = left
		d.indexKey = "$index"
	}

	return
}

func getVIfDirective(raw string) (d VIfDirective) {
	d.condition = strings.Trim(raw, " ")
	return
}
