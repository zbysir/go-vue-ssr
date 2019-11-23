<template>
  <div class="VueToNuxtLogo">
    <div class="Triangle Triangle--two">
      我是一个DIV
    </div>
    name: {{name}}
    info: <text v-bind:age="age" v-bind:sex="sex"></text>
  </div>
</template>


"<div>....." + componentText(map[string]interface{"age": look(data, "age")}, slot)+"</div>"


func componentText(data map[string]interface, slot)(string){
return "<p>" + look(data, "sex") + slot +"</p>"
}
