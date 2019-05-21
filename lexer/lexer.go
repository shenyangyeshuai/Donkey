package lexer

import (
	"monkey/token"
)

/*
 * to "peek" further into the input ,
 * and look after the current character to see what comes up next
 *
 * `readPosition' always points to the "next" character in the input.
 *
 * `position' points to the character in the input that corresponds to the ch byte.
 *
 * `ch' only support ASCII code instead of UTF-8 for simplicity.
 * In order to fully support Unicode and UTF-8 ,
 * we would need to change l.ch from a byte to rune ,
 * and change the way we read the next characters,
 * since they could be multiple bytes wide now.
 */
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examinationposution
}

func New(input string) *Lexer {
	l := &Lexer{input: input}

	// Let's use `readChar()' in our New() function ,
	// so our *Lexer is in a fully working state before anyone calls NextToken() ,
	// with l.ch , l.position and l.readPosition already initialized
	l.readChar()

	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.ch {
	// 前面都是单字母的 identifiers
	case '+': // operators
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '=': // delimiters
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default: // token.IDENT or token.ILLEGAL , 其中 token.IDENT 分好几种情况
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok // `return' here is necessary because we shouldn't `readChar()' again
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber() // the same as `return' above
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	// Before returning the token we advance our pointers into the input ,
	// so when we call `NextToken()' again the l.ch field is already updated
	l.readChar() // 前进一步 , 为下次的 token "做好准备"

	return tok
}

//
//
// package-specific functions (methods)
//
//

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	// 像 ruby 支持 ! 和 ? 符号的, 就在这个函数里面增加条件
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch == '_')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// `peekChar()' is really similar to `readChar()' ,
// except that it doesn't increment l.position and l.readPosition .
// We only want to "peek" ahead in the input and not move around in it .
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code "NUL" (represent for "end of file")
	} else {
		l.ch = l.input[l.readPosition] // the next char
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readIdentifier() string {
	pos := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

// 这个 `readNumber()' 的实现与上面的 `readIdentifier()' 的实现很像啊!
// 有重构的机会了
func (l *Lexer) readNumber() string {
	pos := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
