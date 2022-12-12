package simperium

import "fmt"

const yellow = "\033[0;33m"
const none = "\033[0m"

func printDebugMessage(message string) {
	fmt.Println(yellow)
	fmt.Print(message)
	fmt.Println(none)
}
