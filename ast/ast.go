package ast

import "monkey/token"

// Node 表示 AST 的一个节点。
// AST 中的每个节点都必须实现 Node 接口。
type Node interface {
	// TokenLiteral 返回与其相关联的词法单元的字面量，仅用于调试和测试。
	TokenLiteral() string
}

// Statement 表示 AST 的语句节点。
type Statement interface {
	Node
	statementNode()
}

// Expression 表示 AST 的表达式节点
type Expression interface {
	Node
	expressionNode()
}

// Program 节点将成为词法分析器生成的每个 AST 的根节点。
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

type LetStatement struct {
	Token token.Token // token.LET 词法单元
	Name  *Identifier // 用来保存绑定的标识符
	Value Expression  // 用于表示产生值的表达式
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token // token.IDENT 词法单元
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
