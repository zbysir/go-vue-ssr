<!DOCTYPE html>
<html :lang="lang" :data-site-id="siteId">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport"
        content="width=device-width,initial-scale=1.0,minimum-scale=1.0,maximum-scale=1.0,user-scalable=no">
  <title>{{title}}</title>
  <!--动画-->
  <link rel="stylesheet" href="//static.f.cdn-static.cn/3.7.0/animate.min.css" type="text/css"/>
  <!--swiper-->
  <link rel="stylesheet" href="//static.f.cdn-static.cn/swiper/swiper.min.css">
  <link rel="stylesheet" href="/dist/fonts/iconfont.css" type="text/css"/>
  <link rel="stylesheet" :href="'/dist/css/styles.'+hash['styles.css']+'.css'" type="text/css"/>
  <!--全屏滚动-->
  <link rel="stylesheet" href="//static.f.cdn-static.cn/fullpage/3.0.5/fullpage.css" type="text/css"/>
  <!-- 页面样式 -->
  <link rel="stylesheet" :href="pageCssUrl" type="text/css"/>

  <!-- 网站样式 -->
  <style v-html="siteCss"></style>
</head>

<body
  id="render"
  class="render"
  :class="{'full-page-body': pageLayoutSetting.isFullPageScroll}"
  v-mount:react="{react: siteRaw.react, cp: 'site'}"
>
<div id="head-code" v-html="siteRaw.code.header"></div>

<script type="text/javascript" src="//static.f.cdn-static.cn/swiper/swiper.min.js"></script>
<script type="text/javascript" src="//static.f.cdn-static.cn/popper/popper.min.js"></script>
<script type="text/javascript" src="//static.f.cdn-static.cn/axios/0.18.0/axios.min.js"></script>
<script type="text/javascript" src="//static.f.cdn-static.cn/lodash.js/4.17.10/lodash.min.js"></script>
<script type="text/javascript" src="//static.f.cdn-static.cn/wow/wow.min.js"></script>
<script type="text/javascript" src="//static.f.cdn-static.cn/fullpage/3.0.5/fullpage.extensions.min.js"></script>

<script :src="'/dist/js/main.'+hash['main.js']+'.js'"></script>

<script v-get="{'mount':'pageMount'}">
  window.$mount(pageMount)
</script>

<script src="//static.f.cdn-static.cn/lazysizes.min.js" async></script>
</body>
</html>