package main

import (
	"fmt"
	"os"
)

func main() {
	dbpath := os.Getenv("SQLITE_PATH")
	if dbpath == "" {
		dbpath = "saved.sqlite"
	}
	fmt.Printf("Using db: %s\n", dbpath)
	a := App{}
	a.Initialize(dbpath)
	a.Run(":5555")
}
