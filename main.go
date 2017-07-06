package main

func main() {
	a := App{}
	a.Initialize("apiuser", "wj5np47dn", "links")
	a.Run(":5555")
}
