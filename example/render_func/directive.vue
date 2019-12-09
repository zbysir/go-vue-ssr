<template>
  <div v-animate="{time: '5s', xclass:xclass}">
    <template v-if="show">
      <p v-animate="{type: 'up'}">
        test animate
      </p>
    </template>
    <template v-else>
      !show
    </template>

    <p v-set:swiper='{id: id, value: {"swiper": {"speed": speed}}}'></p>

    <p v-get></p>
    333
  </div>
</template>
