<!DOCTYPE html>
<html lang="zh">
<head>
  <meta charset="UTF-8">
  <title>{{title}}</title>
</head>
<body>
<h1 v-html="title" style="text-align: center; margin-top: 100px"></h1>
<info :name="title" :slogan="slogan" :logo="logo" style="padding: 20px" :height="height+1"></info>
<div info> author: {{info.author}} {{info[slogan]}} {{[slogan]}}</div>

</body>
</html>