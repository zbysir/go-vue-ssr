# genera

渲染vue的方法有两个方向:

其一和nuxt一样 基于vue官方提供的vuessr虚拟节点的渲染方式, 它支持vue的全特性, 但性能仅能满足不太复杂的网站.

其二就是将vue再编译成更高效的代码运行 也就是传统的基于字符串拼接的模板渲染方式, 这种方式能有效避免节点太多所造成的递归/动态等性能问题.

项目的目标就是高效渲染+优雅的模板语法, 故使用上述第二个方法.

## go-vue-ssr命令
使用go-vue-ssr命令可以生成代码
```
go-vue-ssr -src=./exaple/helloworld -to=./ -pkg=./ -pkg=vuetpl
```

#### 命令说明
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
**参数**

- src: 存放vue文件的文件夹, 支持查找子目录, 但不允许重复的文件名(因为文件名会当做组件名).
- to: 存放生成代码的目录
- pkg: go package name

此命令将在当前目录下生成所有需要的Go代码, 也就是运行时不会依赖github.com/zbysir/go-vue-ssr包.

不过在github.com/zbysir/go-vue-ssr/pkg/ssrtool里有一些处理动态数据(interface{})的工具方法可以使用, 如
```
a:= map[string]interface{}{
    "info": map[string]interface{}{
        "name": "bysir",
    },
}

// 使用LookInterface方法可以方便的得到a.info.name的值.
rinterface.GetStr(a, "info.name")
```

## 编译原理

### 处理vue模板
vue的模板其实是标准的html.

所以使用golang.org/x/net/html包解析HTML, 得到html节点树之后再根据attr处理vue特殊的指令, 如v-if v-for, 最终得到vue节点.

### 处理js
在v-if或者\{\{}}中需要使用一些简单的js表达式, 如 v-if="a!=b && a!=c", 这样的表达式需要翻译成golang才能运行, 翻译成golang需要使用到js的AST,

使用golang的库`github.com/robertkrimen/otto`实现解析js代码, 没有使用node+web server实现的原因是内联的golang库性能更好, 但缺点是不支持ES6的高级语法, 如`{[a]: 1}`,
请避免在模板中使用这些高级语法.

### 动态节点 / 静态节点 / 半动态节点
**静态节点**
静态节点在编译时就会生成静态的字符串.

如
```
<span class="m"></span>
```
在最终生成的代码中是这样:
```
"<span class=\"m\"></span>"
```

**动态节点**
动态节点的生成发生在运行时.

满足以下条件都是动态节点:
- 拥有`指令`: 由于指令中可以修改节点的属性, 只能在运行时动态生成html
- `组件的root节点`

动态节点使用方法生成: r.Tag(tagName, options).

拥有`指令` 或者 是`组件的root节点` 则统一为动态节点

**半动态节点**
带有 动态class/style/attr的节点由于需要在运行时确定class/style/attr, 但由于也只需要修改这些属性, 所以最终生成的代码是
```
<div + mixinClass() + mixinStyle + mixinAttr()>children...</div>
```

半动态节点相比动态节点少了方法的调用, 性能会更好一些.

------

**[回到首页](.)**