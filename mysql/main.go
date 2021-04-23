package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type user struct {
	id   int
	age  int
	name string
}

var db *sql.DB

func initDB() (err error) {
	dsn := "root:123456789a@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return nil
}

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

func main() {

	initDB()

	// u, err := queryUserById(1)
	u, err := queryUserById(100)

	if err != nil {
		fmt.Printf("%v \n", err)
		return
	}

	fmt.Println(u)

	defer db.Close()
}
