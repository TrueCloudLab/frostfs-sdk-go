// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // Query

import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

type BaseQueryVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseQueryVisitor) VisitPolicy(ctx *PolicyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitRepStmt(ctx *RepStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitCbfStmt(ctx *CbfStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitSelectStmt(ctx *SelectStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitClause(ctx *ClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitFilterExpr(ctx *FilterExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitFilterStmt(ctx *FilterStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitExpr(ctx *ExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitFilterKey(ctx *FilterKeyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitFilterValue(ctx *FilterValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitNumber(ctx *NumberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitKeyword(ctx *KeywordContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitIdent(ctx *IdentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitIdentWC(ctx *IdentWCContext) interface{} {
	return v.VisitChildren(ctx)
}
