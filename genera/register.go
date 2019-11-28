package genera

var components = map[string]ComponentFunc{}
func init(){components = map[string]ComponentFunc{"component": XComponent_component,"slot": XComponent_slot,"class": XComponent_class,"helloworld": XComponent_helloworld,"text": XComponent_text,"vFor": XComponent_vFor,"xslot": XComponent_xslot,"xstyle": XComponent_xstyle,}}