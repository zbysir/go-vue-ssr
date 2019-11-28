<template>
  <div>
    <div v-for="item in list">
      <slot name="abc" v-bind:msg="item">
        备选slot内容
      </slot>
      <slot name="abcd">
        备选slot abcd内容
      </slot>
    </div>
  </div>
</template>
