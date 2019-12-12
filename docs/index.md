# Go-vue-ssr
Vue server side render but golang. [https://bysir-zl.github.io/go-vue-ssr](https://bysir-zl.github.io/go-vue-ssr)

Hey vue go
## Cause
服务端渲染相较于前端渲染有以下好处:
- 利于内容型网站的SEO.
- 在性能更差的手机端浏览体验更佳.

而服务端渲染又有两个方向:
- 现代js框架vue/react所出的服务端渲染方案
- 传统的模板引擎, 如[raymond](https://github.com/aymerick/raymond)

各有优缺点
- js框架性能不好(这是后话了)
- 传统模板引擎在代码复杂的情况下并不美观, 在处理class/style方面也没有现代js框架方便.

由于代码洁癖与对于未知事物的好奇, 我初期还是选用的nuxt, 但后来发现它并不是银弹:
- 前后端同构需要[客户端激活](https://ssr.vuejs.org/zh/guide/hydration.html)步骤, 也就是在前端重新渲染一遍页面, 如果你的网站大多是静态的内容那么这一步就会造成很大的性能浪费(请不要小看客户端激活所带来的性能消耗).
- 由于不是专职note语言, 所以在面临高级问题上(如并发/缓存)举步维艰, 这对于后期发展不利.

> 关于vuessr性能问题可以看这两篇文章:
> - [实测Vue SSR的渲染性能：避开20倍耗时](https://mp.weixin.qq.com/s?__biz=MzUxMzcxMzE5Ng==&mid=2247485601&amp;idx=1&amp;sn=97a45254a771d13789faed81316b465a&source=41#wechat_redirect)
> - [How to Drastically Reduce Estimated Input Latency and Time to Interactive of SSR Vue.js Applications](https://markus.oberlehner.net/blog/how-to-drastically-reduce-estimated-input-latency-and-time-to-interactive-of-ssr-vue-applications/)

当面临现实的效率问题时, 不得不妥协而使用传统的模板引擎, 但他们实际不是专为现代html而生, 所以都不如vue模板好用(看).

有什么办法能让喜欢vue和go的你我更舒心的编写代码呢?

这就是这个项目诞生的原因.

它将尽力保留vue的特性, 如组件化, [Custom Directives](https://vuejs.org/v2/guide/custom-directive.html), [Class and Style Bindings](https://vuejs.org/v2/guide/class-and-style.html), 相信这些现代特性对于编写html代码是有利的.

## Who need Go-vue-ssr
项目的目的是高效渲染+优雅的模板语法, 并没有实现vue的js部分的特性,
所以它更适用于如官网/活动页等功能不强的页面, 而不适用于如后台管理系统这样功能性强的系统.

## Feature
基于字符串拼接 而不是 虚拟节点来渲染vue组件, 当然这样做有好有坏.

好处就是性能至少能提升1个数量级, 坏处就是舍去虚拟节点也就无法实现vue的数据绑定特性.

## Example
> 完整代码[在这](https://github.com/bysir-zl/go-vue-ssr/tree/master/example/helloworld)

编写vue组件代码如下: (仅支持template块)
```vue
<template>
  <div style="text-align: center">
    <p v-text="slogan" style="padding: 10px 0"></p>
    <img height="50px" alt="todo logo" :src="logo">
  </div>
</template>
```
执行go-vue-ssr编译vue模板
```sh
go-vue-ssr -src=./vue -to=./ -pkg=main
```
编写调用代码
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

## Usage

### step 1: install
```
go get github.com/bysir-zl/go-vue-ssr
```
### step 2: genera
```
go-vue-ssr -src=./exaple/helloworld -to=./internal/vuetpl -pkg=vuetpl -pkg=vuetpl
```
此命令将在./internal/vuetpl里生成go代码

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

> 更多细节请查看文档: [编译](genera.md)

### step3: run
生成的代码可以直接运行返回html.

```go
r := vuetpl.NewRender()
html = r.Component_helloworld()
```

## Supported Vue Template Syntax
- [Text](https://vuejs.org/v2/guide/syntax.html#Text)
  - mustache syntax (double curly braces)
  - v-text (use html.escape)
- [Raw Html](https://vuejs.org/v2/guide/syntax.html#Raw-HTML)
  - v-html
- [Attributes](https://vuejs.org/v2/guide/syntax.html#Attributes)
  - v-bind (support shorthands)
- [Arguments](https://vuejs.org/v2/guide/syntax.html#Attributes)
  - v-bind (support shorthands)
- [Custom Directives](https://vuejs.org/v2/guide/custom-directive.html)
  - emm it's different with vue's custom Directives, see [Tips-CustomDirectives](docs/tips.md#CustomDirectives)
- Class and Style Bindings
  - [Object-Syntax](https://vuejs.org/v2/guide/class-and-style.html#Object-Syntax)
  - [Array Syntax](https://vuejs.org/v2/guide/class-and-style.html#Array-Syntax)
  - [With-Components](https://vuejs.org/v2/guide/class-and-style.html#With-Components)
- [Conditional Rendering](https://vuejs.org/v2/guide/conditional.html)
  - v-if
  - v-else-if
  - v-else
- [List Rendering](https://vuejs.org/v2/guide/list.html)
  - v-for (only on Array, not support Object/Range)
- [Slots](https://vuejs.org/v2/guide/components-slots.html)
  - [Compilation Scope](https://vuejs.org/v2/guide/components-slots.html#Compilation-Scope)
  - [Fallback Content](https://vuejs.org/v2/guide/components-slots.html#Fallback-Content)
  - [Named Slots](https://vuejs.org/v2/guide/components-slots.html#Named-Slots)
  - [Scoped Slots](https://vuejs.org/v2/guide/components-slots.html#Scoped-Slots)
- [Dynamic Components](https://vuejs.org/v2/guide/components-dynamic-async.html)

- Using JavaScript Expressions (by AST)
  - `+ && || !`
  - `function call` e.g. \{\{calcHeight(srcHeight)}}
  - `.length`
  - `'list-' + id`
**not support**
- v-on
- v-show
- filter: please use function instead of it, e.g. \{\{calcHeight(srcHeight)}}
- inject / provider
- v-once

**other**
- prototype: 放在Prototype里的变量可以在任何组件中使用, 如调用全局的方法.

------

**下一篇: [手册](guide.md)**
