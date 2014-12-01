package lexer

import (
	"../parser"
	"../runtimebuilder"
	"log"
	"strconv"
	"strings"
)

type LexerState int

const (
	None LexerState = iota
	ProcedureDef
	ProcedureDefName
	ProcedureDefParamsStart
	ProcedureDefParamName
	ProcedureDefParamType
	ProcedureDefParamNext
	ProcedureDefParamsEnd
	ProcedureDefReturnType
	ProcedureDefBodyStart
	ProcedureDefReturn
	ProcedureDefReturnValue
	ProcedureDefBodyEnd
	ProcedureCall
	ProcedureCallParamsStart
	ProcedureCallParamValue
	ProcedureCallParamNext
	ProcedureCallParamsEnd
	VariableDef
	VariableDefName
	VariableDefType
	VariableDefAssign
	VariableDefValue
	VariableUse
	StringDefStart
	StringDefValue
	StringDefEnd
	InstructionEnd
)

type Keyword string

const (
	ProcDef    Keyword = "proc"
	ProcMain           = "main"
	VarDef             = "let"
	Return             = "return"
	PrintProc          = "print"
	SumProc            = "suma"
	RestaProc          = "resta"
	DivideProc         = "divide"
	IntType            = "int"
	StringType         = "string"
)

type Symbol string

const (
	ParensLeft       Symbol = "("
	ParensRight             = ")"
	CurlyLeft               = "{"
	CurlyRight              = "}"
	Assign                  = "="
	Comma                   = ","
	InstructionBreak        = ";"
	Quote                   = "\""
	EOL                     = "\n"
)

var (
	currentState LexerState
	parentState  LexerState = None
	stringBuffer []string
)

func SetNextState(s LexerState) {
	currentState = s
}

func SetParentState(s LexerState) {
	parentState = s
}

func GetState() LexerState {
	return currentState
}

func GetParentState() LexerState {
	return parentState
}

func Lex(t *parser.Token) bool {
	switch t.Content {
	case string(ProcDef):
		SetNextState(ProcedureDef)
		runtimebuilder.AddProcedure()
	case string(ParensLeft):
		if GetState() == ProcedureDefName {
			SetNextState(ProcedureDefParamsStart)
		} else if GetState() == ProcedureCall {
			SetNextState(ProcedureCallParamsStart)
		} else if GetState() != StringDefStart && GetState() != StringDefValue {
			log.Fatalf("Línea %d: Caracter '(' no esperado", parser.LineNumber)
		}
	case string(IntType):
		if GetState() == ProcedureDefParamName {
			SetNextState(ProcedureDefParamType)
			runtimebuilder.LastProcedure().LastParameter().SetType(runtimebuilder.IntVarType)
		} else if GetState() == VariableDefName {
			SetNextState(VariableDefType)
			runtimebuilder.LastProcedure().LastVariable().SetType(runtimebuilder.IntVarType)
		} else if GetState() == ProcedureDefParamsEnd {
			SetNextState(ProcedureDefReturnType)
			runtimebuilder.LastProcedure().SetReturnType(runtimebuilder.IntVarType)
		} else {
			log.Fatalf("Línea %d: Definición de tipo 'int' no esperada", parser.LineNumber)
		}
	case string(StringType):
		if GetState() == ProcedureDefParamName {
			SetNextState(ProcedureDefParamType)
			runtimebuilder.LastProcedure().LastParameter().SetType(runtimebuilder.StringVarType)
		} else if GetState() == VariableDefName {
			SetNextState(VariableDefType)
			runtimebuilder.LastProcedure().LastVariable().SetType(runtimebuilder.StringVarType)
		} else if GetState() == ProcedureDefParamsEnd {
			SetNextState(ProcedureDefReturnType)
			runtimebuilder.LastProcedure().SetReturnType(runtimebuilder.StringVarType)
		} else {
			log.Fatalf("Línea %d: Definición de tipo 'string' no esperada", parser.LineNumber)
		}
	case string(Comma):
		if GetState() == ProcedureDefParamType {
			SetNextState(ProcedureDefParamNext)
		} else if GetState() == ProcedureCallParamValue {
			SetNextState(ProcedureCallParamNext)
		} else {
			log.Fatalf("Línea %d: Caracter ',' no esperado", parser.LineNumber)
		}
	case string(ParensRight):
		if GetState() == ProcedureDefParamsStart || GetState() == ProcedureDefParamType {
			SetNextState(ProcedureDefParamsEnd)
		} else if GetState() == ProcedureCallParamsStart || GetState() == ProcedureCallParamValue {
			if GetParentState() != None {
				SetNextState(GetParentState())
				SetParentState(None)
			} else {
				SetNextState(ProcedureCallParamsEnd)
			}
		} else if GetState() != StringDefStart && GetState() != StringDefValue {
			log.Fatalf("Línea %d: Caracter ')' no esperado", parser.LineNumber)
		}
	case string(CurlyLeft):
		if GetState() == ProcedureDefParamsEnd || GetState() == ProcedureDefReturnType {
			if GetState() != ProcedureDefReturnType {
				runtimebuilder.LastProcedure().SetReturnType(runtimebuilder.NoneVarType)
			}
			SetNextState(ProcedureDefBodyStart)
		} else if GetState() != StringDefStart && GetState() != StringDefValue {
			log.Fatalf("Línea %d: Caracter '{' no esperado", parser.LineNumber)
		}
	case string(CurlyRight):
		if GetState() == InstructionEnd {
			if runtimebuilder.LastProcedure().ReturnVal == nil && runtimebuilder.LastProcedure().GetReturnType() != runtimebuilder.NoneVarType {
				log.Fatalf("Línea %d: El procedimiento '%s' especifica un tipo, pero la instrucción de 'return' no fue hallada",
					parser.LineNumber, runtimebuilder.LastProcedure().GetName())
			}
			SetNextState(ProcedureDefBodyEnd)
		} else if GetState() != StringDefStart && GetState() != StringDefValue {
			log.Fatalf("Línea %d: Caracter '}' no esperado", parser.LineNumber)
		}
	case string(VarDef):
		if GetState() != ProcedureDefBodyStart && GetState() != InstructionEnd {
			log.Fatalf("Línea %d: Definición de variable no esperada", parser.LineNumber)
		} else {
			SetNextState(VariableDef)
			runtimebuilder.LastProcedure().AddVariable()
		}
	case string(Assign):
		if GetState() == VariableDefType {
			SetNextState(VariableDefAssign)
		} else if GetState() == VariableDefName {
			log.Fatalf("Línea %d: Se esperaba declaración de tipo, pero se encontró '='", parser.LineNumber)
		} else {
			log.Fatalf("Línea %d: Caracter '=' no esperado", parser.LineNumber)
		}
	case string(InstructionBreak):
		if GetState() == ProcedureCallParamsEnd || GetState() == VariableDefValue || GetState() == ProcedureDefReturnValue {
			SetNextState(InstructionEnd)
		} else if GetState() != StringDefValue {
			log.Fatalf("Línea %d: Caracter ';' no esperado", parser.LineNumber)
		}
	case string(Return):
		if GetState() != InstructionEnd && GetState() != ProcedureDefBodyStart {
			log.Fatalf("Línea %d: Instrucción 'return' no esperada", parser.LineNumber)
		} else {
			if runtimebuilder.LastProcedure().GetReturnType() == runtimebuilder.NoneVarType {
				log.Fatalf("Línea %d: 'return' sobrante; la función %s no especifica un tipo", parser.LineNumber, runtimebuilder.LastProcedure().GetName())
			}
			SetNextState(ProcedureDefReturn)
		}
	case string(Quote):
		s := GetState()
		if s == StringDefValue {
			if GetParentState() != None {
				SetNextState(GetParentState())
				SetParentState(None)
			} else {
				SetNextState(StringDefEnd)
			}
			if GetState() == VariableDefValue {
				if runtimebuilder.LastProcedure().LastVariable().Type == runtimebuilder.StringVarType {
					runtimebuilder.LastProcedure().LastVariable().StringValue = strings.Join(stringBuffer, " ")
					stringBuffer = make([]string, 0)
				}
			} else if GetState() == ProcedureCallParamValue {
				if len(stringBuffer) > 0 {
					runtimebuilder.LastProcedure().LastInvocation().AddParameter(&runtimebuilder.Value{Type: runtimebuilder.StringVarType, StringValue: strings.Join(stringBuffer, " ")})
					stringBuffer = make([]string, 0)
				}
			}
		}
		if (s == ProcedureCallParamsStart || s == ProcedureCallParamNext) || s == VariableDefAssign {
			if s == ProcedureCallParamsStart || s == ProcedureCallParamNext {
				SetParentState(ProcedureCallParamValue)
			} else if s == VariableDefAssign {
				SetParentState(VariableDefValue)
			}
			SetNextState(StringDefStart)
			stringBuffer = make([]string, 0)
		}
	case string(EOL):
		if GetState() == StringDefValue {
			log.Fatalf("Línea %d: Los strings multi–línea no están permitidos", parser.LineNumber)
		}
	default:
		switch s := GetState(); {
		case s == ProcedureDef:
			SetNextState(ProcedureDefName)
			runtimebuilder.LastProcedure().SetName(t.Content)
		case s == VariableDef:
			SetNextState(VariableDefName)
			runtimebuilder.LastProcedure().LastVariable().SetName(t.Content)
		case s == ProcedureDefReturn:
			// check if var exists in runtimebuilder or if it is a value
			ex, v := runtimebuilder.LastProcedure().VariableExists(t.Content)
			if ex {
				// bind variable to proc return value
				SetNextState(ProcedureDefReturnValue)
				runtimebuilder.LastProcedure().ReturnVal = v
			} else {
				log.Fatalf("Línea %d: La variable '%s' no existe", parser.LineNumber, t.Content)
			}
		case s == ProcedureDefParamsStart || s == ProcedureDefParamNext:
			SetNextState(ProcedureDefParamName)
			runtimebuilder.LastProcedure().AddParameter()
			runtimebuilder.LastProcedure().LastParameter().SetName(t.Content)
		case s == ProcedureCallParamsStart || s == ProcedureCallParamNext:
			// check if var exists in runtimebuilder or if it is a value
			SetNextState(ProcedureCallParamValue)
			i, err := strconv.Atoi(t.Content)
			if err != nil { // it is a named variable
				ex, v := runtimebuilder.LastProcedure().VariableExists(t.Content)
				if ex {
					runtimebuilder.LastProcedure().LastInvocation().AddParameter(v)
				} else {
					log.Fatalf("Línea %d: La variable '%s' no existe", parser.LineNumber, t.Content)
				}
			} else { // it is an int literal
				runtimebuilder.LastProcedure().LastInvocation().AddParameter(&runtimebuilder.Value{Type: runtimebuilder.IntVarType, IntValue: i})
			}
		case s == VariableDefAssign:
			_, err := strconv.Atoi(t.Content) // this is a number conversion check
			if err != nil {
				// check if proc exists in runtimebuilder
				// check if current proc return value is not nil
				ex, proc := runtimebuilder.ProcedureExists(t.Content)
				if ex {
					SetParentState(VariableDefValue)
					SetNextState(ProcedureCall)
					// check if proc and variable are the same type
					// bind variable to proc invocation
					if runtimebuilder.LastProcedure().LastVariable().Type == proc.GetReturnType() {
						// add proc invocation and bind variable
						runtimebuilder.LastProcedure().AddInvocation()
						runtimebuilder.LastProcedure().LastInvocation().SetProcedure(proc)
						runtimebuilder.LastProcedure().LastInvocation().SetBoundVariable(runtimebuilder.LastProcedure().LastVariable())
					} else {
						log.Fatalf("Línea %d: La variable '%s' y el procedimiento '%s' no son del mismo tipo",
							parser.LineNumber, runtimebuilder.LastProcedure().LastVariable().Name, proc.GetName())
					}
				} else {
					log.Fatalf("Línea %d: El procedimiento '%s' no existe", parser.LineNumber, t.Content)
				}
			} else {
				SetNextState(VariableDefValue)
				runtimebuilder.LastProcedure().LastVariable().SetValue(t.Content)
			}
		case s == ProcedureDefBodyStart || s == InstructionEnd:
			// check if proc exists in runtimebuilder
			// check if current proc return value is not nil
			ex, proc := runtimebuilder.ProcedureExists(t.Content)
			if ex {
				SetNextState(ProcedureCall)
				runtimebuilder.LastProcedure().AddInvocation()
				runtimebuilder.LastProcedure().LastInvocation().SetProcedure(proc)
			} else {
				log.Fatalf("Línea %d: El procedimiento '%s' no existe", parser.LineNumber, t.Content)
			}
		case s == StringDefStart || s == StringDefValue:
			SetNextState(StringDefValue)
			stringBuffer = append(stringBuffer, t.Content)
		default:
			log.Fatalf("Línea %d: Se esperaba ';', pero se encontró '%s'", parser.LineNumber, t.Content)
		}
	}
	return true
}
