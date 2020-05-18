<template>
  <div v-bind:class="{'a': true}" class="b">
    <span class="d" v-bind:class="{c: true}" :a="1" v-html="data.msg">
    </span>
    <div v-if="a"></div>

    <template v-else v-for="id in a">
      <div></div>
    </template>



    <!-- 子组件将会被异步渲染 -->
    <div v-for="item in data.c">
      <async>
        <bench :data="item"></bench>
      </async>
    </div>
  </div>
</template>