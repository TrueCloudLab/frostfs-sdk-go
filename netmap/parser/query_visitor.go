// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // Query

import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// A complete Visitor for a parse tree produced by Query.
type QueryVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by Query#policy.
	VisitPolicy(ctx *PolicyContext) interface{}

	// Visit a parse tree produced by Query#repStmt.
	VisitRepStmt(ctx *RepStmtContext) interface{}

	// Visit a parse tree produced by Query#cbfStmt.
	VisitCbfStmt(ctx *CbfStmtContext) interface{}

	// Visit a parse tree produced by Query#selectStmt.
	VisitSelectStmt(ctx *SelectStmtContext) interface{}

	// Visit a parse tree produced by Query#clause.
	VisitClause(ctx *ClauseContext) interface{}

	// Visit a parse tree produced by Query#filterExpr.
	VisitFilterExpr(ctx *FilterExprContext) interface{}

	// Visit a parse tree produced by Query#filterStmt.
	VisitFilterStmt(ctx *FilterStmtContext) interface{}

	// Visit a parse tree produced by Query#expr.
	VisitExpr(ctx *ExprContext) interface{}

	// Visit a parse tree produced by Query#filterKey.
	VisitFilterKey(ctx *FilterKeyContext) interface{}

	// Visit a parse tree produced by Query#filterValue.
	VisitFilterValue(ctx *FilterValueContext) interface{}

	// Visit a parse tree produced by Query#number.
	VisitNumber(ctx *NumberContext) interface{}

	// Visit a parse tree produced by Query#keyword.
	VisitKeyword(ctx *KeywordContext) interface{}

	// Visit a parse tree produced by Query#ident.
	VisitIdent(ctx *IdentContext) interface{}

	// Visit a parse tree produced by Query#identWC.
	VisitIdentWC(ctx *IdentWCContext) interface{}
}
