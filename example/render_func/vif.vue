<template>
  <div >
    <div v-if="name" :html="1" v-bind:a="1">{{name}}</div>
    <div v-else>empty name</div>
  </div>
</template>
