// Package osexitanalyzer - пакет статического анализа.
// Проверяет наличие вызова os.Exit() в main функции main пакета
package osexitanalyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	packageName = "main"
	funcName
	ignoreFileNamePattern = "go-build"
)

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for os exit in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.String() != packageName {
			continue
		}

		srcFile := pass.Fset.File(file.Pos())
		if srcFile != nil && strings.Contains(srcFile.Name(), ignoreFileNamePattern) {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			if x, ok := node.(*ast.FuncDecl); ok && x.Name.String() == funcName {
				ast.Inspect(x.Body, func(nn ast.Node) bool {
					switch xx := nn.(type) {
					case *ast.CallExpr:
						if s, selOk := xx.Fun.(*ast.SelectorExpr); selOk {
							if from, identOk := s.X.(*ast.Ident); identOk {
								if from.Name == "os" && s.Sel.Name == "Exit" {
									pass.Reportf(s.Pos(), "os exit in the main func")
								}
							}
						}
					}
					return true
				})

				return false
			}
			return true
		})
	}

	return nil, nil
}
