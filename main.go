package main

import (
	"fmt"
	"log"

	"go-micrograd/value"
	"go-micrograd/visualize"
)

func main() {
	// Let us define the following expression tree:

	// a = b + c
	a := value.New(8, "a")
	b := value.New(3, "b")
	c := a.Add(b, "c")

	// d = c * f
	d := c.Mult(value.New(2, "f"), "d")

	// generate graph.htm which can be rendered in UI for
	// seeing the expression tree graph.
	if err := visualize.HTMLFile(d, "out/graph.html"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote out/graph.html")
}
