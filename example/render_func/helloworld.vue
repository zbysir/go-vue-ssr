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

    <text
        v-bind:age="age"
        v-bind:sex="sex"
        v-bind:data="{a: 1}"
        v-bind:list="list">
      我的信息:
      <div>-- defaultSlot --</div>

      <div v-for="item in list">
        {{item}}
      </div>

      <div>-- vFor --</div>
      <vFor v-bind:list="list">
      </vFor>
      <div>-- /vFor --</div>

      <div>-- /defaultSlot --</div>
    </text>

    <text>
      <div v-slot:abc="a">
        <div>-- 具名slot --</div>
        我是具名slot内容
        <div>-- /具名slot --</div>
      </div>
    </text>

    <div v-if="!isShow">显示隐藏</div>

    <class v-bind:class="customClass"></class>
  </div>
</template>

