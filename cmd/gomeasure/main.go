package main

import (
	"go/token"
	"strconv"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/pijng/goinject"
)

type measureModifier struct{}

func (mm measureModifier) Modify(f *dst.File, dec *decorator.Decorator, res *decorator.Restorer) *dst.File {

	for _, decl := range f.Decls {
		funcDecl, isFunc := decl.(*dst.FuncDecl)
		if !isFunc {
			continue
		}

		spanStmt := buildSpan(funcDecl.Name.Name)
		funcDecl.Body.List = append(spanStmt.List, funcDecl.Body.List...)
	}

	return f
}

func main() {
	goinject.Process(measureModifier{})
}

// span := StartSpan("main")
// defer() {
//   span.End()
// }()

func buildSpan(funcName string) dst.BlockStmt {
	return dst.BlockStmt{
		List: []dst.Stmt{
			&dst.AssignStmt{
				Lhs: []dst.Expr{
					&dst.Ident{Name: "span"},
				},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{
					&dst.CallExpr{
						Fun: &dst.Ident{Path: "gomeasure", Name: "StartSpan"},
						Args: []dst.Expr{
							&dst.BasicLit{Kind: token.STRING, Value: strconv.Quote(funcName)},
						},
					},
				},
			},
			&dst.DeferStmt{
				Call: &dst.CallExpr{
					Fun: &dst.FuncLit{
						Type: &dst.FuncType{
							Params: &dst.FieldList{
								List: []*dst.Field{},
							},
						},
						Body: &dst.BlockStmt{
							List: []dst.Stmt{
								&dst.ExprStmt{
									X: &dst.CallExpr{
										Fun:  &dst.Ident{Path: "span", Name: "End"},
										Args: []dst.Expr{},
									},
								},
							},
						},
					},
					// Args: []dst.Expr{
					// },
				},
			},
		},
	}
}
