package main

import "os"

func main() {
	dbpath := os.Getenv("SQLITE_PATH")
	if dbpath == "" {
		dbpath = "saved.sqlite"
	}
	a := App{}
	a.Initialize(dbpath)
	a.Run(":5555")
}
