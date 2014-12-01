package runtimebuilder

import (
	"log"
	"fmt"
	"strconv"
)

type SumProc Procedure

type SubProc Procedure

type MultProc Procedure

type PrintProc Procedure

func (s *SumProc) Run(v []*Value) *Value {
	if len(v) != 2 {
		log.Fatal("número de parámetros para 'suma' incorrectos")
	}
	for k := range v {
		if v[k] == nil {
			log.Fatal("error fatal: valor sin inicializar en 'suma'")
		}
		if v[k].Type != IntVarType {
			log.Fatal("'suma' requiere parámetros de tipo 'int'")
		}
	}
	return &Value{
		Type:     IntVarType,
		IntValue: v[0].IntValue + v[1].IntValue,
	}
}

func (s *SumProc) GetName() string {
	return s.Name
}

func (s *SumProc) GetReturnType() VarType {
	return s.ReturnType
}

func (s *SubProc) Run(v []*Value) *Value {
	if len(v) != 2 {
		log.Fatal("número de parámetros para 'resta' incorrectos")
	}
	for k := range v {
		if v[k] == nil {
			log.Fatal("error fatal: valor sin inicializar en 'resta'")
		}
		if v[k].Type != IntVarType {
			log.Fatal("'resta' requiere parámetros de tipo 'int'")
		}
	}
	return &Value{
		Type:     IntVarType,
		IntValue: v[0].IntValue - v[1].IntValue,
	}
}

func (s *SubProc) GetName() string {
	return s.Name
}

func (s *SubProc) GetReturnType() VarType {
	return s.ReturnType
}

func (m *MultProc) Run(v []*Value) *Value {
	if len(v) != 2 {
		log.Fatal("número de parámetros para 'mult' incorrectos")
	}
	for k := range v {
		if v[k] == nil {
			log.Fatal("error fatal: valor sin inicializar en 'mult'")
		}
		if v[k].Type != IntVarType {
			log.Fatal("'mult' requiere parámetros de tipo 'int'")
		}
	}
	return &Value{
		Type:     IntVarType,
		IntValue: v[0].IntValue * v[1].IntValue,
	}
}

func (m *MultProc) GetName() string {
	return m.Name
}

func (m *MultProc) GetReturnType() VarType {
	return m.ReturnType
}

func (p *PrintProc) Run(v []*Value) *Value {
	if len(v) != 1 {
		log.Fatal("número de parámetros para 'mult' incorrectos")
	}
	if v[0] == nil {
		log.Fatal("error fatal: valor sin inicializar en 'print'")
	}
	if v[0].Type == StringVarType {
		//log.Fatal("'print' requiere parámetros de tipo 'string'") <-- old logic
		s, _ := strconv.Unquote("\"" + v[0].StringValue + "\"")
		fmt.Print(s)
	} else if v[0].Type == IntVarType {
		fmt.Printf("%d", v[0].IntValue)
	}
	return nil
}

func (p *PrintProc) GetName() string {
	return p.Name
}

func (p *PrintProc) GetReturnType() VarType {
	return p.ReturnType
}
