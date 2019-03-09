package main

import (
	"fmt"
	"os"

	"github.com/contd/links/app"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbpath := os.Getenv("SQLITE_PATH")
	if dbpath == "" {
		dbpath = "./data/saved.sqlite"
	}
	fmt.Printf("Using db: %s\n", dbpath)
	a := app.App{}
	a.Initialize(dbpath)
	a.Run(":5555")
}
