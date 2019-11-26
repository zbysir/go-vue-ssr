<template>
  <p>
    age: {{age}} sex: {{sex}}
    <vFor v-bind:list="list">
      <slot></slot>
    </vFor>
  </p>
</template>
