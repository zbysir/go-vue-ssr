<!DOCTYPE html>
<html lang="zh">
<head>
  <meta charset="UTF-8">
  <title>{{title}}</title>
</head>
<body>
<h1 v-html="title" style="text-align: center; margin-top: 100px"></h1>
<info :name="title" :slogan="slogan" :logo="logo" style="padding: 20px" :height="height+1"></info>
<v-on :msg="'hello event'"></v-on>

<script v-on-handler/>
<script>
    // 为dom添加事件
    for (var i in vOnBinds) {
        var item = vOnBinds[i];
        var dom = document.querySelector('[data-von-' + item.DomSelector + ']')
        dom.addEventListener(item.Event, function (item, dom) {
            return function (event) {
                if (window[item.Func]) {
                    window[item.Func].call(window, event, ...item.Args)
                } else {
                    console.error('not found function: ' + item.Func)
                }
            }
        }(item, dom))
    }
</script>

</body>
</html>