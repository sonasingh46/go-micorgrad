package main

import (
	"fmt"
	"log"

	"go-micrograd/value"
	"go-micrograd/visualize"
)

func main() {
	// Let us define 'a' and 'b'
	// add 'a' and 'b' and store the result in 'c'

	a := value.New(8, "a")
	b := value.New(3, "b")
	c := a.Add(b, "c")

	// generate graph.htm which can be rendered in UI for
	// seeing the expression tree graph.
	if err := visualize.HTMLFile(c, "out/graph.html"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote out/graph.html")
}
