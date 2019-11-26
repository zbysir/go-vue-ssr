<template>
  <div v-if="isShow" class="VueToNuxtLogo">
    <div
        class="Triangle Triangle--two"
        v-bind:class="customClass"
        style="background: #f99; ">
      我是一个DIV
    </div>

    name: {{name}}
    info:

    <text v-bind:age="age" v-bind:sex="sex" v-bind:list="list">
      我的信息:
      <span v-for="item in list">
        {{item}}
      </span>
    </text>

    <div v-if="isShow">显示隐藏</div>
  </div>
</template>

