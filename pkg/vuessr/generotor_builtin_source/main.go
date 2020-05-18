package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	sourceFiles := []string{"./generotor_builtin_source/source.go"}
	target := "./generator_builtin_gen.go"
	pkg := "vuessr"
	beginTag := []byte("// begin")

	source := ""
	for _, sourceFile := range sourceFiles {
		fileBs, err := ioutil.ReadFile(sourceFile)
		if err != nil {
			panic(err)
		}

		if !bytes.Contains(fileBs, beginTag) {
			panic("source file not has `// begin` tag")
		}

		// 合并多个文件中的import

		code := bytes.Split(fileBs, beginTag)[1]
		code = bytes.Trim(code, "\n")

		source += fmt.Sprintf("\n\n// src: %s\n%s", sourceFile, code)
	}

	to := fmt.Sprintf(`// generate by ./generotor_builtin_source/main.go
package %s

const builtinCode = `+"`%s`", pkg, source)

	err := ioutil.WriteFile(target, []byte(to), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
