package main

import (
	"fmt"
	"os"

	"github.com/BartoszBurgiel/cloud/reciever"
)

func main() {

	rec, err := reciever.NewRecieverFromConfig(
		os.Args[1],
		os.Args[2],
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	ok, err := rec.AskForPackage()
	fmt.Println(ok, err)
}
