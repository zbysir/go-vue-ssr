// Code generated by go-vue-ssr: https://github.com/bysir-zl/go-vue-ssr

package bench_string

func NewRender() *Render {
	r := &Render{}
	r.Components = map[string]ComponentFunc{
		"bench":     r.Component_bench,
		"component": r.Component_component,
		"slot":      r.Component_slot,
	}
	return r
}