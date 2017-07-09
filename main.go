package main

func main() {
	a := App{}
	a.Initialize("saved.sqlite")
	a.Run(":5555")
}
