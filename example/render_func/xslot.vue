<template>
  <div>
    <text v-bind:list="list">
      <div v-slot:abc="a">
        我是具名slot props msg: {{a.msg}}
        我是具名slot 所属组件属性 age: {{age}}
      </div>
    </text>
  </div>
</template>
