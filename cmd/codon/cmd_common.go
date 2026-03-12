package main

import (
	"errors"
	"fmt"
	"os"
)

var errInvalidCmd = errors.New("invalid command")

func defaultRoot(args []string) string {
	if len(args) > 2 {
		return args[2]
	}
	return ".codon"
}

func parseCmd(args []string) (string, []string) {
	if len(args) < 2 {
		return "", nil
	}
	return args[1], args[2:]
}

func exitIfErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Usage: codon <load|validate|emit> [path]")
}
