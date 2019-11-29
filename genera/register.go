package genera

var components = map[string]ComponentFunc{}
func init(){components = map[string]ComponentFunc{"xstyle": XComponent_xstyle,"vFor": XComponent_vFor,"xattr": XComponent_xattr,"slot": XComponent_slot,"class": XComponent_class,"text": XComponent_text,"xslot": XComponent_xslot,"component": XComponent_component,"helloworld": XComponent_helloworld,}}