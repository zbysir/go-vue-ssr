# vue-ssr
vue server side render but golang

## cause
- https://markus.oberlehner.net/blog/how-to-drastically-reduce-estimated-input-latency-and-time-to-interactive-of-ssr-vue-applications/
- https://mp.weixin.qq.com/s?__biz=MzUxMzcxMzE5Ng==&mid=2247485601&amp;idx=1&amp;sn=97a45254a771d13789faed81316b465a&source=41#wechat_redirect

## feature
基于字符串拼接 而不是 虚拟节点来渲染vue组件, 这样做当然是有好有坏的

好处就是性能至少能提升1个数量级, 坏处就是舍去了vue的数据绑定特性.

但如果实现vue的数据绑定就一定会有虚拟节点, 这又会导致性能问题(如nuxt的客户端激活).

如果你的项目追求性能, 舍弃vue的特性也不是不能接受.

## usage

### step 1: install
```
go get github.com/bysir-zl/vue-ssr
```
### step 2: genera
```
vue-ssr -src=./exaple/helloworld -to=./internal/vuetpl
```
将在./internal/vuetpl里生成go代码

### step3: run
```go
vuetpl.XComponent_helloworld()
```

## vue features
**support**
- v-if
- v-for
- v-bind
- dynamically style
- dynamically class
- named slot
- slot scope
- component

**not support**
- v-on
- v-show

**todo**
