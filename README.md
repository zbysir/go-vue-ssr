# vue-ssr
vue server side render but golang

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
vue-ssr -src=./exaple/helloworld -to=./internal/vuetpl
```
将在./internal/vuetpl里生成go代码

### step3: run
```go
html = vuetpl.XComponent_helloworld()
```

## vue features
**support**
- v-if
- v-for
- v-bind (support shorthands)
- dynamically style
- dynamically class
- named slot
- slot scope
- component
- expression by AST
  - `+ && || !`
  - `function call`
  - `.length`  

**not support**
- v-on
- v-show

**todo**
- v-else
- v-else-if
- inject / provider
