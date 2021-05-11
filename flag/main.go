package main

import (
	"flag"
	"fmt"
)

func main() {

	// (字段名, 默认值, 提示 -- help 的提示)
	name := flag.String("name", "lisi", "填写姓名")
	age := flag.Int("age", 18, "填写年龄")
	married := flag.Bool("married", false, "结婚了吗？")

	var grade int

	flag.IntVar(&grade, "grade", 7, "班级信息")

	flag.Parse()

	fmt.Println(*name)
	fmt.Println(*age)
	fmt.Println(*married)
	fmt.Println(grade)
}
