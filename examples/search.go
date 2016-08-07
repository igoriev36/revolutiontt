package main

import (
	"fmt"
	"github.com/chrisabruce/revolutiontt"
	"os"
)

func main() {
	r := new(revolutiontt.RevolutionTT)
	r.Connect(os.Getenv("rev_username"), os.Getenv("rev_password"))
	results, _ := r.Search("Caddyshack")

	for _, v := range results {
		fmt.Println(v.Title + ": " + v.DownloadUrl)
	}
}
