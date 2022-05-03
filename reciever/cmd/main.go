package main

import (
	"fmt"
	"os"

	"github.com/BartoszBurgiel/cloud/reciever"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Not enough arguments provided.")
		return
	}

	rec, err := reciever.NewRecieverFromConfig(
		os.Args[1],
		os.Args[2],
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	rec.AskForPackage()
}
