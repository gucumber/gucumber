package cucumber

import (
	"fmt"
	"os"
)

func RunMain() {
	runner, err := GlobalContext.RunDir(os.Args[1])
	if err != nil {
		panic(err)
	}

	if len(runner.Unmatched) > 0 {
		fmt.Println("Some steps were missing, you can add them by using the following step definition stubs: ")
		fmt.Println("")
		fmt.Print(runner.MissingMatcherStubs())
	}
}
