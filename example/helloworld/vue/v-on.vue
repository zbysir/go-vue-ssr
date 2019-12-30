<template>
  <div>
    <button v-on:click="buttonClick(msg)">click me to alert msg</button>
    <button @click="buttonClick2(msg)">click me to alert2 msg</button>
  </div>
</template>
