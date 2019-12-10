package ssrtool

import "testing"

func TestLookJson(t *testing.T) {
	t.Logf("%v", LookJson([]byte(`{"a": 1}`), "a") == 1.0)
}
