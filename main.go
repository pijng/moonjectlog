package main

import (
	"fmt"
	"go/token"
	"strconv"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/pijng/moonject"
)

type logModifier struct{}

func (lm logModifier) Modify(f *dst.File, dec *decorator.Decorator, res *decorator.Restorer) *dst.File {
	for _, decl := range f.Decls {
		decl, isFunc := decl.(*dst.FuncDecl)
		if !isFunc {
			continue
		}

		span := buildSpan(decl.Name.Name)
		decl.Body.List = append(span, decl.Body.List...)
	}

	return f
}

func main() {
	moonject.Process(logModifier{})
}

func buildSpan(funcName string) []dst.Stmt {
	return []dst.Stmt{
		&dst.ExprStmt{
			X: &dst.CallExpr{
				Fun: &dst.Ident{Path: "fmt", Name: "Println"},
				Args: []dst.Expr{
					&dst.BasicLit{Kind: token.STRING, Value: strconv.Quote(fmt.Sprintf("Calling [%s] func", funcName))},
				},
			},
		},
	}
}
