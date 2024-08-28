/*
Main package with main func printing reverse string
*/
package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	fmt.Println(reverse.String("Hello, OTUS!"))
}
