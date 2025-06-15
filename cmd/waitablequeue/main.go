package main

import (
	"os"
)

func main() {
	os.Exit(realMain(os.Stdout, 60, 20, 30))
}
