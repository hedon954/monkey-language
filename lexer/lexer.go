package lexer

import "monkey/token"

type Lexer struct {
	input        string // 输入的源代码字符串
	position     int    // 所输入字符串中的当前位置（指向当前字符）
	readPosition int    // 所输入字符串中的当前读取位置（指向当前字符之后的下一个字符）
	ch           byte   // 当前正在查看的字符
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar 的目的是读取 input 的下一个字符，并前移其在 input 中的位置。
// 具体过程如下：
//  1. 检查是否已经达到 input 的末尾，如果是，则将 l.ch 置为 0，这是 NULL 字符的 ASCII 编码，用于表示“尚未读取任何内容”或“文件结尾”。
//  2. 如果还没达到 input 的末尾，则：
//     2.1 将 l.ch 设置为下一个字符。
//     2.2 将 position 更新为刚用过的 l.readPosition
//     2.3 将 l.readPosition++，使其指向下一个将要读取的字符位置
//
// NOTE:
//
//	这意味着该语法分析器仅支持 **ASCII** 字符，不支持所有的 Unicode 字符。
//	如果要完全支持 Unicode 和 UTC-8，这需要将 l.ch 改成 rune 类型，并修改响应字段的逻辑。
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// NextToken 首先检查当前正在查看的字符 l.ch 并返回相应的词法单元 token.TokenType，
// 返回调用 readChar 前进一个字符。
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skinWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
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
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) { // 处理标识符和关键字
			// 读取一个完整的字母
			tok.Literal = l.readIdentifier()
			// 判断是标识符还是关键字
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) { // 处理数字字面量
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func newToken(tokType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokType, Literal: string(ch)}
}

// readIdentifier 读取一个完整的标识符
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber 读取一个完整的数字
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// peekChar 跟 readChar 很相似，但是它不会前移 l.position 和 l.readPosition，
// 仅仅只是窥视输入中的下一个字符。
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// isLetter 判断给定的字符是否为字母
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit 判断给定的字符是否为数字
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// skinWhitespace 跳过空白字符
func (l *Lexer) skinWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}
