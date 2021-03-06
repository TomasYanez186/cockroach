// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License included
// in the file licenses/BSL.txt and at www.mariadb.com/bsl11.
//
// Change Date: 2022-10-01
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by the Apache License, Version 2.0,
// included in the file licenses/APL.txt and at
// https://www.apache.org/licenses/LICENSE-2.0

package tree

// ConstantEvalVisitor replaces constant TypedExprs with the result of Eval.
type ConstantEvalVisitor struct {
	ctx *EvalContext
	err error

	fastIsConstVisitor fastIsConstVisitor
}

var _ Visitor = &ConstantEvalVisitor{}

// MakeConstantEvalVisitor creates a ConstantEvalVisitor instance.
func MakeConstantEvalVisitor(ctx *EvalContext) ConstantEvalVisitor {
	return ConstantEvalVisitor{ctx: ctx, fastIsConstVisitor: fastIsConstVisitor{ctx: ctx}}
}

// Err retrieves the error field in the ConstantEvalVisitor.
func (v *ConstantEvalVisitor) Err() error { return v.err }

// VisitPre implements the Visitor interface.
func (v *ConstantEvalVisitor) VisitPre(expr Expr) (recurse bool, newExpr Expr) {
	if v.err != nil {
		return false, expr
	}
	return true, expr
}

// VisitPost implements the Visitor interface.
func (v *ConstantEvalVisitor) VisitPost(expr Expr) Expr {
	if v.err != nil {
		return expr
	}

	typedExpr, ok := expr.(TypedExpr)
	if !ok || !v.isConst(expr) {
		return expr
	}

	value, err := typedExpr.Eval(v.ctx)
	if err != nil {
		// Ignore any errors here (e.g. division by zero), so they can happen
		// during execution where they are correctly handled. Note that in some
		// cases we might not even get an error (if this particular expression
		// does not get evaluated when the query runs, e.g. it's inside a CASE).
		return expr
	}
	if value == DNull {
		// We don't want to return an expression that has a different type; cast
		// the NULL if necessary.
		var newExpr TypedExpr
		newExpr, v.err = ReType(DNull, typedExpr.ResolvedType())
		if v.err != nil {
			return expr
		}
		return newExpr
	}
	return value
}

func (v *ConstantEvalVisitor) isConst(expr Expr) bool {
	return v.fastIsConstVisitor.run(expr)
}
