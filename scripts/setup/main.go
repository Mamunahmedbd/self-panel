package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "mim_it:P@$$mim_it!1@tcp(localhost:3308)/mim_it?parseTime=True"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS tickets (
		id bigint(20) NOT NULL AUTO_INCREMENT,
		subject varchar(255) NOT NULL,
		description text NOT NULL,
		status enum('open','pending','closed') NOT NULL DEFAULT 'open',
		priority enum('low','medium','high') NOT NULL DEFAULT 'medium',
		client_id int(11) NOT NULL,
		client_username varchar(255) NOT NULL,
		created_at datetime NOT NULL,
		updated_at datetime NOT NULL,
		PRIMARY KEY (id),
		KEY client_id (client_id),
		KEY client_username (client_username),
		KEY status (status),
		KEY created_at (created_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tickets table created successfully or already exists.")
}
