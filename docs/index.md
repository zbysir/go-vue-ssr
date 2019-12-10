使用golang渲染vue

## cause

项目诞生的目的是为了解决现有node项目的性能问题 与 [客户端激活](https://ssr.vuejs.org/zh/guide/hydration.html)的性能问题
- https://markus.oberlehner.net/blog/how-to-drastically-reduce-estimated-input-latency-and-time-to-interactive-of-ssr-vue-applications/
- https://mp.weixin.qq.com/s?__biz=MzUxMzcxMzE5Ng==&mid=2247485601&amp;idx=1&amp;sn=97a45254a771d13789faed81316b465a&source=41#wechat_redirect

虽然vuessr有优缺点, 但我认为vue的模板格式清晰易懂, 用来替换golang的tpl引擎或者其他模板引擎(如raymond)是有价值的.

总的来说, 如果你的项目追求的是性能, 对于vue特性需求不大或者压根不需要vue特性, 那你就可以试一试这个项目.

## feature
基于字符串拼接 而不是 虚拟节点来渲染vue组件, 当然这样做有好有坏.

好处就是性能至少能提升1个数量级, 坏处就是舍去虚拟节点也就无法实现vue的数据绑定特性.

## usage

### step 1: install
```
go get github.com/bysir-zl/go-vue-ssr
```
### step 2: genera
```
go-vue-ssr -src=./exaple/helloworld -to=./internal/vuetpl
```
将在./internal/vuetpl里生成go代码

所有运行渲染所需要的代码都会保存在vuetpl包里, 也就是运行时不会依赖github.com/bysir-zl/go-vue-ssr包, 

不过在github.com/bysir-zl/go-vue-ssr/pkg/ssrtool里有一些处理动态数据(interface{})的工具方法可以使用, 如
```
a:= map[string]interface{}{
    "info": map[string]interface{}{
        "name": "bysir",
    },
}

// 使用LookInterface方法可以方便的得到a.info.name的值.
ssrtool.LookInterface(a, "info.name")
``` 

### step 3: run
```go
r:=vuetpl.NewRender()
html = r.XComponent_helloworld()
```

## vue features
**support**
- v-if v-else v-else-if
- v-for
- v-bind (support shorthands)
- dynamically style
- dynamically class
- v-slot
- slot scope
- component
- expression by AST
  - `+ && || !`
  - `function call`
  - `.length`
- function call: eg. {{calcHeight(srcHeight)}}
- directive
- v-html (use html.escape)
- v-text

**not support**
- v-on
- v-show
- filter: please use function instead of it, e.g. {{calcHeight(srcHeight)}}
- inject / provider

**other**
- prototype: 放在Prototype里的变量可以在任何组件中使用.

## 编译原理

### 处理vue模板
vue的模板其实是标准的html.

所以使用golang.org/x/net/html包解析HTML, 得到Token之后再根据attr处理vue特殊的指令, 如v-if v-for, 最终得到vue节点.

### 处理js
在v-if或者{{}}中需要使用一些简单的js表达式, 如 v-if="a!=b && a!=c", 这样的表达式需要翻译成golang才能运行, 翻译成golang需要使用到js的AST,

最开始的想法是使用golang实现或者找一个现有的库去解析JS, 但奈何没有找到, 实现起来也十分麻烦, 所以还是使用Node+acorn封装了一个API供以调用.

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
