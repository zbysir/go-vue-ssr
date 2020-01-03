<template>
  <div v-if="name" :html="1" v-bind:a="1">
    {{name}}
    <div v-if="name2">{{name2}}</div>
    <div v-else>123123</div>
  </div>
  <div v-else-if="name2"> !name AND !name2</div>
  <div v-else> !name AND !name2</div>
</template>
