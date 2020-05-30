package parser

import (
	"encoding/json"
	"testing"
)

func TestGoHtmlParse(t *testing.T) {
	p := GoHtml{}
	x, err := p.Parse(`./test_src/template_in_head.html`)

	if err != nil {
		t.Fatal(err)
	}

	bs, _ := json.MarshalIndent(x, " ", " ")
	t.Logf("%s", bs)
}
