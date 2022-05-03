package main

import (
	"fmt"
	"os"

	"github.com/BartoszBurgiel/cloud/client"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments provided.")
		return
	}
	c, err := client.NewClientFromConfigFile(
		os.Args[1],
		os.Args[2:],
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = c.Sumbit(); err != nil {
		fmt.Println(err)
	}
}
