package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

const (
	tkInt   = "INT"
	tkFloat = "FLOAT"
	tkMinus = "MINUS"
	tkPlus  = "PLUS"
	tkMul   = "MUL"
	tkDiv   = "DIV"
	tkLp    = "LP"
	tkRp    = "RP"
)

const (
	space = ' '
	tab   = '\t'
	dot   = '.'
	plus  = '+'
	minus = '-'
	mul   = '*'
	slash = '/'
	lp    = '('
	rp    = ')'
)

const (
	endChar = '#'
	none    = "none"
)

var num = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

// Interfaces

type laxerFunctions interface {
	makeTokens() ([]token, error)
	makeNumbers() token
	advance()
}

// Token

type token struct {
	typ   string
	value interface{}
}

// Laxer

type laxer struct {
	text        string
	pos         uint
	currentChar byte
}

func (l *laxer) advance() {
	l.pos += 1
	l.currentChar = l.text[l.pos]
}

func (l *laxer) makeTokens() ([]token, error) {
	tokens := make([]token, 0, len(l.text))

	for l.currentChar != endChar {
		switch {
		case l.currentChar == space || l.currentChar == tab:
			l.advance()
		case bytes.Contains(num, []byte{l.currentChar}):
			tokens = append(tokens, l.makeNumbers())
		case l.currentChar == plus:
			tokens = append(tokens, token{tkPlus, none})
			l.advance()

		case l.currentChar == minus:
			tokens = append(tokens, token{tkMinus, none})
			l.advance()

		case l.currentChar == mul:
			tokens = append(tokens, token{tkMul, none})
			l.advance()

		case l.currentChar == slash:
			tokens = append(tokens, token{tkDiv, none})
			l.advance()

		case l.currentChar == lp:
			tokens = append(tokens, token{tkLp, none})
			l.advance()

		case l.currentChar == rp:
			tokens = append(tokens, token{tkRp, none})
			l.advance()

		default:
			return []token{}, errors.New(fmt.Sprintf("Parse error, invalid caractere %s  in \" %s \" ", string(l.currentChar), l.text[:len(l.text)-1]))
		}
	}

	return tokens, nil
}

func (l *laxer) makeNumbers() token {
	numStr := ""
	dotCount := 0

	for l.currentChar != endChar && (bytes.Contains(num, []byte{l.currentChar}) || bytes.Contains([]byte(string(dot)), []byte{l.currentChar})) {
		if l.currentChar == dot {
			if dotCount == 1 {
				break
			}
			dotCount += 1
			numStr += "."
			l.advance()
		} else {
			numStr += string(l.currentChar)
			l.advance()
		}
	}

	if dotCount == 0 {
		number, err := strconv.ParseInt(numStr, 10, 32)

		if err != nil {
			panic(err)
		}

		return token{tkInt, number}
	} else {
		number, err := strconv.ParseFloat(numStr, 64)

		if err != nil {
			panic(err)
		}

		return token{tkFloat, number}
	}
}

func New(function string) ([]token, error) {
	function = function + string(endChar)
	lxr := laxer{function, 0, function[0]}

	return lxr.makeTokens()
}
