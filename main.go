package main

import (
	"go-wc/mapReduceBigFile"
	"go-wc/mapReduceDir"
	"go-wc/single"
	"os"
)

func main() {
	var mode string
	if len(os.Args) > 1 {
		mode = os.Args[1]
	} else {
		mode = "single"
	}

	switch mode {
	case "mrbf":
		mapReduceBigFile.PrintWC()
		break
	case "mrd":
		mapReduceDir.PrintWC("/home/vit/Books/TXT_Books")
		break
	case "single":
	default:
		single.PrintWC()
	}
}
