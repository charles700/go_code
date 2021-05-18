wire  -- 依赖注入
---
[Go 每日一库之 wire](https://zhuanlan.zhihu.com/p/110453784)

## 安装
> go get github.com/google/wire/cmd/wire

## 简单使用
1. 项目文件中定义好 相关的结构体 和 New 函数
2. 新建 `wire.go` 文件, 编写代码
```golang
//+build wireinject
package main

import "github.com/google/wire"

func InitMission(name string) Mission {
  wire.Build(NewMonster, NewPlayer, NewMission)
  return Mission{}
}
```
3. 终端之命令, 会生成 wire_gen.go 文件
> wire

4. 在 项目 main 函数中直接使用 wire_gen.go 文件中的方法
```golang
func main() {
  mission := InitMission("dj")

  mission.Start()
}
```


## 传参
- 为不同的结构体 New 函数定义 不同的参数类型
- 在 wire.go 文件中，InitXXX 方法传入对应类型即可。


## 错误处理
参考[文章](https://zhuanlan.zhihu.com/p/110453784)内容

