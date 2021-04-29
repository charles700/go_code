reflect 反射
---- 
变量的信息
- 类型信息：预先定义好的元信息。
- 值信息：程序运行过程中可动态变化的。


## reflect.TypeOf()
可以获得任意值的类型信息
> t := reflect.TypeOf(a)
> t.Kind() -- 种类
> t.Name() -- 具体的类型


```golang
type Cat struct {
}
var a = Cat{}
t = reflect.TypeOf(a) // string
t.Kind() //  struct 
t.Name() //  Cat
```
## reflect.ValueOf()
获取值类型
> reflect.ValueOf()返回的是reflect.Value类型，其中包含了原始值的值信息。
> reflect.Value与原始值之间可以互相转换。

