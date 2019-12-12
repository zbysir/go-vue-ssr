# Guide

Go-vue-ssr是简单的, 它的使用方法应该和Vue手册一样, 如果你不了解Vue可以直接去到[Vue官网](https://vuejs.org/)查阅资料.

如果在使用过程中, 你有遇到和Vue特性不一致的地方, 或者有需要优先支持的特性欢迎提Issue.

## Install
```
go get github.com/bysir-zl/go-vue-ssr
```
> 在这一步中可能会使用带代理, 如何为go get设置代理? 关键字: go get proxy

## Genera code(compile)
在你的项目下执行:
```
go-vue-ssr -src=./exaple/helloworld -to=./internal/vuetpl -pkg=vuetpl -pkg=vuetpl
```
```
$ go-vue-ssr -h

NAME:
   go-vue-ssr - Vue to Go compiler

USAGE:
   go-vue-ssr [global options] command [command options] [arguments...]

VERSION:
   0.0.1

DESCRIPTION:
   Hey vue go

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --src value    The .vue files dir
   --to value     Dist dir (default: "./internal/vuetpl")
   --pkg value    pkg name
   --help, -h     show help
   --version, -v  print the version

```
- src: 存放vue文件的文件夹, 支持子目录, 但不允许 

## Run
现在运行生成的代码就可以啦
```
r := vuetpl.NewRender()
html = r.Component_compName()
```

------

**下一篇: [Tips](tips.md)**
