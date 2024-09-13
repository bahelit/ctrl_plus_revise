package main

import (
	"log"

	"github.com/bahelit/ctrl_plus_revise/pkg/clipboard"
)

func main() {
	err := clipboard.WriteAll("Hello Clipboard!")
	if err != nil {
		log.Println("clipboard write all error: ", err)
	}

	text, err := clipboard.ReadAll()
	if err != nil {
		log.Println("clipboard read all error: ", err)
		return
	}

	if text != "" {
		log.Println("text is: ", text)
		// Output: 日本語
	}
}
