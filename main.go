package main

import (
	"github.com/tristanwietsma/tabler/lib"
	"log"
	"os"
)

func main() {
	for _, path := range os.Args[1:] {
		infile := lib.InputFile{}
		if err := infile.Init(path); err != nil {
			log.Printf("%v", err)
			continue
		}
		infile.Write()
	}
}
