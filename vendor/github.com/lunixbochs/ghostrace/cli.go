package main

import (
	"fmt"
	"os"

	"github.com/lunixbochs/ghostrace/ghost"
)

func main() {
	trace, err := ghost.DefaultCli(os.Args, nil)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	// TODO: no way to kill target, so this should be more than a channel
	for sc := range trace {
		fmt.Fprintf(os.Stderr, "%+v\n", sc)
	}
}
