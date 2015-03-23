package gucumber

func RunMain(args ...string) {
	if len(args) == 0 {
		args = []string{"features"}
	}

	if err := BuildAndRunDir(args[0]); err != nil {
		panic(err)
	}
}
