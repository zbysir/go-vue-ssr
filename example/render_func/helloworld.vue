<template>
  <div class="VueToNuxtLogo">
    <div class="Triangle Triangle--two">
      我是一个DIV
    </div>
    <span v-for="item in list">{{item}}</span>
    name: {{name}}
    info:
    <text v-bind:age="age" v-bind:sex="sex">
      我的信息:
    </text>
  </div>
</template>

