<template>
    <div>
        <span v-for="(item, index) in list">{{index}}: {{item}} <slot></slot></span>
    </div>
</template>