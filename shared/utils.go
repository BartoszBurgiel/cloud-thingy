package shared

import (
	"fmt"
	"time"
)

// DisplayProgressBar simulates a simple progress bar to indicate downloading/uploading of the packages
func DisplayProgressBar(title string, finish chan bool) {
	fmt.Printf("%s...", title)
	for {
		//fmt.Println("A")
		select {
		case _, ok := <-finish:
			if ok {
				fmt.Println("\nFinished!")
				return
			}
		default:
			fmt.Print(".")
			time.Sleep(time.Second)
		}
	}
}
