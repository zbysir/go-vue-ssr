package ssrtool

import (
	"testing"
)

func TestFormat(t *testing.T) {
	s := FormatHtml("<div> 123 <p>456</p><div>444</div></div><div>x</div>", 2)
	//bs,_:=json.Marshal(s)
	t.Logf("'%s'", s)
}
