package helpers

import (
	"fmt"
	"os"

	"github.com/dimiro1/banner"
)

//Banner outputs gizo banner in terminal
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
