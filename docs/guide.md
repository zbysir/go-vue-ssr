# Guide
这一篇将介绍本项目如何使用.

## Install
```
go get github.com/zbysir/go-vue-ssr
```
> 在这一步中可能会使用带代理, 如何为go get设置代理? 关键字: go get proxy

## Genera code(compile)
在你的项目下执行:
```
go-vue-ssr -src=./exaple/helloworld -to=./ -pkg=./ -pkg=vuetpl
```
根据vue文件不同, 生成的代码可能如下
```
│  buildin.go -- 内置的方法
│  info.vue.go -- 组件渲染函数
│  new.go -- 注册component的代码
```

> tips: 你可以使用 `-watch`参数来启用监听文件变化自动编译.

关于如何使用go-vue-ssr命令请看 [生成](genera.md)

## Run
现在运行生成的代码就可以啦
```go
package main

func main()  {
    r := NewRender()
    htmlStr := r.Component_page(&Options{
    	Props: map[string]interface{}{
    		"slogan": "Hey vue go",
    		"logo":   "https://avatars2.githubusercontent.com/u/13434040?s=88&v=4",
    	},
    })
    print(htmlStr)
}
```

现在你就可以编写简单的vue组件了, Go-vue-ssr是简单的, 组件的编写方法应该和Vue手册一样, 如果你不了解Vue可以直接去到[Vue官网](https://vuejs.org/)查阅资料.

如果在使用过程中, 你有遇到和Vue特性不一致的地方, 或者有需要优先支持的特性欢迎提Issue.

不过既然你都读到了这里, 下一篇[Tips](tips.md)也应该读一读, 它会介绍一些小问题与技巧.

------

**下一篇: [Tips](tips.md)**
