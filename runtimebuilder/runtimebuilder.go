package runtimebuilder

import (
	"errors"
	"log"
	"strconv"
)

type VarType string

const (
	NoneVarType   VarType = "none"
	IntVarType            = "int"
	StringVarType         = "string"
)

type Value struct {
	Name        string
	Type        VarType
	IntValue    int
	StringValue string
}

type Procedure struct {
	Name           string
	ReturnType     VarType
	Variables      []*Value
	ParameterCount int
	Parameters     []*Value
	Procs          []*ProcInvocation
	ReturnVal      *Value
}

type ProcInvocation struct {
	Proc     Executable
	Params   []*Value
	BoundVar *Value
}

type Executable interface {
	Run([]*Value) *Value
	GetName() string
	GetReturnType() VarType
}

var procedures = []Executable{
	&SumProc{
		Name:       "suma",
		ReturnType: IntVarType,
	},
	&SubProc{
		Name:       "sub",
		ReturnType: IntVarType,
	},
	&MultProc{Name: "mult",
		ReturnType: IntVarType,
	},
	&PrintProc{
		Name:       "print",
		ReturnType: NoneVarType,
	},
}

func ProcedureExists(n string) (bool, Executable) {
	for k := range procedures {
		if procedures[k].GetName() == n {
			return true, procedures[k]
		}
	}
	return false, nil
}

func AddProcedure() {
	if len(procedures) == 0 {
		procedures = make([]Executable, 0)
	}
	procedures = append(procedures, &Procedure{})
}

func LastProcedure() *Procedure {
	if len(procedures) == 0 {
		return nil
	}
	return procedures[len(procedures)-1].(*Procedure)
}

func (p *Procedure) SetName(n string) {
	p.Name = n
}

func (p *Procedure) GetName() string {
	return p.Name
}

func (p *Procedure) SetReturnType(t VarType) {
	p.ReturnType = t
}

func (p *Procedure) GetReturnType() VarType {
	return p.ReturnType
}

func (p *Procedure) AddParameter() {
	if len(p.Parameters) == 0 {
		p.Parameters = make([]*Value, 0)
	}
	p.Parameters = append(p.Parameters, &Value{})
}

func (p *Procedure) LastParameter() *Value {
	if len(p.Parameters) == 0 {
		return nil
	}
	return p.Parameters[len(p.Parameters)-1]
}

func (p *Procedure) AddVariable() {
	if len(p.Variables) == 0 {
		p.Variables = make([]*Value, 0)
	}
	p.Variables = append(p.Variables, &Value{})
}

func (p *Procedure) LastVariable() *Value {
	if len(p.Variables) == 0 {
		return nil
	}
	return p.Variables[len(p.Variables)-1]
}

func (p *Procedure) VariableExists(n string) (bool, *Value) {
	for k := range p.Variables {
		if p.Variables[k].Name == n {
			return true, p.Variables[k]
		}
	}
	for k := range p.Parameters {
		if p.Parameters[k].Name == n {
			return true, p.Parameters[k]
		}
	}
	return false, nil
}

func (p *Procedure) AddInvocation() {
	if len(p.Procs) == 0 {
		p.Procs = make([]*ProcInvocation, 0)
	}
	p.Procs = append(p.Procs, &ProcInvocation{})
}

func (p *Procedure) LastInvocation() *ProcInvocation {
	if len(p.Procs) == 0 {
		return nil
	}
	return p.Procs[len(p.Procs)-1]
}

func (i *ProcInvocation) SetProcedure(e Executable) {
	i.Proc = e
}

func (i *ProcInvocation) AddParameter(v *Value) {
	if len(i.Params) == 0 {
		i.Params = make([]*Value, 0)
	}
	i.Params = append(i.Params, v)
}

func (i *ProcInvocation) SetBoundVariable(v *Value) {
	i.BoundVar = v
}

func (p *Procedure) Run(v []*Value) *Value {
	if len(v) != len(p.Parameters) {
		log.Fatalf("la cantidad de parámetros utilizados para invocar '%s' no coincide con la cantidad de parámetros en su definición\n\ndefinición: %d; invocación: %d", p.Name, len(p.Parameters), len(v))
	}
	for k := range v {
		if p.Parameters[k].Type == v[k].Type {
			if p.Parameters[k].Type == IntVarType {
				p.Parameters[k].IntValue = v[k].IntValue
			} else if p.Parameters[k].Type == StringVarType {
				p.Parameters[k].StringValue = v[k].StringValue
			}
		} else {
			log.Fatalf("el tipo especificado para el parámetro '%d' no coincide con el tipo del valor proporcionado",
				p.Parameters[k].Name)
		}
	}
	for k := range p.Procs {
		val := p.Procs[k].Proc.Run(p.Procs[k].Params) // remember to bind the return variable
		if p.Procs[k].BoundVar != nil {
			// do binding logic here
			p.Procs[k].BoundVar.IntValue = val.IntValue
			p.Procs[k].BoundVar.StringValue = val.StringValue
		}
	}
	return p.ReturnVal
}

func (v *Value) SetName(n string) {
	v.Name = n
}

func (v *Value) SetType(t VarType) {
	v.Type = t
}

func (v *Value) SetValue(val string) error {
	if v.Type == IntVarType {
		i, err := strconv.Atoi(val)
		if err != nil {
			return errors.New("valor no válido para variable de tipo 'int'")
		}
		v.IntValue = i
		return nil
	} else if v.Type == StringVarType {
		v.StringValue = val
		return nil
	} else {
		return errors.New("error de intérprete; variable sin tipo")
	}
	return nil
}

func PrintStack() {
	log.Print("\n\nSTACK:\n\n")
	for k := range procedures {
		log.Printf("\n\n%s:\n\n%v", procedures[k].GetName(), procedures[k])
	}
}

func Run() {
	ex, proc := ProcedureExists("main")
	if ex {
		proc.Run(nil)
		print("\n")
	} else {
		log.Fatalf("La función 'main' no fue hallada")
	}
}
