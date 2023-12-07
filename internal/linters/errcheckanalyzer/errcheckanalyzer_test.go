package errcheckanalyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestErrCheckAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), ErrCheckAnalyzer, "./...")
}
