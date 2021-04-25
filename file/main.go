package main

import (
	"fmt"
	"os"
)

func main() {

	fileObj, err := os.Open("./main.go")

	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Printf("%T \n", fileObj)
	fileInfo, err := fileObj.Stat()

	if err != nil {
		fmt.Println(err)
		return
	}

	// 输出文件大小
	fmt.Printf("%d b \n", fileInfo.Size()) // 249 b
	fmt.Printf("%s \n", fileInfo.Name())   //

}
