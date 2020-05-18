// generate by ./generotor_builtin_source/main.go
package vuessr

const builtinCode = `

// src: ./generotor_builtin_source/source.go
import (
	"encoding/json"
	"fmt"
	"html"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Render struct {
	// 全局变量, 可以理解为js中的windows, 每个组件中都可以直接读取到这个对象中的值.
	// 其中可以Set签名为function的方法, 供{{func(a)}}语法使用.
	Global *Global
	// 注册的动态组件
	Components map[string]ComponentFunc
	// 指令
	directives    map[string]DirectivesFunc
	writerCreator func() Writer
}

func (r Render) NewWriter() Writer {
	return r.writerCreator()
}

func newRender(options ...RenderOption) *Render {
	r := &Render{
		Global:     &Global{NewScope()},
		Components: nil,
		directives: nil,
		writerCreator: func() Writer {
			return NewBufferSpans()
		},
	}

	for _, o := range options {
		o(r)
	}

	return r
}

type RenderOption func(r *Render)

func WithWriter(c func() Writer) RenderOption {
	return func(r *Render) {
		r.writerCreator = c
	}
}

type Global struct {
	*Scope
}

func (p *Global) Func(name string, f Function) {
	p.Scope.Set(name, f)
}

func (p *Global) Var(name string, v interface{}) {
	p.Scope.Set(name, v)
}

// for {{func(a)}}
type Function func(args ...interface{}) interface{}

type DirectivesBinding struct {
	Value interface{}
	Arg   string
	Name  string
}

type DirectivesFunc func(b DirectivesBinding, options *Options)

func emptyFunc(args ...interface{}) interface{} {
	if len(args) != 0 {
		return args[0]
	}
	return nil
}

// js中的作用域
type Scope struct {
	p      *Scope
	values map[string]interface{}
}

func (s *Scope) ParentScope() *Scope {
	return s.p
}

// 设置暂时只支持在当前作用域设置变量
func (s *Scope) Set(k string, v interface{}) {
	s.values[k] = v
}

// 查找作用域中的变量, 返回变量所在的map
func (s *Scope) Find(k string) map[string]interface{} {
	curr := s
	for curr != nil {
		if _, ok := curr.values[k]; ok {
			return curr.values
		}

		curr = curr.p
	}

	return nil
}

func NewScope() *Scope {
	return &Scope{
		p:      nil,
		values: map[string]interface{}{},
	}
}

// 获取作用域中的变量
// 会向上查找
func (s *Scope) Get(k ...string) (v interface{}) {
	var rootExist bool
	var ok bool

	curr := s
	for curr != nil {
		v, rootExist, ok = shouldLookInterface(curr.values, k...)
		// 如果root存在, 则说明就应该读取当前作用域, 否则向上层作用域查找
		if rootExist {
			if !ok {
				return nil
			} else {
				return
			}
		}

		curr = curr.p
	}

	return
}

type Writer interface {
	// 如果需要实现异步计算, 则需要将span存储, 在最后统一计算出string.
	WriteSpan(Span)
	// 如果是同步计算, 使用WriteString会将string结果直接存储或者拼接
	WriteString(string)
	Result() string
}

type Span interface {
	Result() string
}

// 将多个Promise拼接为一个, 以减少内存与链的长度
type BufferSpan struct {
	s *strings.Builder
}

func (p *BufferSpan) Result() string {
	return p.s.String()
}

func (p *BufferSpan) WriteString(s string) {
	p.s.WriteString(s)
}

func NewBufferSpan(s string) Span {
	var b strings.Builder
	b.WriteString(s)
	return &BufferSpan{
		s: &b,
	}
}

// buffer块, 同步计算
type BufferWriter struct {
	s *strings.Builder
}

func (p BufferWriter) WriteSpan(span Span) {
	p.s.WriteString(span.Result())
}

func (p BufferWriter) WriteString(s string) {
	p.s.WriteString(s)
}

func (p BufferWriter) Result() string {
	return p.s.String()
}

func NewBufferSpans() Writer {
	var b strings.Builder
	return &BufferWriter{
		s: &b,
	}
}

// ListSpans将存储Span链表, 在最后计算结果, 可以实现并行计算.
type ListSpans struct {
	Value Span
	Next  *ListSpans
	Last  *ListSpans // 用于在append时提升速度
}

func (p *ListSpans) WriteSpans(s Writer) {
	switch t := s.(type) {
	case *ListSpans:
		if t == nil || t.Value == nil {
			return
		}

		if p.Value == nil {
			if t.Next != nil {
				// 跳过s的第一个元素, 将值存储到自己
				// 注意: 如果s只有一个元素, 由于s.last存储的是s自己, p.Last也赋值为s.last的话, 如果跳过s, 就导致了p.Last存储了一个被抛弃(跳过)的元素, 当下次赋值p.Last.Next就会出错
				p.Value = t.Value
				p.Last = t.Last
				p.Next = t.Next
			} else {
				// 如果s只有一个元素, 则抛弃s, 由p自己存储此元素
				p.WriteSpan(t.Value)
			}
			return
		}

		if p.Last == nil || t.Last == nil {
			panic("last不能为空")
		}

		// TODO 如果Last和t第一个元素可以合并, 则再合并一次
		p.Last.Next = t
		p.Last = t.Last
	default:
		panic("listSpan support Append listSpan only")
	}
}

func (l *ListSpans) WriteString(s string) {
	l.WriteSpan(NewBufferSpan(s))
}

func (p *ListSpans) WriteSpan(s Span) {
	if p.Value == nil {
		p.Value = s
		p.Last = p
		return
	}

	// 如果s是StringSpan并且p.Last也是StringSpan的话, 就将s的值附加到Last上
	// 以减少链的长度
	if ss, ok := s.(*BufferSpan); ok {
		if ls, ok := p.Last.Value.(*BufferSpan); ok {
			ls.WriteString(ss.Result())
			return
		}
	}

	last := &ListSpans{
		Value: s,
	}

	p.Last.Next = last
	p.Last = last
}

func (l *ListSpans) Result() string {
	if l == nil || l.Value == nil {
		return ""
	}

	b := strings.Builder{}

	for cur := l; cur != nil; cur = cur.Next {
		b.WriteString(cur.Value.Result())
	}

	return b.String()
}

func (l *ListSpans) Length() int {
	if l == nil || l.Value == nil {
		return 0
	}

	i := 0
	for cur := l; cur != nil; cur = cur.Next {
		i++
	}

	return i
}

func NewListSpans() Writer {
	return &ListSpans{}
}

type ChanSpan struct {
	c       chan string
	getOnce sync.Once
	setOnce sync.Once
	r       string
}

func (p *ChanSpan) Result() string {
	p.getOnce.Do(func() {
		p.r = <-p.c
	})
	return p.r
}

func (p *ChanSpan) Done(s string) {
	p.setOnce.Do(func() {
		p.c <- s
	})
}

func NewChanSpan() *ChanSpan {
	return &ChanSpan{
		c: make(chan string, 1),
	}
}

// 注册指令
func (r *Render) Directive(name string, f DirectivesFunc) {
	if r.directives == nil {
		r.directives = map[string]DirectivesFunc{}
	}

	r.directives[name] = f
}

// 内置组件Slot, 将渲染父级传递的slot.
func (r *Render) Component_slot(w Writer, options *Options) {
	name := options.Attrs["name"]
	if name == "" {
		name = "default"
	}
	props := options.Props
	injectSlotFunc, ok := options.P.Slots[name]

	// 如果没有传递slot 则使用自身默认的slot
	if !ok {
		injectSlotFunc = options.Slots["default"]
	}

	injectSlotFunc.Exec(w, props)
}

func (r *Render) Component_async(w Writer, options *Options) {
	scope := extendScope(r.Global.Scope, options.Props)
	options.Directives.Exec(r, options)
	_ = scope

	s := NewChanSpan()
	// 异步子节点计算
	go func() {
		w := r.NewWriter()
		options.Slots.Exec(w, "default", nil)
		s.Done(w.Result())
	}()

	w.WriteSpan(s)

	return
}

func (r *Render) Component_component(w Writer, options *Options) {
	is, ok := options.Props["is"].(string)
	if !ok {
		return
	}
	if c, ok := r.Components[is]; ok {
		c(w, options)
		return
	}
	w.WriteString(fmt.Sprintf("<p>not register com: %s</p>", is))
}

func (r *Render) Component_template(w Writer, options *Options) {
	// exec directive
	options.Directives.Exec(r, options)

	options.Slots.Exec(w, "default", nil)
}

// 动态tag
// 何为动态tag:
// - 每个组件的root层tag(attr受到上层传递的props影响)
// - 有自己定义指令(自定义指令需要修改组件所有属性, 只能由动态tag实现)
func (r *Render) tag(w Writer, tagName string, isRoot bool, options *Options) {
	// exec directive
	options.Directives.Exec(r, options)

	var p *Options
	if isRoot {
		p = options.P
	}

	// attr
	attr := mixinClass(p, options.Class, options.PropsClass) +
		mixinStyle(p, options.Style, options.PropsStyle) +
		mixinAttr(p, options.Attrs, options.Props)

	w.WriteString(fmt.Sprintf("<%s%s>", tagName, attr))
	options.Slots.Exec(w, "default", nil)
	w.WriteString(fmt.Sprintf("</%s>", tagName))

	return
}

// 渲染组件需要的结构
// tips: 此结构应该尽量的简单, 方便渲染才能性能更好.
type Options struct {
	Props      Props                  // 本节点的数据(不包含class和style)
	PropsClass interface{}            // :class
	PropsStyle map[string]interface{} // :style
	Attrs      map[string]string      // 本节点静态的attrs (除去class和style)
	Class      []string               // 本节点静态class
	Style      map[string]string      // 本节点静态style
	Slots      Slots                  // 当前组件所有的插槽代码(v-slot指令和默认的子节点), 支持多个不同名字的插槽, 如果没有名字则是"default"
	// 有两种情况
	// -  如果渲染的是元素（div等html元素），那么P是它所属的组件数据 ①
	// -  如果渲染的是组件，那么P是它的父级组件数据 ②
	// 在以下场景会用到 (后面的数字指的是属于上方的哪一种情况)
	// - 渲染插槽. (根据name取到所属组件的slot) ①
	// - 读取上层传递的PropsClass, 在root tag会读取上层的class等作用在自己身上. ①
	// - Inject ①
	// - Provide ①/②
	P             *Options
	Directives    directives // 多个指令
	VonDirectives []vonDirective
	// 组件模板中能够访问的所有值, 由Prototype+Props组成, 在指令中可以修改这个值达到声明变量的目的
	// tips: 由于渲染顺序, 修改只会影响到子节点
	Scope   *Scope
	Provide map[string]interface{}
}

func (o *Options) SetProvide(d map[string]interface{}) {
	if o.Provide == nil {
		o.Provide = d
	} else {
		o.Provide = map[string]interface{}{}
		for k, v := range d {
			o.Provide[k] = v
		}
	}
	return
}

// GetProvide会循环向上层查找Provide
func (o *Options) GetProvide(k string) (v interface{}) {
	// 向上查找
	curr := o
	for curr != nil {
		if curr.Provide != nil {
			if v, ok := curr.Provide[k]; ok {
				return v
			}
		}

		curr = curr.P
	}

	return nil
}

type directive struct {
	Name  string
	Value interface{}
	Arg   string
}

type vonDirective struct {
	Event string
	Func  string
	Args  []interface{}
}

type directives []directive

func (ds directives) Exec(r *Render, options *Options) {
	for _, d := range ds {
		if f, ok := r.directives[d.Name]; ok {
			f(DirectivesBinding{
				Value: d.Value,
				Arg:   d.Arg,
				Name:  d.Name,
			}, options)
		}
	}
}

type Props map[string]interface{}

func (p Props) CanBeAttr() Props {
	htmlAttr := map[string]struct{}{
		"id":  {},
		"src": {},
	}

	a := Props{}
	for k, v := range p {
		if _, ok := htmlAttr[k]; ok {
			a[k] = v
			continue
		}

		if strings.HasPrefix(k, "data-") {
			a[k] = v
			continue
		}
	}
	return a
}

type Slots map[string]NamedSlotFunc

func (s Slots) Exec(w Writer, name string, slotProps Props) {
	if s == nil {
		return
	}
	if f, ok := s[name]; ok {
		f(w, slotProps)
		return
	}

	return
}

// 组件的render函数
type ComponentFunc func(w Writer, options *Options)

// 用来生成slot的方法
// 由于slot具有自己的作用域, 所以只能使用闭包实现(而不是字符串).
type NamedSlotFunc func(w Writer, slotProps Props)

func (f NamedSlotFunc) Exec(w Writer, slotProps Props) {
	if f == nil {
		return
	}

	f(w, slotProps)
}

// 混合动态和静态的标签, 主要是style/class需要混合
// todo) 如果style/class没有冲突, 则还可以优化
// tip: 纯静态的class应该在编译时期就生成字符串, 而不应调用这个
// classProps: 支持 obj, array, string
// options: 上层组件的options
func mixinClass(options *Options, staticClass []string, classProps interface{}) (str string) {
	var class []string
	// 静态
	for _, c := range staticClass {
		if c != "" {
			class = append(class, c)
		}
	}

	// 本身的props
	for _, c := range getClassFromProps(classProps) {
		if c != "" {
			class = append(class, c)
		}
	}

	if options != nil {
		// 上层传递的props
		if options.PropsClass != nil {
			for _, c := range getClassFromProps(options.PropsClass) {
				if c != "" {
					class = append(class, c)
				}
			}
		}

		// 上层传递的静态class
		if len(options.Class) != 0 {
			for _, c := range options.Class {
				if c != "" {
					class = append(class, c)
				}
			}
		}
	}

	if len(class) != 0 {
		str = " class=\"" + strings.Join(class, " ") + "\""
	}

	return
}

// 构建style, 生成如style="color: red"的代码, 如果style代码为空 则只会返回空字符串
func mixinStyle(options *Options, staticStyle map[string]string, styleProps map[string]interface{}) (str string) {
	style := map[string]string{}

	// 静态
	for k, v := range staticStyle {
		style[k] = v
	}

	// 当前props
	ps := getStyleFromProps(styleProps)
	for k, v := range ps {
		style[k] = v
	}

	if options != nil {
		// 上层传递的props
		if options.PropsStyle != nil {
			ps := getStyleFromProps(options.PropsStyle)
			for k, v := range ps {
				style[k] = v
			}
		}

		// 上层传递的静态style
		for k, v := range options.Style {
			style[k] = v
		}
	}

	styleCode := genStyle(style)
	if styleCode != "" {
		str = " style=\"" + styleCode + "\""
	}

	return
}

// 生成除了style和class的attr
func mixinAttr(options *Options, staticAttr map[string]string, propsAttr map[string]interface{}) string {
	attrs := map[string]string{}

	// 静态
	for k, v := range staticAttr {
		attrs[k] = v
	}

	// 当前props
	ps := getStyleFromProps(propsAttr)
	for k, v := range ps {
		attrs[k] = v
	}

	if options != nil {
		// 上层传递的props
		if options.Props != nil {
			for k, v := range getStyleFromProps(options.Props.CanBeAttr()) {
				attrs[k] = v
			}
		}

		// 上层传递的静态style
		for k, v := range options.Attrs {
			attrs[k] = v
		}
	}

	c := genAttr(attrs)
	if c == "" {
		return ""
	}

	return " " + c
}

func getSortedKey(m map[string]string) (keys []string) {
	keys = make([]string, len(m))
	index := 0
	for k := range m {
		keys[index] = k
		index++
	}
	if len(m) < 2 {
		return keys
	}

	sort.Strings(keys)

	return
}

func genStyle(style map[string]string) string {
	sortedKeys := getSortedKey(style)

	var st strings.Builder
	for _, k := range sortedKeys {
		v := style[k]
		if st.Len() != 0 {
			st.WriteByte(' ')
		}
		st.WriteString(k + ": " + v + ";")
	}

	return st.String()
}

func genAttr(attr map[string]string) string {
	sortedKeys := getSortedKey(attr)

	var st strings.Builder
	for _, k := range sortedKeys {
		v := attr[k]
		if st.Len() != 0 {
			st.WriteByte(' ')
		}

		if v != "" {
			st.WriteString(k + "=" + "\"" + v + "\"")
		} else {
			st.WriteString(k)
		}
	}

	return st.String()
}

func getStyleFromProps(styleProps map[string]interface{}) map[string]string {
	st := map[string]string{}
	for k, v := range styleProps {
		switch v := v.(type) {
		case string:
			st[k] = escape(v)
		default:
			bs, _ := json.Marshal(v)
			st[k] = escape(string(bs))
		}
	}
	return st
}

// classProps: 支持 obj, array, string
func getClassFromProps(classProps interface{}) []string {
	if classProps == nil {
		return nil
	}
	var cs []string
	switch t := classProps.(type) {
	case []string:
		cs = t
	case string:
		cs = []string{t}
	case map[string]interface{}:
		var c []string
		for k, v := range t {
			if interfaceToBool(v) {
				c = append(c, k)
			}
		}
		sort.Strings(c)
		cs = c
	case []interface{}:
		var c []string
		for _, v := range t {
			cc := getClassFromProps(v)
			c = append(c, cc...)
		}

		cs = c
	}

	for i := range cs {
		cs[i] = escape(cs[i])
	}

	return cs
}

func lookInterface(data interface{}, keys ...string) (desc interface{}) {
	m, _, ok := shouldLookInterface(data, keys...)
	if !ok {
		return nil
	}

	return m
}

func lookInterfaceToSlice(data interface{}, key string) (desc []interface{}) {
	m, _, ok := shouldLookInterface(data, key)
	if !ok {
		return nil
	}

	return interface2Slice(m)
}

// 扩展map, 实现作用域
func extendMap(src map[string]interface{}, ext ...map[string]interface{}) (desc map[string]interface{}) {
	desc = make(map[string]interface{}, len(src))
	for k, v := range src {
		desc[k] = v
	}
	for _, m := range ext {
		for k, v := range m {
			desc[k] = v
		}
	}
	return desc
}

func extendScope(parent *Scope, data map[string]interface{}) *Scope {
	return &Scope{
		p:      parent,
		values: data,
	}
}

func interfaceToStr(s interface{}, escaped ...bool) (d string) {
	switch a := s.(type) {
	case int, string, float64:
		d = fmt.Sprintf("%v", a)
	default:
		bs, _ := json.Marshal(a)
		d = string(bs)
	}

	if len(escaped) == 1 && escaped[0] {
		d = escape(d)
	}
	return
}

// 字符串false,0 会被认定为false
func interfaceToBool(s interface{}) (d bool) {
	if s == nil {
		return false
	}
	switch a := s.(type) {
	case bool:
		return a
	case int, float64, float32, int8, int64, int32, int16:
		return a != 0
	case string:
		return a != "" && a != "false" && a != "0"
	default:
		return true
	}

	return
}

func interfaceToFloat(s interface{}) (d float64) {
	if s == nil {
		return 0
	}
	switch a := s.(type) {
	case int:
		return float64(a)
	case int32:
		return float64(a)
	case int64:
		return float64(a)
	case float64:
		return a
	case float32:
		return float64(a)
	default:
		return 0
	}
}

// 用来模拟js两个变量相加
// 如果两个变量都是number, 则相加后也是number
// 只有有一个不是number, 则都按字符串处理相加
func interfaceAdd(a, b interface{}) interface{} {
	an, ok := isNumber(a)
	if !ok {
		return interfaceToStr(a) + interfaceToStr(b)
	}
	bn, ok := isNumber(b)
	if !ok {
		return interfaceToStr(a) + interfaceToStr(b)
	}

	return an + bn
}

func interfaceLess(a, b interface{}) interface{} {
	an, ok := isNumber(a)
	if !ok {
		return interfaceToStr(a) < interfaceToStr(b)
	}
	bn, ok := isNumber(b)
	if !ok {
		return interfaceToStr(a) < interfaceToStr(b)
	}

	return an < bn
}

func interfaceGreater(a, b interface{}) interface{} {
	an, ok := isNumber(a)
	if !ok {
		return interfaceToStr(a) > interfaceToStr(b)
	}
	bn, ok := isNumber(b)
	if !ok {
		return interfaceToStr(a) > interfaceToStr(b)
	}

	return an > bn
}

func isNumber(s interface{}) (d float64, is bool) {
	if s == nil {
		return 0, false
	}
	switch a := s.(type) {
	case int:
		return float64(a), true
	case int32:
		return float64(a), true
	case int64:
		return float64(a), true
	case float64:
		return a, true
	case float32:
		return float64(a), true
	default:
		return 0, false
	}
}

// 用于{{func(a)}}语法
func interfaceToFunc(s interface{}) (d Function) {
	if s == nil {
		return emptyFunc
	}

	switch a := s.(type) {
	case func(args ...interface{}) interface{}:
		return a
	case Function:
		return a
	default:
		panic(a)
		return emptyFunc
	}
}

func interface2Slice(s interface{}) (d []interface{}) {
	switch a := s.(type) {
	case []interface{}:
		return a
	case []map[string]interface{}:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	case []int:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	case []int64:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	case []int32:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	case []string:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	case []float64:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	}
	return
}

// shouldLookInterface会返回interface(map[string]interface{})中指定的keys路径的值
func shouldLookInterface(data interface{}, keys ...string) (desc interface{}, rootExist bool, exist bool) {
	if len(keys) == 0 {
		return data, true, true
	}

	currKey := keys[0]

	switch data := data.(type) {
	case map[string]interface{}:
		// 对象
		c, ok := data[currKey]
		if !ok {
			return
		}
		rootExist = true
		desc, _, exist = shouldLookInterface(c, keys[1:]...)
		return

	case []interface{}:
		// 数组
		switch currKey {
		case "length":
			// length
			return len(data), true, true
		default:
			// index
			index, ok := strconv.ParseInt(currKey, 10, 64)
			if ok != nil {
				return
			}

			if int(index) >= len(data) || index < 0 {
				return
			}
			return shouldLookInterface(data[index], keys[1:]...)
		}
	case string:
		switch currKey {
		case "length":
			// length
			return len(data), true, true
		default:
		}
	}

	return
}

func escape(src string) string {
	return html.EscapeString(src)
}`