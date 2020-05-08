package version

// 当version改变，vue编译缓存就会失效。
const Version = "0.0.17"

// 0.0.9
// fix <!doctype html>

// 0.0.10
// fix unsafe string in attr

// 0.0.11
// use github.com/robertkrimen/otto to parse js code

// 0.0.12
// support watch file and recompile
// use the next package to watch file: github.com/radovskyb/watcher

// 0.0.13
// Optimization code: scope

// 0.0.14
// support inject and provide
// use directive: v-provide AND v-inject

// 0.0.15
// fix empty slot

// 0.0.16
// fix panic when nil slot called
