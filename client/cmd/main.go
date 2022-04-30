package main

import (
	"fmt"
	"os"

	"github.com/BartoszBurgiel/cloud/client"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Not enough arguments provided.")
		return
	}
	c, err := client.NewClientFromConfigFile(
		os.Args[1],
		os.Args[2],
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(c.Sumbit())
}
