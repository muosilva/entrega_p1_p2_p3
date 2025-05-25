package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Parser struct {
	S        string
	Position int
}

type Node interface {
	Gen(g *BFGen, cell int)
}

type BinOp struct {
	Op          byte
	Left, Right Node
}

type Number struct {
	Val int
}

type BFGen struct {
	SB       strings.Builder
	Position int
}

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	text := strings.TrimSpace(string(data))
	parts := strings.SplitN(text, "=", 2)
	if len(parts) != 2 {
		fmt.Fprintln(os.Stderr, "uso: VAR=EXPR")
		os.Exit(1)
	}
	varName, expr := parts[0], parts[1]

	p := &Parser{S: expr}
	ast := p.parseExpr()

	gen := &BFGen{}

	for _, c := range varName + "=" {
		for _, b := range []byte(string(c)) {
			gen.moveTo(10)
			gen.zero()
			gen.inc(int(b))
			gen.SB.WriteByte('.')
		}
	}

	ast.Gen(gen, 0)

	result := eValNode(ast)
	for _, ch := range strconv.Itoa(result) {
		gen.moveTo(10)
		gen.zero()
		gen.inc(int(ch))
		gen.SB.WriteByte('.')
	}

	fmt.Println(gen.String())
}

func (p *Parser) peek() rune {
	if p.Position >= len(p.S) {
		return 0
	}
	return rune(p.S[p.Position])
}

func (p *Parser) consume() rune {
	ch := p.peek()
	if ch != 0 {
		p.Position++
	}
	return ch
}

func (p *Parser) parseExpr() Node {
	node := p.parseTerm()
	for {
		switch p.peek() {
		case '+', '-':
			Op := byte(p.consume())
			Right := p.parseTerm()
			node = &BinOp{Op: Op, Left: node, Right: Right}
		default:
			return node
		}
	}
}

func (p *Parser) parseTerm() Node {
	node := p.parseFactor()
	for p.peek() == '*' {
		p.consume()
		Right := p.parseFactor()
		node = &BinOp{Op: '*', Left: node, Right: Right}
	}
	return node
}

func (p *Parser) parseFactor() Node {
	if p.peek() == '(' {
		p.consume()
		node := p.parseExpr()
		if p.peek() == ')' {
			p.consume()
		}
		return node
	}
	start := p.Position
	for unicode.IsDigit(p.peek()) {
		p.consume()
	}
	numStr := p.S[start:p.Position]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "número inválido: %S\n", numStr)
		os.Exit(1)
	}
	return &Number{Val: num}
}

func (n *Number) Gen(g *BFGen, cell int) {
	g.moveTo(cell)
	g.zero()
	g.inc(n.Val)
}

func (b *BinOp) Gen(g *BFGen, cell int) {
	switch b.Op {
	case '+':
		b.Left.Gen(g, cell)
		b.Right.Gen(g, cell+1)
		g.emitAdd(cell+1, cell)
	case '-':
		b.Left.Gen(g, cell)
		b.Right.Gen(g, cell+1)
		g.emitSub(cell+1, cell)
	case '*':
		b.Left.Gen(g, cell)
		b.Right.Gen(g, cell+1)
		g.emitMul(cell, cell+1, cell+2, cell+3)
	}
}

func eValNode(n Node) int {
	switch v := n.(type) {
	case *Number:
		return v.Val
	case *BinOp:
		L := eValNode(v.Left)
		R := eValNode(v.Right)
		switch v.Op {
		case '+':
			return L + R
		case '-':
			return L - R
		case '*':
			return L * R
		}
	}
	return 0
}

func (g *BFGen) moveTo(c int) {
	for g.Position < c {
		g.SB.WriteByte('>')
		g.Position++
	}
	for g.Position > c {
		g.SB.WriteByte('<')
		g.Position--
	}
}

func (g *BFGen) zero() {
	g.SB.WriteString("[-]")
}

func (g *BFGen) inc(n int) {
	for i := 0; i < n; i++ {
		g.SB.WriteByte('+')
	}
}

func (g *BFGen) emitLoOp(c int, body func()) {
	g.moveTo(c)
	g.SB.WriteByte('[')
	body()
	g.moveTo(c)
	g.SB.WriteByte(']')
}

func (g *BFGen) emitAdd(src, dst int) {
	g.emitLoOp(src, func() {
		g.SB.WriteByte('-')
		g.moveTo(dst)
		g.SB.WriteByte('+')
		g.moveTo(src)
	})
	g.moveTo(dst)
}

func (g *BFGen) emitSub(src, dst int) {
	g.emitLoOp(src, func() {
		g.SB.WriteByte('-')
		g.moveTo(dst)
		g.SB.WriteByte('-')
		g.moveTo(src)
	})
	g.moveTo(dst)
}

func (g *BFGen) emitMul(a, b, res, tmp int) {
	g.moveTo(res)
	g.zero()
	g.moveTo(tmp)
	g.zero()
	g.emitLoOp(a, func() {
		g.moveTo(a)
		g.SB.WriteByte('-')
		g.emitLoOp(b, func() {
			g.SB.WriteByte('-')
			g.moveTo(res)
			g.SB.WriteByte('+')
			g.moveTo(tmp)
			g.SB.WriteByte('+')
			g.moveTo(b)
		})
		g.emitLoOp(tmp, func() {
			g.SB.WriteByte('-')
			g.moveTo(b)
			g.SB.WriteByte('+')
			g.moveTo(tmp)
		})
	})
	g.emitLoOp(res, func() {
		g.SB.WriteByte('-')
		g.moveTo(a)
		g.SB.WriteByte('+')
		g.moveTo(res)
	})
	g.moveTo(a)
}

func (g *BFGen) String() string {
	return g.SB.String()
}
