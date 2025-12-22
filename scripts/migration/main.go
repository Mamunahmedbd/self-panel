package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/migrate"
)

func main() {
	dsn := "mim_it:P@$$mim_it!1@tcp(localhost:3308)/mim_it?parseTime=True"
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()

	// Dump migration changes to stdout
	if err := client.Schema.WriteTo(context.Background(), log.Writer(), migrate.WithGlobalUniqueID(true)); err != nil {
		log.Fatalf("failed printing schema changes: %v", err)
	}
	fmt.Println("\nSuccess")
}
