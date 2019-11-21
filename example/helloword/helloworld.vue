<template>
  <div class="VueToNuxtLogo">
    <div class="Triangle Triangle--two">
      我是一个DIV
    </div>
    <div class="Triangle Triangle--one" />
    <div class="Triangle Triangle--three" />
    <div class="Triangle Triangle--four" />
    name: {{name}}
    info: <text v-bind:age="age" v-bind:sex="sex"></text>
  </div>
</template>