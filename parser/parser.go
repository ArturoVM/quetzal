package parser

import (
	"bufio"
	"errors"
	"os"
)

type TokenType int

const (
	Whitespace TokenType = iota
	Word
	Symbol
)

type Token struct {
	Content string
	Type    TokenType
}

var sourceReader *bufio.Reader
var LineNumber int = 1

func tokenBuilder(c string, t TokenType) *Token {
	return &Token{
		c,
		t,
	}
}

func NextToken() (*Token, error) {
	if sourceReader == nil {
		return nil, errors.New("No source file loaded")
	}
	wordBuffer := make([]rune, 0)
	for {
		r, _, err := sourceReader.ReadRune()
		if err != nil {
			return nil, err
		}
		switch r {
		case ' ':
			if len(wordBuffer) > 0 {
				return tokenBuilder(string(wordBuffer), Word), nil
			}
		case '\t':
			if len(wordBuffer) > 0 {
				return tokenBuilder(string(wordBuffer), Word), nil
			}
		case ';':
			if len(wordBuffer) > 0 {
				sourceReader.UnreadRune()
				return tokenBuilder(string(wordBuffer), Word), nil
			}
			return tokenBuilder(";", Symbol), nil
		case '(':
			if len(wordBuffer) > 0 {
				sourceReader.UnreadRune()
				return tokenBuilder(string(wordBuffer), Word), nil
			}
			return tokenBuilder("(", Symbol), nil
		case ')':
			if len(wordBuffer) > 0 {
				sourceReader.UnreadRune()
				return tokenBuilder(string(wordBuffer), Word), nil
			}
			return tokenBuilder(")", Symbol), nil
		case '{':
			if len(wordBuffer) > 0 {
				sourceReader.UnreadRune()
				return tokenBuilder(string(wordBuffer), Word), nil
			}
			return tokenBuilder("{", Symbol), nil
		case '}':
			if len(wordBuffer) > 0 {
				sourceReader.UnreadRune()
				return tokenBuilder(string(wordBuffer), Word), nil
			}
			return tokenBuilder("}", Symbol), nil
		case ',':
			if len(wordBuffer) > 0 {
				sourceReader.UnreadRune()
				return tokenBuilder(string(wordBuffer), Word), nil
			}
			return tokenBuilder(",", Symbol), nil
		case '=':
			if len(wordBuffer) > 0 {
				sourceReader.UnreadRune()
				return tokenBuilder(string(wordBuffer), Word), nil
			}
			return tokenBuilder("=", Symbol), nil
		case '"':
			if len(wordBuffer) > 0 {
				sourceReader.UnreadRune()
				return tokenBuilder(string(wordBuffer), Word), nil
			}
			return tokenBuilder("\"", Symbol), nil
		case '\n':
			LineNumber += 1
			if len(wordBuffer) > 0 {
				return tokenBuilder(string(wordBuffer), Word), nil
			}
			return tokenBuilder("\n", Symbol), nil
		default:
			wordBuffer = append(wordBuffer, r)
		}
	}
}

func Read(path string) error {
	sourceFile, err := os.Open(path)
	if err != nil {
		return err
	}
	sourceReader = bufio.NewReader(sourceFile)
	return nil
}
