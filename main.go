package main

import (
	"flag"
	"fmt"
)

func main() {

	flag.Parse()
	path := flag.Arg(0)
	fmt.Println(path)
	if path == "" {
		//path = "GreekNation.txt"
		PanicString("must set file")
	}
	m := NewModule(path)
	m.Parse()
	fmt.Println(m.Display())
	m.Format(path)
}
