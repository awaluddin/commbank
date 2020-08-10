package main

import (
	"flag"
	"fmt"
	"strconv"
)

func main() {

	// get number from -num
	numStr := flag.String("num", "4", "input number")
	flag.Parse()

	// convert string number to integer
	n, err := strconv.Atoi(*numStr)
	if err != nil {
		// failed convert string to int, then set 0
		fmt.Println("please make sure you enter a number type")
		n = 0
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			// write white space
			fmt.Print(" ")
		}
		for j := 0; j <= i; j++ {
			// write hastag
			fmt.Print("#")
		}
		// new line
		fmt.Println()
	}
}
