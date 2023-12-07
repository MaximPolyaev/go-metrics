// Статический анализатор, состоящий из:
// - стандартных статических анализаторов пакета golang.org/x/tools/go/analysis/passes;
// - всех анализаторов класса SA пакета staticcheck.io;
// - не менее одного анализатора остальных классов пакета staticcheck.io;
// - двух или более любых публичных анализаторов;
// - собственный анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакета main;
// - собственный анализатор, проверяющий игнорирование ошибок.
package main

import (
	"github.com/MaximPolyaev/go-metrics/internal/linters/errcheckanalyzer"
	"github.com/MaximPolyaev/go-metrics/internal/linters/osexitanalyzer"
	ifshort "github.com/esimonov/ifshort/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// checkVarsAndPackageNamesRule проверка формата имен переменных и пакетов (https://staticcheck.dev/docs/checks/#ST1003)
const checkVarsAndPackageNamesRule = "ST1003"

func main() {
	myChecks := []*analysis.Analyzer{
		// проверка игнорироваться ошибок
		errcheckanalyzer.ErrCheckAnalyzer,
		// проверка вызова os.Exit() в main
		osexitanalyzer.OsExitAnalyzer,
		// проверка соответствия сигнатуры паттерна и аргументов printf
		printf.Analyzer,
		// проверка затенения переменных
		shadow.Analyzer,
		// проверка формата тегов структур
		structtag.Analyzer,
		// проверка сдвигов превышающего целого числа
		shift.Analyzer,
		// проверка возможности использовать короткий синтаксис if
		ifshort.Analyzer,
	}

	// все анализаторы класса SA пакета staticcheck.io
	for _, v := range staticcheck.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}

	for _, v := range stylecheck.Analyzers {
		if v.Analyzer.Name == checkVarsAndPackageNamesRule {
			myChecks = append(myChecks, v.Analyzer)
		}
	}

	multichecker.Main(myChecks...)
}
