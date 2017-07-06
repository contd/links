package main

func main() {
	a := App{}
	a.Initialize("apiuser", "apipass", "links")
	a.Run(":5555")
}
