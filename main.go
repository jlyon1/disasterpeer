package main

func main() {
	a := NewServer("0.0.0.0:8080")
	a.Serve()
}
