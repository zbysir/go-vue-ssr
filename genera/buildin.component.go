package genera
// 内置组件

func XComponent_slot(options *Options)string {
	name:=options.Attrs["name"]
	props:=options.Props
	injectSlotFunc:= options.P.Slot[name]

	// 如果没有传递slot 则使用默认的code
	if injectSlotFunc == nil {
		return options.Slot["default"](nil)
	}

	return injectSlotFunc(props)
}

func XComponent_component(options *Options)string{
	is,ok:=options.Props["is"].(string)
	if !ok{
		return ""
	}
	return components[is](options)
}