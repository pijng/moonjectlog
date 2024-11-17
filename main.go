package main

import (
	"fmt"
	"go/token"
	"strconv"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/pijng/goinject"
)

type logModifier struct{}

func (lm logModifier) Modify(f *dst.File, dec *decorator.Decorator, res *decorator.Restorer) *dst.File {
	for _, decl := range f.Decls {
		decl, isFunc := decl.(*dst.FuncDecl)
		if !isFunc {
			continue
		}

		spanStmt := buildSpan(decl.Name.Name)
		decl.Body.List = append(spanStmt.List, decl.Body.List...)
	}

	return f
}

func main() {
	goinject.Process(logModifier{})
}

func buildSpan(funcName string) dst.BlockStmt {
	return dst.BlockStmt{
		List: []dst.Stmt{
			&dst.AssignStmt{
				Lhs: []dst.Expr{
					&dst.Ident{Name: "now"},
				},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{
					&dst.CallExpr{
						Fun:  &dst.Ident{Path: "time", Name: "Now"},
						Args: []dst.Expr{},
					},
				},
			},
			&dst.DeferStmt{
				Call: &dst.CallExpr{
					Fun: &dst.FuncLit{
						Type: &dst.FuncType{
							Params: &dst.FieldList{
								List: []*dst.Field{
									{
										Names: []*dst.Ident{{Name: "t"}},
										Type: &dst.SelectorExpr{
											X:   &dst.Ident{Name: "time"},
											Sel: &dst.Ident{Name: "Time"},
										},
									},
								},
							},
						},
						Body: &dst.BlockStmt{
							List: []dst.Stmt{
								&dst.ExprStmt{
									X: &dst.CallExpr{
										Fun: &dst.Ident{Path: "fmt", Name: "Println"},
										Args: []dst.Expr{
											&dst.BasicLit{Kind: token.STRING, Value: strconv.Quote(fmt.Sprintf("[%s] took: ", funcName))},
											&dst.CallExpr{
												Fun: &dst.SelectorExpr{
													X:   &dst.Ident{Name: "time"},
													Sel: &dst.Ident{Name: "Since"},
												},
												Args: []dst.Expr{
													&dst.Ident{Name: "t"},
												},
											},
										},
									},
								},
							},
						},
					},
					Args: []dst.Expr{
						&dst.Ident{Name: "now"},
					},
				},
			},
		},
	}
}
