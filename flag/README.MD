## flag

## flag.[Type] -- 基本使用

- 返回 指针类型
- 运行 `go run flag/main.go --name zhangsan --age=20`

```golang
	// (字段名, 默认值, 帮助信息)
	name := flag.String("name", "lisi", "填写姓名")
	age := flag.Int("age", 18, "填写年龄")
	married := flag.Bool("married", false, "结婚了吗？")

	flag.Parse()

	fmt.Println(*name)
	fmt.Println(*age)
	fmt.Println(*married)
```

## flag.[type]Var

- 先定义变量，在赋值

```golang
	var grade int

	flag.IntVar(&grade, "grade", 7, "班级信息")

	flag.Parse()

    fmt.Println(grade)
```

## flag.Parse()

- 在所有 flag 都注册之后，调用：flag.Parse() 来解析命令行参数写入注册的 flag 里。
  支持的命令行参数格式有以下几种：

1. -flag xxx （使用空格，一个-符号）
2. --flag xxx （使用空格，两个-符号）
3. -flag=xxx （使用等号，一个-符号）
4. --flag=xxx （使用等号，两个-符号）

## --help

> go run flag/main.go --help

## 其他

```golang
	fmt.Println(flag.Args())    //返回命令行参数后的其他参数
	fmt.Println(flag.NArg()) 	//返回命令行参数后的其他参数个数
	fmt.Println(flag.NFlag()) 	//返回使用的命令行参数个数
```
