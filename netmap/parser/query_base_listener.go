// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // Query

import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// BaseQueryListener is a complete listener for a parse tree produced by Query.
type BaseQueryListener struct{}

var _ QueryListener = &BaseQueryListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseQueryListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseQueryListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseQueryListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseQueryListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterPolicy is called when production policy is entered.
func (s *BaseQueryListener) EnterPolicy(ctx *PolicyContext) {}

// ExitPolicy is called when production policy is exited.
func (s *BaseQueryListener) ExitPolicy(ctx *PolicyContext) {}

// EnterRepStmt is called when production repStmt is entered.
func (s *BaseQueryListener) EnterRepStmt(ctx *RepStmtContext) {}

// ExitRepStmt is called when production repStmt is exited.
func (s *BaseQueryListener) ExitRepStmt(ctx *RepStmtContext) {}

// EnterCbfStmt is called when production cbfStmt is entered.
func (s *BaseQueryListener) EnterCbfStmt(ctx *CbfStmtContext) {}

// ExitCbfStmt is called when production cbfStmt is exited.
func (s *BaseQueryListener) ExitCbfStmt(ctx *CbfStmtContext) {}

// EnterSelectStmt is called when production selectStmt is entered.
func (s *BaseQueryListener) EnterSelectStmt(ctx *SelectStmtContext) {}

// ExitSelectStmt is called when production selectStmt is exited.
func (s *BaseQueryListener) ExitSelectStmt(ctx *SelectStmtContext) {}

// EnterClause is called when production clause is entered.
func (s *BaseQueryListener) EnterClause(ctx *ClauseContext) {}

// ExitClause is called when production clause is exited.
func (s *BaseQueryListener) ExitClause(ctx *ClauseContext) {}

// EnterFilterExpr is called when production filterExpr is entered.
func (s *BaseQueryListener) EnterFilterExpr(ctx *FilterExprContext) {}

// ExitFilterExpr is called when production filterExpr is exited.
func (s *BaseQueryListener) ExitFilterExpr(ctx *FilterExprContext) {}

// EnterFilterStmt is called when production filterStmt is entered.
func (s *BaseQueryListener) EnterFilterStmt(ctx *FilterStmtContext) {}

// ExitFilterStmt is called when production filterStmt is exited.
func (s *BaseQueryListener) ExitFilterStmt(ctx *FilterStmtContext) {}

// EnterExpr is called when production expr is entered.
func (s *BaseQueryListener) EnterExpr(ctx *ExprContext) {}

// ExitExpr is called when production expr is exited.
func (s *BaseQueryListener) ExitExpr(ctx *ExprContext) {}

// EnterFilterKey is called when production filterKey is entered.
func (s *BaseQueryListener) EnterFilterKey(ctx *FilterKeyContext) {}

// ExitFilterKey is called when production filterKey is exited.
func (s *BaseQueryListener) ExitFilterKey(ctx *FilterKeyContext) {}

// EnterFilterValue is called when production filterValue is entered.
func (s *BaseQueryListener) EnterFilterValue(ctx *FilterValueContext) {}

// ExitFilterValue is called when production filterValue is exited.
func (s *BaseQueryListener) ExitFilterValue(ctx *FilterValueContext) {}

// EnterNumber is called when production number is entered.
func (s *BaseQueryListener) EnterNumber(ctx *NumberContext) {}

// ExitNumber is called when production number is exited.
func (s *BaseQueryListener) ExitNumber(ctx *NumberContext) {}

// EnterKeyword is called when production keyword is entered.
func (s *BaseQueryListener) EnterKeyword(ctx *KeywordContext) {}

// ExitKeyword is called when production keyword is exited.
func (s *BaseQueryListener) ExitKeyword(ctx *KeywordContext) {}

// EnterIdent is called when production ident is entered.
func (s *BaseQueryListener) EnterIdent(ctx *IdentContext) {}

// ExitIdent is called when production ident is exited.
func (s *BaseQueryListener) ExitIdent(ctx *IdentContext) {}

// EnterIdentWC is called when production identWC is entered.
func (s *BaseQueryListener) EnterIdentWC(ctx *IdentWCContext) {}

// ExitIdentWC is called when production identWC is exited.
func (s *BaseQueryListener) ExitIdentWC(ctx *IdentWCContext) {}
