package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	sourceFile := "./generotor_builtin_code/source.go"
	target := "./generator_builtin_code_gen.go"
	pkg := "vuessr"
	beginTag := []byte("// begin")

	source, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		panic(err)
	}

	if !bytes.Contains(source, beginTag) {
		panic("source file not has `// begin` tag")
	}

	source = bytes.Split(source, beginTag)[1]

	to := fmt.Sprintf(`// generate by ./generotor_builtin_code/main.go
package %s

const builtinCode = `+"`%s`", pkg, source)

	err = ioutil.WriteFile(target, []byte(to), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
