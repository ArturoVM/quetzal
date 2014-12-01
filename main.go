package main

import (
	"./interpreter"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		println("\n\tArgumentos insuficientes o en demasía.")
		println("\n\tUso: quetzal [path]\n")
		os.Exit(1)
	}
	interpreter.Interpret(os.Args[1])
}
