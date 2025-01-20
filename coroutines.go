package toolbox

import (
	"fmt"
	"runtime/debug"
)

// Recover 出错打印错误堆栈
func Recover() {
	if err := recover(); err != nil {
		fmt.Printf("Panic: %v\n%s", err, string(debug.Stack()))
		return
	}
}

// Try try-catch
func Try(f func()) {
	defer Recover()
	f()
}

// Go 带recover模式带go
func Go(f func()) {
	go func() {
		defer Recover()
		f()
	}()
}

// GoWithRecover wraps a `go func()` with recover()
func GoWithRecover(handler func(), recoverHandler func(r interface{})) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("goroutine panic: %v\n%s", r, string(debug.Stack()))
				if recoverHandler != nil {
					go func() {
						defer func() {
							if p := recover(); p != nil {
								fmt.Printf("goroutine panic: %v\n%s", r, string(debug.Stack()))
							}
						}()
						recoverHandler(r)
					}()
				}
			}
		}()
		handler()
	}()
}
