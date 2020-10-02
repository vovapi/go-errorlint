package errorlint

import (
	"flag"
	"golang.org/x/tools/go/analysis"
	"sort"
)

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:  "errorlint",
		Doc:   "Source code linter for Go software that can be used to find code that will cause problems with the error wrapping scheme introduced in Go 1.13.",
		Run:   run,
		Flags: flagSet,
	}
}

var (
	flagSet     flag.FlagSet
	checkErrorf bool
)

func init() {
	flagSet.BoolVar(&checkErrorf, "errorf", false, "Check whether fmt.Errorf uses the %w verb for formatting errors. See the readme for caveats")
}

func run(pass *analysis.Pass) (interface{}, error) {
	lints := []Lint{}
	if checkErrorf {
		l := LintFmtErrorfCalls(pass.Fset, *pass.TypesInfo)
		lints = append(lints, l...)
	}
	l := LintErrorComparisons(pass.Fset, *pass.TypesInfo)
	lints = append(lints, l...)
	l = LintErrorTypeAssertions(pass.Fset, *pass.TypesInfo)
	lints = append(lints, l...)
	sort.Sort(ByPosition(lints))

	for _, l := range lints {
		pass.Report(analysis.Diagnostic{Pos: l.Pos, Message: l.Message})
	}
	return nil, nil
}
