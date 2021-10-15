package main

import (
	"fmt"
	"time"
	"unsafe"
)

func main() {
	// 常用类型的宽度
	var i int
	var s string
	var b bool
	var ai []int32
	var as []string
	var mb map[string]bool
	var mi map[string]interface{}
	var ms map[string]struct{}
	var st struct{}

	fmt.Printf("int sizeof %v\n",unsafe.Sizeof(i)) // 8
	fmt.Printf("string sizeof %v\n",unsafe.Sizeof(s)) // 16
	fmt.Printf("bool sizeof %v\n",unsafe.Sizeof(b)) // 1
	fmt.Printf("[]int32 sizeof %v\n",unsafe.Sizeof(ai)) // 24
	fmt.Printf("[]string sizeof %v\n",unsafe.Sizeof(as)) // 24
	fmt.Printf("map[string]bool sizeof %v\n",unsafe.Sizeof(mb)) // 8
	fmt.Printf("map[string]interface{} sizeof %v\n",unsafe.Sizeof(mi)) // 8
	fmt.Printf("map[string]struct{} sizeof %v\n",unsafe.Sizeof(ms)) // 8

	fmt.Printf("struct{} sizeof %v\n",unsafe.Sizeof(st)) // 0 用于占位

	var t T
	t.Hello()

	set := Set{}
	set.Append("aaa")
	set.Append("aaa")
	set.Append("bbb")
	set.Append("ccc")
	set.Remove("aaa")
	fmt.Println(set.Exits("aaa")) // false

	// 使用场景3 空通道
	// 实现 timeout的处理
	ch := make(chan struct{})
	go func() {
		time.Sleep(10 * time.Second)
		close(ch)
	}()
	fmt.Println("aaaaaaa")
	<-ch
	fmt.Println("ooooooo")
}

// 使用场景1 方法的接收者
type T struct {}

func (t *T) Hello()  {
	fmt.Println("hello empty struct")
}

// 使用场景2 实现集合类型
type Set map[string]struct{}

func (s Set) Append(k string)  {
	s[k] = struct{}{}
}

func (s Set) Remove(k string)  {
	delete(s,k)
}

func (s Set) Exits(k string) bool {
	_,ok := s[k]
	return ok
}