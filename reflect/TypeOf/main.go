package main

import (
	"fmt"
	"reflect"
)

type Cat struct {
}

func main() {
	var a = Cat{}
	fmt.Printf("type: %v %v\n", reflect.TypeOf(a).Kind(), reflect.TypeOf(a).Name()) // type: struct {}

}
