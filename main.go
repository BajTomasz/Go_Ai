package main

import (
	"Go_Ai/Tasks"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Enter the task number")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "0":
		Tasks.Helloapi()
	case "1":
		Tasks.Moderation()
	case "2":
		Tasks.Blogger()
	case "3":
		Tasks.Liar()
	case "4":
		Tasks.Inprompt()
	case "5":
		Tasks.Embedding()
	case "6":
		Tasks.Whisper()
	case "7":
		Tasks.Functions()
	case "8":
		Tasks.Rodo()
	case "9":
		Tasks.Scraper()
	case "10":
		Tasks.Whoami()
	case "11":
		Tasks.Search()
	case "12":
		Tasks.People()
	case "13":
		Tasks.Knowledge()
	case "14":
		Tasks.Tools()
	case "15":
		Tasks.Gnome()
	case "16":
		Tasks.Ownapi()
	case "17":
		Tasks.Ownapipro()
	default:
		fmt.Println("Enter the task number")
		os.Exit(1)
	}
}
