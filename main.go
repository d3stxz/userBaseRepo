package main

import (
	"fmt"
	us "UserRepo/userPostgres"
)


func main() {
	us.Host = "localhost"
	us.Dbname = "postgres"
	us.User = "postgres"
	us.Password = "1029384756"

	// err := us.DeleteUser()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	

	rows, err :=  us.ListUsers()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, row := range rows {
		fmt.Println(row)
	}
}