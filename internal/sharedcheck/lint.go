package sharedcheck

import (
	"go/ast"
	"go/types"

	"github.com/xklalala/go-tools/go/ast/astutil"
	"github.com/xklalala/go-tools/go/ir"
	"github.com/xklalala/go-tools/go/ir/irutil"
	"github.com/xklalala/go-tools/internal/passes/buildir"

	"golang.org/x/tools/go/analysis"
)

func CheckRangeStringRunes(pass *analysis.Pass) (interface{}, error) {
	for _, fn := range pass.ResultOf[buildir.Analyzer].(*buildir.IR).SrcFuncs {
		cb := func(node ast.Node) bool {
			rng, ok := node.(*ast.RangeStmt)
			if !ok || !astutil.IsBlank(rng.Key) {
				return true
			}

			v, _ := fn.ValueForExpr(rng.X)

			// Check that we're converting from string to []rune
			val, _ := v.(*ir.Convert)
			if val == nil {
				return true
			}
			Tsrc, ok := val.X.Type().Underlying().(*types.Basic)
			if !ok || Tsrc.Kind() != types.String {
				return true
			}
			Tdst, ok := val.Type().(*types.Slice)
			if !ok {
				return true
			}
			TdstElem, ok := Tdst.Elem().(*types.Basic)
			if !ok || TdstElem.Kind() != types.Int32 {
				return true
			}

			// Check that the result of the conversion is only used to
			// range over
			refs := val.Referrers()
			if refs == nil {
				return true
			}

			// Expect two refs: one for obtaining the length of the slice,
			// one for accessing the elements
			if len(irutil.FilterDebug(*refs)) != 2 {
				// TODO(dh): right now, we check that only one place
				// refers to our slice. This will miss cases such as
				// ranging over the slice twice. Ideally, we'd ensure that
				// the slice is only used for ranging over (without
				// accessing the key), but that is harder to do because in
				// IR form, ranging over a slice looks like an ordinary
				// loop with index increments and slice accesses. We'd
				// have to look at the associated AST node to check that
				// it's a range statement.
				return true
			}

			pass.Reportf(rng.Pos(), "should range over string, not []rune(string)")

			return true
		}
		if source := fn.Source(); source != nil {
			ast.Inspect(source, cb)
		}
	}
	return nil, nil
}
