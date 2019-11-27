package vuessr

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 用来生成模板字符串代码
// 目的是为了解决递归渲染节点造成的性能问题, 不过这是一个难题, 先尝试, 不行就算了.

func genComponentRenderFunc(app *App, pkgName, name string, file string) string {
	ve, err := ParseVue(file)
	if err != nil {
		panic(err)
	}
	code, _ := ve.RenderFunc(app)

	// 处理多余的纯字符串拼接: "a"+"b" => "ab"
	//code = strings.Replace(code, `"+"`, "", -1)

	return fmt.Sprintf("package %s\n\n"+
		"func XComponent_%s(options *Options)string{\n"+
		"%s:= %s\n_ = %s\n"+
		"return %s"+
		"}", pkgName, name, DataKey, PropsKey, DataKey, code)
}

// 生成并写入文件夹
func genAllFile(src, desc string) (err error) {
	// 生成文件夹
	err = os.MkdirAll(desc, os.ModePerm)
	if err != nil {
		return
	}

	// 删除老的.vue.go文件

	del, err := walkDir(desc, ".vue.go")
	if err != nil {
		return
	}

	for _, v := range del {
		err = os.Remove(v)
		if err != nil {
			return
		}
	}

	// 生成新的
	vue, err := walkDir(src, ".vue")
	if err != nil {
		return
	}

	app := NewApp()

	for _, v := range vue {
		_, fileName := filepath.Split(v)
		name := strings.Split(fileName, ".")[0]
		app.Component(name)
	}

	for _, v := range vue {
		_, fileName := filepath.Split(v)
		name := strings.Split(fileName, ".")[0]
		_, pkgName := filepath.Split(desc)
		code := genComponentRenderFunc(app, pkgName, name, v)
		err = ioutil.WriteFile(desc+string(os.PathSeparator)+fileName+".go", []byte(code), 0666)
		if err != nil {
			return
		}
	}

	return
}

func walkDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)

	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		//遍历目录
		if err != nil {
			return err
		}
		if fi.IsDir() {
			// 忽略目录
			return nil
		}

		if strings.HasSuffix(filename, suffix) {
			files = append(files, filename)
		}

		return nil
	})

	return
}
