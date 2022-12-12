package debug

import "fmt"

func PrintPlain(a ...any) {
	fmt.Print(NoColor)
	fmt.Print(a...)
}

func PrintInfo(a ...any) {
	fmt.Print(Cyan)
	fmt.Print(a...)
	fmt.Print(NoColor)
}

func PrintWarning(a ...any) {
	fmt.Print(Yellow)
	fmt.Print(a...)
	fmt.Print(NoColor)
}

func PrintError(a ...any) {
	fmt.Print(Red)
	fmt.Print(a...)
	fmt.Print(NoColor)
}

func PrintDebug(a ...any) {
	fmt.Print(Magenta)
	fmt.Print(a...)
	fmt.Print(NoColor)
}
