package main

import "os"

func main() {
	root := defaultRoot(os.Args)
	cmd, _ := parseCmd(os.Args)

	switch cmd {
	case "load":
		exitIfErr(runLoad(root))
	case "validate":
		exitIfErr(runValidate(root))
	case "emit":
		exitIfErr(runEmit(root))
	default:
		usage()
		exitIfErr(errInvalidCmd)
	}
}
