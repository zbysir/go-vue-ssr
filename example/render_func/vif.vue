<template>
  <div>
    <div v-if="name" :html="1" v-bind:a="1">{{name}}</div>
    <div v-else-if="name2">{{name2}}</div>
    <template v-else> !name AND !name2</template>
    <div>
      <div v-slot:a>
        vsa
      </div>
    </div>
  </div>
</template>
