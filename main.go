package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Dilshod@2005"
	dbname   = "demo"
)

func main() {

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to the database!")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id serial primary key,
			username varchar(50) not null,
			email varchar(100) unique not null,
			password varchar(100) NOT NULL,
			created_at timestamp default current_timestamp
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Users table created successfully!")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS friendships (
			user_id INT NOT NULL,
			friend_id INT NOT NULL,
			status VARCHAR(20) NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(user_id),
			FOREIGN KEY (friend_id) REFERENCES users(user_id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("friendships table creatttted")

}

func sendRequest(db *sql.DB, fromUserID, toUserID int) error {
	_, err := db.Exec("INSERT INTO friendships (user_id, friend_id, status) VALUES ($1, $2, 'send')", fromUserID, toUserID)
	return err
}

func acceptRequest(db *sql.DB, fromUserID, toUserID int) error {
	_, err := db.Exec("UPDATE friendships SET status = 'accepted' WHERE user_id = $1 AND friend_id = $2", toUserID, fromUserID)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO friendships (user_id, friend_id, status) VALUES ($1, $2, 'accepted')", toUserID, fromUserID)
	return err
}

func block(db *sql.DB, userID, blockedUserID int) error {
	_, err := db.Exec("UPDATE friendships SET status = 'blocked' WHERE user_id = $1 AND friend_id = $2", userID, blockedUserID)
	return err
}

func unblock(db *sql.DB, userID, unblockedUserID int) error {
	_, err := db.Exec("DELETE FROM friendships WHERE user_id = $1 AND friend_id = $2", userID, unblockedUserID)
	return err
}
