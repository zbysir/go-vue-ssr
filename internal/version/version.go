package version

// 当version改变，编译缓存就会失效。
const Version = "0.0.10"

// 0.0.9
//  fix <!doctype html>

// 0.0.10
// fix unsafe string in attr

// 0.0.11
// use github.com/robertkrimen/otto to parse js code