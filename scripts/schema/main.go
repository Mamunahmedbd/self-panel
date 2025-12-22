package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide table name")
	}
	tableName := os.Args[1]

	dsn := "mim_it:P@$$mim_it!1@tcp(localhost:3308)/mim_it?parseTime=True"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Printf("--- Schema for %s ---\n", tableName)
	rows, err := db.Query(fmt.Sprintf("DESCRIBE %s", tableName))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var field, typ, null, key, def, extra sql.NullString
		if err := rows.Scan(&field, &typ, &null, &key, &def, &extra); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Field: %-20s Type: %-20s Null: %-5s Key: %-5s Default: %-15s Extra: %s\n",
			field.String, typ.String, null.String, key.String, def.String, extra.String)
	}
}
