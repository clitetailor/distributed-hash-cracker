package main

import (
	"fmt"
	"flag"
)

func main() {
	someFlag := flag.String("my-flag", "Hello World!", "this is my flag")
	flag.Parse()
	fmt.Println(*someFlag)
}
