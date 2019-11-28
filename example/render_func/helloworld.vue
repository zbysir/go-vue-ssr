<template>
  <div v-if="isShow" class="VueToNuxtLogo">
    <div
        class="Triangle Triangle--two"
        v-bind:class="customClass"
        style="background: #f99; ">
      我是一个DIV
    </div>

    name: {{name + " " + name}}
    info:

    <text v-bind:list="list">
      <div v-slot:abc="a">
        1我是具名slot props msg: {{a.msg}}
        1我是具名slot 所属组件属性 age: {{age}}
      </div>
    </text>
    <text v-bind:list="list">
      <div v-slot:abc="a">
        2我是具名slot props msg: {{a.msg}}
        2我是具名slot 所属组件属性 age: {{age}}
      </div>
    </text>

    <div v-if="!isShow">显示隐藏</div>

    <class v-bind:class="customClass" class="d e f"></class>
    <xstyle v-bind:style="{'padding': '50px'}" style="margin: 50px"></xstyle>
  </div>
</template>

