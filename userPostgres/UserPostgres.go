package userPostgres
// CREATE TABLE userdata (
// 	ID SERIAL PRIMARY KEY,
// 	Name VARCHAR(30),
// 	Lastname VARCHAR(30),
// 	Telephone VARCHAR(30),
// 	Description VARCHAR(30)
// )

import (
	"fmt"
	"database/sql"

	_ "github.com/lib/pq"
)

var (
	Host = ""
	port = 5432
	User = ""
	Password = ""
	Dbname = ""
)

type UserData struct {
	Name string
	Lastname string
	Telephone string
	Description string
}


func openConnection() (*sql.DB, error) {
	comm := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host,port,User, Password, Dbname)
	db, err := sql.Open("postgres", comm)
	if err != nil {
		return nil, fmt.Errorf("Open: %w", err)
	}
	return db, nil
}


func ListUsers() ([]UserData, error) {
	db, err := openConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT name, lastname, telephone, description FROM userdata"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	userdata := []UserData{}
	for rows.Next() {
		var name, lastname, telephone, description string
		err := rows.Scan(&name,&lastname,&telephone,&description)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		data := UserData{
			Name: name,
			Lastname: lastname,
			Telephone: telephone,
			Description: description,
		}
		userdata = append(userdata, data)
	}
	return userdata, nil
}

func exists(telephone string) (bool, error) {
	db, err := openConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()
	var exists bool
	// QueryRow используется, когда ожидается ОДНА строка (одно число)
	query := "SELECT EXISTS(SELECT 1 FROM userdata WHERE telephone=$1)"
	err = db.QueryRow(query, telephone).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil



}

func InsertUser(name, lastname, telephone, description string) (UserData, error) {
	userdata :=  UserData{}
	exists, err := exists(telephone)
	if err != nil {
		return userdata, err
	}
	db, err := openConnection()
	if err != nil {
		return userdata, err
	}
	defer db.Close()
	if exists {
		return userdata, fmt.Errorf("пользователь по номер %s уже добавлен",telephone)
	}

	query := "INSERT INTO userdata(name, lastname, telephone, description) VALUES ($1,$2,$3,$4)"
	_, err = db.Exec(query, name, lastname, telephone, description)
	if err != nil {
		return userdata, fmt.Errorf("exec: %w", err)
	}
	userdata = UserData{
		Name: name,
		Lastname: lastname,
		Telephone: telephone,
		Description: description,
	}
	return userdata, nil
}

func DeleteUsers(telephons ...string) error {
	if len(telephons) == 0 {
		return fmt.Errorf("параметр telephons не передан")
	}
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	for _, tel := range telephons {
		exists, err := exists(tel)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("номер не найден: %s", tel)
		}

		query := "DELETE FROM userdata WHERE telephone=$1"
		_, err = db.Exec(query, tel)
		if err != nil {
			return fmt.Errorf("exec: %w", err)
		}
	}
	return nil
}


func DeleteAllUsers() error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	query := "DELETE FROM userdata"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
