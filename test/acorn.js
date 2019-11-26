// run "cnpm i acorn" first
const {Parser} = require("acorn")

console.log(JSON.stringify(Parser.parse("(a + b) || v"), " ","    "))
