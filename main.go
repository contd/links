package main

func main() {
	a := App{}
	a.Initialize("/data/saved.sqlite")
	a.Run(":5555")
}
