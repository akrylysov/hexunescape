package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/akrylysov/hexunescape"
)

func main() {
	flag.Usage = func() {
		fmt.Println("usage: hexunescape [path]\n\npath is optional, defaults to stdin.")
		flag.PrintDefaults()
	}
	flag.Parse()

	in := os.Stdin
	if path := flag.Arg(0); path != "" {
		var err error
		if in, err = os.Open(path); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer in.Close()
	}
	if err := hexunescape.Unescape(os.Stdout, in); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
