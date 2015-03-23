package gucumber

import "fmt"

func RunMain(args ...string) {
	if len(args) == 0 {
		args = []string{"features"}
	}

	runner, err := GlobalContext.RunDir(args[0])
	if err != nil {
		panic(err)
	}

	if len(runner.Unmatched) > 0 {
		fmt.Println("Some steps were missing, you can add them by using the following step definition stubs: ")
		fmt.Println("")
		fmt.Print(runner.MissingMatcherStubs())
	}
}
