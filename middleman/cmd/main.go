package main

import (
	"fmt"

	"github.com/BartoszBurgiel/cloud/middleman"
)

func main() {
	middleman := middleman.NewMiddlemanFromEnv()
	fmt.Println(middleman.Start())
}
