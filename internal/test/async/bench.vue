<template>
  <div v-bind:class="{'a': true}" class="b">
    <span class="d" v-bind:class="{c: true}" :a="1">
        {{data.msg}}
    </span>

    <!-- 子组件将会被异步渲染 -->
    <async>
      <div v-for="item in data.c">
        <bench :data="item"></bench>
      </div>
    </async>
  </div>
</template>