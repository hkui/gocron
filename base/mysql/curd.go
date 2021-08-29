package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var (
		err error
		db  *sql.DB
	)
	//user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
	db, err = sql.Open("mysql", "root:mysql2019@tcp(39.100.78.46:3306)/cron?utf8")
	if err != nil {
		fmt.Println(err)
		return
	}
	//插入单条数据
	stmt, err := db.Prepare("INSERT INTO user  SET uname=?,pwd=?")
	checkErr(err)
	res, err := stmt.Exec("hkui", "123456")
	checkErr(err)
	fmt.Println(res.LastInsertId())

}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
