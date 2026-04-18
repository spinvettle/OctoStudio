package main

import (
	"fmt"
	"unsafe"
)

func main() {
	a := make(map[string]string)
	fmt.Printf("类型是: %T\n", a)
	fmt.Println(unsafe.Sizeof(a))
}
