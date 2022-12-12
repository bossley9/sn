package simperium

import "fmt"

const magenta = "\033[0;35m"
const none = "\033[0m"

func printDebugMessage(message string) {
	fmt.Println(magenta)
	fmt.Print(message)
	fmt.Println(none)
}
