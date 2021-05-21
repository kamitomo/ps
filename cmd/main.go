package main

import (
	"flag"
	"log"
	"os"

	"github.com/tomohiro8/ps"
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	vm := ps.NewVM()
	err = vm.Execute(f)
	if err != nil {
		log.Fatal(err)
	}
}
