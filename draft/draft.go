package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print(err.Error())
		} else {
			fmt.Print(text)
		}
	}
}