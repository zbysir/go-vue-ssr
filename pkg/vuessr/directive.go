package vuessr

import (
	"fmt"
	"github.com/bysir-zl/go-vue-ssr/pkg/vuessr/ast_from_api"
	"golang.org/x/net/html"
	"strings"
)

// 代码生成指令, 处理所有需要编译期间执行的指令. 暂不支持运行时自定义指令.
// 不过可以自定义编译时指令.

// 编译时指令分为两种, 一种是生成代码之后执行的指令, 一种是生成代码之前执行的指令
// - 生成代码之前的指令中, 可以操作所有VueElement中的属性, 来达到提前修改数据的目的.
// - 生成代码之后的指令中, 可以处理生成的代码, 可以使用正则修改代码, 更加灵活.

// 指令:
// 指令会影响当前节点的渲染, 返回修改后的go代码
// 有一个特殊的指令: v-slot, 会将节点代码改为空, 并且写入到namedSlotCode里.
type GenCodeDirective interface {
	Exec(e *VueElement, app *App, code string) (resCode string, namedSlotCode map[string]string)
}

type GenCodeDirectives map[string]GenCodeDirective

func (d GenCodeDirectives) Exec(e *VueElement, app *App, code string) (descCode string, namedSlotCode map[string]string) {
	namedSlotCode = map[string]string{}
	for _, v := range d {
		var n2 map[string]string
		code, n2 = v.Exec(e, app, code)
		for k, v := range n2 {
			namedSlotCode[k] = v
		}
	}
	return code, namedSlotCode
}

type VForDirective struct {
	arrayKey string
	itemKey  string
	indexKey string
}

func (e VForDirective) Exec(el *VueElement, app *App, code string) (descCode string, namedSlotCode map[string]string) {
	vfArray := e.arrayKey
	vfItem := e.itemKey
	vfIndex := e.indexKey
	// 将自己for, 将子代码的data字段覆盖, 实现作用域的修改
	return fmt.Sprintf(`func ()string{
  var c = ""

  for index, item := range lookInterfaceToSlice(%s, "%s") {
    c += func(xdata map[string]interface{}) string{
        %s := extendMap(map[string]interface{}{
          "%s": index,
          "%s": item,
        }, xdata)

        return %s
    }(%s)
  }
return c
}()`, DataKey, vfArray, DataKey, vfIndex, vfItem, code, DataKey), nil
}

type VIfDirective struct {
	Condition string
	ElseIf    []VIfDirectiveElseIf
}

type VIfDirectiveElseIf struct {
	Types     string // elseif or else
	Condition string // 条件表达式
	Code      string // 子节点代码
}

func (e VIfDirective) Exec(el *VueElement, app *App, code string) (descCode string, namedSlotCode map[string]string) {
	condition, err := ast_from_api.JsCode2Go(e.Condition, DataKey)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(`func ()string{
  if interfaceToBool(%s) {return %s}
  return ""
}()`, condition, code), nil
}

func (e VIfDirective) AddElseIf(el *VueElement) () {

}

type VSlotDirective struct {
	slotName string
	propsKey string
}

func (e VSlotDirective) Exec(el *VueElement, app *App, code string) (descCode string, namedSlotCode map[string]string) {
	// 插槽支持传递props, 需要有自己的作用域, 所以需要使用闭包实现
	code = fmt.Sprintf(`func(props map[string]interface{}) string{
	%s := extendMap(map[string]interface{}{"%s": props}, %s)
	return %s
}`, DataKey, e.propsKey, DataKey, code)

	namedSlotCode = map[string]string{
		e.slotName: code,
	}

	// 插槽会将原来的子代码去掉, 并将代码放在namedSlot里.
	descCode = `""`
	return
}

// raw: 指令的值
func getVForDirective(attr html.Attribute) (d VForDirective) {
	val := attr.Val

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

func getVIfDirective(attr html.Attribute) (d VIfDirective) {
	d.Condition = strings.Trim(attr.Val, " ")
	return
}

func getVSlotDirective(attr html.Attribute) (d VSlotDirective) {
	oriKey := attr.Key
	key := oriKey
	ss := strings.Split(oriKey, ":")
	if len(ss) == 2 {
		key = ss[1]
	}
	d.slotName = key
	d.propsKey = attr.Val
	// 不应该为空, 否则可能会导致生成的go代码有误
	if d.propsKey == "" {
		d.propsKey = "slotProps"
	}

	return
}
