package main

import (
	"fmt"
	"runtime"
)

// call (skip int)
// - skip : 是要提升的堆栈帧数，0-当前函数，1-上一层函数，....
// pc - 是uintptr这个返回的是函数指针
// file - 是函数所在文件名目录
// line - 所在行号
// ok - 是否可以获取到信息
func call(skip int) {
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		funcName := runtime.FuncForPC(pc).Name() //获取函数名
		fmt.Println(fmt.Sprintf("skip=%v  pc=%v file=%s   line=%d  funcName=%s", skip, pc, file, line, funcName))
	}
}

func test(skip int) {
	call(skip)
}

func main() {
	for i := 0; i < 4; i++ {
		test(i)
	}
}
