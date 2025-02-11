package models

import (
	"database/sql"
	"strconv"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func ConnectDatabase() error {
	db, err := sql.Open("sqlite", "./users.db")
	if err != nil {
		return err
	}

	DB = db
	return nil
}

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func GetUsers(count int) ([]User, error) {

	rows, err := DB.Query("SELECT id, email from users LIMIT " + strconv.Itoa(count))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		singleUser := User{}
		err = rows.Scan(&singleUser.Id, &singleUser.Email)

		if err != nil {
			return nil, err
		}

		users = append(users, singleUser)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return users, err
}

func GetUserById(id string) (User, error) {

	stmt, err := DB.Prepare("SELECT id, email from users WHERE id = ?")

	if err != nil {
		return User{}, err
	}

	user := User{}

	sqlErr := stmt.QueryRow(id).Scan(&user.Id, &user.Email)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, sqlErr
	}
	return user, nil
}

func AddUser(newUser User) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO users (email) VALUES (?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newUser.Email)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func UpdateUser(user User, id int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE users SET email = ? WHERE Id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.Email, user.Id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeleteUser(userId int) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE from users where id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(userId)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
