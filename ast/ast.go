package ast

import (
	"monkey/token"
)

// interface: Node
type Node interface {
	TokenLiteral() string
}

// interface: Statment
type Statement interface {
	Node
	statementNode()
}

// interface: Expression
type Expression interface {
	Node
	expressionNode()
}

// struct: Programm
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
