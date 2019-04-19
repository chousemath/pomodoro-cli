package clear

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var commands = map[string]func(){
	"darwin": func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	},
	"linux": func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	},
	"windows": func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	},
}

// Execute performs the clear screen command
func Execute() {
	value, ok := commands[runtime.GOOS] // runtime.GOOS -> linux, windows, darwin etc.
	if ok {                             // if we defined a clear func for that platform:
		value() // we execute it
	} else { //unsupported platform
		fmt.Printf("Your OS (%s) is not supported", runtime.GOOS)
	}
}
