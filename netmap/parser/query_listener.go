// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // Query

import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// QueryListener is a complete listener for a parse tree produced by Query.
type QueryListener interface {
	antlr.ParseTreeListener

	// EnterPolicy is called when entering the policy production.
	EnterPolicy(c *PolicyContext)

	// EnterRepStmt is called when entering the repStmt production.
	EnterRepStmt(c *RepStmtContext)

	// EnterCbfStmt is called when entering the cbfStmt production.
	EnterCbfStmt(c *CbfStmtContext)

	// EnterSelectStmt is called when entering the selectStmt production.
	EnterSelectStmt(c *SelectStmtContext)

	// EnterClause is called when entering the clause production.
	EnterClause(c *ClauseContext)

	// EnterFilterExpr is called when entering the filterExpr production.
	EnterFilterExpr(c *FilterExprContext)

	// EnterFilterStmt is called when entering the filterStmt production.
	EnterFilterStmt(c *FilterStmtContext)

	// EnterExpr is called when entering the expr production.
	EnterExpr(c *ExprContext)

	// EnterFilterKey is called when entering the filterKey production.
	EnterFilterKey(c *FilterKeyContext)

	// EnterFilterValue is called when entering the filterValue production.
	EnterFilterValue(c *FilterValueContext)

	// EnterNumber is called when entering the number production.
	EnterNumber(c *NumberContext)

	// EnterKeyword is called when entering the keyword production.
	EnterKeyword(c *KeywordContext)

	// EnterIdent is called when entering the ident production.
	EnterIdent(c *IdentContext)

	// EnterIdentWC is called when entering the identWC production.
	EnterIdentWC(c *IdentWCContext)

	// ExitPolicy is called when exiting the policy production.
	ExitPolicy(c *PolicyContext)

	// ExitRepStmt is called when exiting the repStmt production.
	ExitRepStmt(c *RepStmtContext)

	// ExitCbfStmt is called when exiting the cbfStmt production.
	ExitCbfStmt(c *CbfStmtContext)

	// ExitSelectStmt is called when exiting the selectStmt production.
	ExitSelectStmt(c *SelectStmtContext)

	// ExitClause is called when exiting the clause production.
	ExitClause(c *ClauseContext)

	// ExitFilterExpr is called when exiting the filterExpr production.
	ExitFilterExpr(c *FilterExprContext)

	// ExitFilterStmt is called when exiting the filterStmt production.
	ExitFilterStmt(c *FilterStmtContext)

	// ExitExpr is called when exiting the expr production.
	ExitExpr(c *ExprContext)

	// ExitFilterKey is called when exiting the filterKey production.
	ExitFilterKey(c *FilterKeyContext)

	// ExitFilterValue is called when exiting the filterValue production.
	ExitFilterValue(c *FilterValueContext)

	// ExitNumber is called when exiting the number production.
	ExitNumber(c *NumberContext)

	// ExitKeyword is called when exiting the keyword production.
	ExitKeyword(c *KeywordContext)

	// ExitIdent is called when exiting the ident production.
	ExitIdent(c *IdentContext)

	// ExitIdentWC is called when exiting the identWC production.
	ExitIdentWC(c *IdentWCContext)
}
