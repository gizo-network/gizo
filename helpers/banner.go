package helpers

import (
	"fmt"
	"os"

	"github.com/dimiro1/banner"
)

func Banner() {
	//! Banner
	in, err := os.Open("banner.txt")
	defer in.Close()
	if err != nil {
		return
	}
	fmt.Println()
	banner.Init(os.Stdout, true, true, in)
	fmt.Println()
}
