<template>
    <div v-bind:id="'id'" class="b" data-src="//baidu.jpg">
        test attr
        <img v-bind:src="img(imgUrl)" alt="标题"/>
    </div>
</template>