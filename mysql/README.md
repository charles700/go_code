处理 `sql.ErrNoRows` 错误 方案
------------------------------------------------


1. `sql.ErrNoRows` 错误并应该再次包装返回给上层。
2. 原因：虽然`sql.ErrNoRows` 算作底层的错误，但是对使用dao层的调用者而言，调用端给前端返回的要么是空对象，要么是空数组，是一种合理的业务逻辑。
3. 处理：
- 直接返回错误 或 返回空值

```golang
func queryUserById(id int) (user, error) {

	sqlStr := "select id, name, age from users where id=?"

	var u user
	err := db.QueryRow(sqlStr, id).Scan(&u.id, &u.name, &u.age)

	if err == sql.ErrNoRows {
		// fmt.Println("没找到记录")
		return u, nil
	}

	if err != nil {
		// fmt.Printf("%T, err:%v\n", err, err)
		return u, errors.WithMessage(err, fmt.Sprintf("test.users : 查询 id=%v 的记录出错 \n", id))
	}

	return u, nil
}
```
