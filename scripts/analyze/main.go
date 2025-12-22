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

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			log.Fatal(err)
		}
		tables = append(tables, table)
	}

	for _, table := range tables {
		fmt.Printf("--- TABLE: %s ---\n", table)
		columnRows, err := db.Query(fmt.Sprintf("DESCRIBE %s", table))
		if err != nil {
			fmt.Printf("Error describing table %s: %v\n", table, err)
			continue
		}

		for columnRows.Next() {
			var field, typ, null, key, def, extra sql.NullString
			if err := columnRows.Scan(&field, &typ, &null, &key, &def, &extra); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Field: %s, Type: %s, Null: %s, Key: %s, Default: %s, Extra: %s\n",
				field.String, typ.String, null.String, key.String, def.String, extra.String)
		}
		columnRows.Close()
		fmt.Println()
	}
}
