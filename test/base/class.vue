<template>
  <div v-bind:class="{'a': true}" class="b">
        <span class="d" v-bind:class="{c: true}" :a="1">
            test class
        </span>
  </div>
</template>