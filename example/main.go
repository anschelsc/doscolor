package main

import (
	"fmt"
	"os"
	"github.com/anschelsc/doscolor"
)

func main() {
	output := doscolor.NewWrapper(os.Stdout)
	if err := output.Save(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintln(output, "Here's some ordinary text.")
	if err := output.SetMask(doscolor.Red | doscolor.Bright, doscolor.Foreground); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintln(output, "Here's some (revolutionary) red text.")
	if err := output.Set(doscolor.White | doscolor.Bright | doscolor.BG(doscolor.Green)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintln(output, "Esperanto colors.")
	if err := output.Restore(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintln(output, "Now it should be back to normal.")
}
