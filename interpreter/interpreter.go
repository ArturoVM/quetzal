package interpreter

import (
	"../lexer"
	"../parser"
	"../runtimebuilder"
	"io"
	"log"
)

func Interpret(path string) {
	err := parser.Read(path)
	if err != nil {
		log.Fatalf("error opening file: %s", err.Error())
	}
	t, err := parser.NextToken()
	if err != nil {
		log.Fatalf("error parsing: %s", err.Error())
	}
	for lexer.Lex(t) {
		t, err = parser.NextToken()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("error fatal: caracter no válido encontrado: %s", err.Error())
		}
	}
	//runtimebuilder.PrintStack()
	runtimebuilder.Run()
}
