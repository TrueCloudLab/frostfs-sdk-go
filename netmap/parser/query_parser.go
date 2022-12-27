// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // Query

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type Query struct {
	*antlr.BaseParser
}

var queryParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func queryParserInit() {
	staticData := &queryParserStaticData
	staticData.literalNames = []string{
		"", "'AND'", "'OR'", "", "'REP'", "'IN'", "'AS'", "'CBF'", "'SELECT'",
		"'FROM'", "'FILTER'", "'*'", "'SAME'", "'DISTINCT'", "'('", "')'", "'@'",
		"", "", "'0'",
	}
	staticData.symbolicNames = []string{
		"", "AND_OP", "OR_OP", "SIMPLE_OP", "REP", "IN", "AS", "CBF", "SELECT",
		"FROM", "FILTER", "WILDCARD", "CLAUSE_SAME", "CLAUSE_DISTINCT", "L_PAREN",
		"R_PAREN", "AT", "IDENT", "NUMBER1", "ZERO", "STRING", "WS",
	}
	staticData.ruleNames = []string{
		"policy", "repStmt", "cbfStmt", "selectStmt", "clause", "filterExpr",
		"filterStmt", "expr", "filterKey", "filterValue", "number", "keyword",
		"ident", "identWC",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 21, 130, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 1, 0, 4, 0, 30, 8, 0, 11,
		0, 12, 0, 31, 1, 0, 3, 0, 35, 8, 0, 1, 0, 5, 0, 38, 8, 0, 10, 0, 12, 0,
		41, 9, 0, 1, 0, 5, 0, 44, 8, 0, 10, 0, 12, 0, 47, 9, 0, 1, 0, 1, 0, 1,
		1, 1, 1, 1, 1, 1, 1, 3, 1, 55, 8, 1, 1, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3,
		1, 3, 3, 3, 64, 8, 3, 1, 3, 3, 3, 67, 8, 3, 1, 3, 1, 3, 1, 3, 1, 3, 3,
		3, 73, 8, 3, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 83,
		8, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 5, 5, 91, 8, 5, 10, 5, 12, 5,
		94, 9, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1,
		7, 3, 7, 107, 8, 7, 1, 8, 1, 8, 3, 8, 111, 8, 8, 1, 9, 1, 9, 1, 9, 3, 9,
		116, 8, 9, 1, 10, 1, 10, 1, 11, 1, 11, 1, 12, 1, 12, 3, 12, 124, 8, 12,
		1, 13, 1, 13, 3, 13, 128, 8, 13, 1, 13, 0, 1, 10, 14, 0, 2, 4, 6, 8, 10,
		12, 14, 16, 18, 20, 22, 24, 26, 0, 3, 1, 0, 12, 13, 1, 0, 18, 19, 2, 0,
		4, 6, 8, 10, 132, 0, 29, 1, 0, 0, 0, 2, 50, 1, 0, 0, 0, 4, 56, 1, 0, 0,
		0, 6, 59, 1, 0, 0, 0, 8, 74, 1, 0, 0, 0, 10, 82, 1, 0, 0, 0, 12, 95, 1,
		0, 0, 0, 14, 106, 1, 0, 0, 0, 16, 110, 1, 0, 0, 0, 18, 115, 1, 0, 0, 0,
		20, 117, 1, 0, 0, 0, 22, 119, 1, 0, 0, 0, 24, 123, 1, 0, 0, 0, 26, 127,
		1, 0, 0, 0, 28, 30, 3, 2, 1, 0, 29, 28, 1, 0, 0, 0, 30, 31, 1, 0, 0, 0,
		31, 29, 1, 0, 0, 0, 31, 32, 1, 0, 0, 0, 32, 34, 1, 0, 0, 0, 33, 35, 3,
		4, 2, 0, 34, 33, 1, 0, 0, 0, 34, 35, 1, 0, 0, 0, 35, 39, 1, 0, 0, 0, 36,
		38, 3, 6, 3, 0, 37, 36, 1, 0, 0, 0, 38, 41, 1, 0, 0, 0, 39, 37, 1, 0, 0,
		0, 39, 40, 1, 0, 0, 0, 40, 45, 1, 0, 0, 0, 41, 39, 1, 0, 0, 0, 42, 44,
		3, 12, 6, 0, 43, 42, 1, 0, 0, 0, 44, 47, 1, 0, 0, 0, 45, 43, 1, 0, 0, 0,
		45, 46, 1, 0, 0, 0, 46, 48, 1, 0, 0, 0, 47, 45, 1, 0, 0, 0, 48, 49, 5,
		0, 0, 1, 49, 1, 1, 0, 0, 0, 50, 51, 5, 4, 0, 0, 51, 54, 5, 18, 0, 0, 52,
		53, 5, 5, 0, 0, 53, 55, 3, 24, 12, 0, 54, 52, 1, 0, 0, 0, 54, 55, 1, 0,
		0, 0, 55, 3, 1, 0, 0, 0, 56, 57, 5, 7, 0, 0, 57, 58, 5, 18, 0, 0, 58, 5,
		1, 0, 0, 0, 59, 60, 5, 8, 0, 0, 60, 66, 5, 18, 0, 0, 61, 63, 5, 5, 0, 0,
		62, 64, 3, 8, 4, 0, 63, 62, 1, 0, 0, 0, 63, 64, 1, 0, 0, 0, 64, 65, 1,
		0, 0, 0, 65, 67, 3, 24, 12, 0, 66, 61, 1, 0, 0, 0, 66, 67, 1, 0, 0, 0,
		67, 68, 1, 0, 0, 0, 68, 69, 5, 9, 0, 0, 69, 72, 3, 26, 13, 0, 70, 71, 5,
		6, 0, 0, 71, 73, 3, 24, 12, 0, 72, 70, 1, 0, 0, 0, 72, 73, 1, 0, 0, 0,
		73, 7, 1, 0, 0, 0, 74, 75, 7, 0, 0, 0, 75, 9, 1, 0, 0, 0, 76, 77, 6, 5,
		-1, 0, 77, 78, 5, 14, 0, 0, 78, 79, 3, 10, 5, 0, 79, 80, 5, 15, 0, 0, 80,
		83, 1, 0, 0, 0, 81, 83, 3, 14, 7, 0, 82, 76, 1, 0, 0, 0, 82, 81, 1, 0,
		0, 0, 83, 92, 1, 0, 0, 0, 84, 85, 10, 4, 0, 0, 85, 86, 5, 1, 0, 0, 86,
		91, 3, 10, 5, 5, 87, 88, 10, 3, 0, 0, 88, 89, 5, 2, 0, 0, 89, 91, 3, 10,
		5, 4, 90, 84, 1, 0, 0, 0, 90, 87, 1, 0, 0, 0, 91, 94, 1, 0, 0, 0, 92, 90,
		1, 0, 0, 0, 92, 93, 1, 0, 0, 0, 93, 11, 1, 0, 0, 0, 94, 92, 1, 0, 0, 0,
		95, 96, 5, 10, 0, 0, 96, 97, 3, 10, 5, 0, 97, 98, 5, 6, 0, 0, 98, 99, 3,
		24, 12, 0, 99, 13, 1, 0, 0, 0, 100, 101, 5, 16, 0, 0, 101, 107, 3, 24,
		12, 0, 102, 103, 3, 16, 8, 0, 103, 104, 5, 3, 0, 0, 104, 105, 3, 18, 9,
		0, 105, 107, 1, 0, 0, 0, 106, 100, 1, 0, 0, 0, 106, 102, 1, 0, 0, 0, 107,
		15, 1, 0, 0, 0, 108, 111, 3, 24, 12, 0, 109, 111, 5, 20, 0, 0, 110, 108,
		1, 0, 0, 0, 110, 109, 1, 0, 0, 0, 111, 17, 1, 0, 0, 0, 112, 116, 3, 24,
		12, 0, 113, 116, 3, 20, 10, 0, 114, 116, 5, 20, 0, 0, 115, 112, 1, 0, 0,
		0, 115, 113, 1, 0, 0, 0, 115, 114, 1, 0, 0, 0, 116, 19, 1, 0, 0, 0, 117,
		118, 7, 1, 0, 0, 118, 21, 1, 0, 0, 0, 119, 120, 7, 2, 0, 0, 120, 23, 1,
		0, 0, 0, 121, 124, 3, 22, 11, 0, 122, 124, 5, 17, 0, 0, 123, 121, 1, 0,
		0, 0, 123, 122, 1, 0, 0, 0, 124, 25, 1, 0, 0, 0, 125, 128, 3, 24, 12, 0,
		126, 128, 5, 11, 0, 0, 127, 125, 1, 0, 0, 0, 127, 126, 1, 0, 0, 0, 128,
		27, 1, 0, 0, 0, 16, 31, 34, 39, 45, 54, 63, 66, 72, 82, 90, 92, 106, 110,
		115, 123, 127,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// QueryInit initializes any static state used to implement Query. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewQuery(). You can call this function if you wish to initialize the static state ahead
// of time.
func QueryInit() {
	staticData := &queryParserStaticData
	staticData.once.Do(queryParserInit)
}

// NewQuery produces a new parser instance for the optional input antlr.TokenStream.
func NewQuery(input antlr.TokenStream) *Query {
	QueryInit()
	this := new(Query)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &queryParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "java-escape"

	return this
}

// Query tokens.
const (
	QueryEOF             = antlr.TokenEOF
	QueryAND_OP          = 1
	QueryOR_OP           = 2
	QuerySIMPLE_OP       = 3
	QueryREP             = 4
	QueryIN              = 5
	QueryAS              = 6
	QueryCBF             = 7
	QuerySELECT          = 8
	QueryFROM            = 9
	QueryFILTER          = 10
	QueryWILDCARD        = 11
	QueryCLAUSE_SAME     = 12
	QueryCLAUSE_DISTINCT = 13
	QueryL_PAREN         = 14
	QueryR_PAREN         = 15
	QueryAT              = 16
	QueryIDENT           = 17
	QueryNUMBER1         = 18
	QueryZERO            = 19
	QuerySTRING          = 20
	QueryWS              = 21
)

// Query rules.
const (
	QueryRULE_policy      = 0
	QueryRULE_repStmt     = 1
	QueryRULE_cbfStmt     = 2
	QueryRULE_selectStmt  = 3
	QueryRULE_clause      = 4
	QueryRULE_filterExpr  = 5
	QueryRULE_filterStmt  = 6
	QueryRULE_expr        = 7
	QueryRULE_filterKey   = 8
	QueryRULE_filterValue = 9
	QueryRULE_number      = 10
	QueryRULE_keyword     = 11
	QueryRULE_ident       = 12
	QueryRULE_identWC     = 13
)

// IPolicyContext is an interface to support dynamic dispatch.
type IPolicyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPolicyContext differentiates from other interfaces.
	IsPolicyContext()
}

type PolicyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPolicyContext() *PolicyContext {
	var p = new(PolicyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = QueryRULE_policy
	return p
}

func (*PolicyContext) IsPolicyContext() {}

func NewPolicyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PolicyContext {
	var p = new(PolicyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryRULE_policy

	return p
}

func (s *PolicyContext) GetParser() antlr.Parser { return s.parser }

func (s *PolicyContext) EOF() antlr.TerminalNode {
	return s.GetToken(QueryEOF, 0)
}

func (s *PolicyContext) AllRepStmt() []IRepStmtContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRepStmtContext); ok {
			len++
		}
	}

	tst := make([]IRepStmtContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRepStmtContext); ok {
			tst[i] = t.(IRepStmtContext)
			i++
		}
	}

	return tst
}

func (s *PolicyContext) RepStmt(i int) IRepStmtContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRepStmtContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRepStmtContext)
}

func (s *PolicyContext) CbfStmt() ICbfStmtContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICbfStmtContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICbfStmtContext)
}

func (s *PolicyContext) AllSelectStmt() []ISelectStmtContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISelectStmtContext); ok {
			len++
		}
	}

	tst := make([]ISelectStmtContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISelectStmtContext); ok {
			tst[i] = t.(ISelectStmtContext)
			i++
		}
	}

	return tst
}

func (s *PolicyContext) SelectStmt(i int) ISelectStmtContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectStmtContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectStmtContext)
}

func (s *PolicyContext) AllFilterStmt() []IFilterStmtContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFilterStmtContext); ok {
			len++
		}
	}

	tst := make([]IFilterStmtContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFilterStmtContext); ok {
			tst[i] = t.(IFilterStmtContext)
			i++
		}
	}

	return tst
}

func (s *PolicyContext) FilterStmt(i int) IFilterStmtContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterStmtContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterStmtContext)
}

func (s *PolicyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PolicyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PolicyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterPolicy(s)
	}
}

func (s *PolicyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitPolicy(s)
	}
}

func (s *PolicyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitPolicy(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Query) Policy() (localctx IPolicyContext) {
	this := p
	_ = this

	localctx = NewPolicyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, QueryRULE_policy)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(29)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == QueryREP {
		{
			p.SetState(28)
			p.RepStmt()
		}

		p.SetState(31)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(34)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == QueryCBF {
		{
			p.SetState(33)
			p.CbfStmt()
		}

	}
	p.SetState(39)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == QuerySELECT {
		{
			p.SetState(36)
			p.SelectStmt()
		}

		p.SetState(41)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(45)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == QueryFILTER {
		{
			p.SetState(42)
			p.FilterStmt()
		}

		p.SetState(47)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(48)
		p.Match(QueryEOF)
	}

	return localctx
}

// IRepStmtContext is an interface to support dynamic dispatch.
type IRepStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetCount returns the Count token.
	GetCount() antlr.Token

	// SetCount sets the Count token.
	SetCount(antlr.Token)

	// GetSelector returns the Selector rule contexts.
	GetSelector() IIdentContext

	// SetSelector sets the Selector rule contexts.
	SetSelector(IIdentContext)

	// IsRepStmtContext differentiates from other interfaces.
	IsRepStmtContext()
}

type RepStmtContext struct {
	*antlr.BaseParserRuleContext
	parser   antlr.Parser
	Count    antlr.Token
	Selector IIdentContext
}

func NewEmptyRepStmtContext() *RepStmtContext {
	var p = new(RepStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = QueryRULE_repStmt
	return p
}

func (*RepStmtContext) IsRepStmtContext() {}

func NewRepStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RepStmtContext {
	var p = new(RepStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryRULE_repStmt

	return p
}

func (s *RepStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *RepStmtContext) GetCount() antlr.Token { return s.Count }

func (s *RepStmtContext) SetCount(v antlr.Token) { s.Count = v }

func (s *RepStmtContext) GetSelector() IIdentContext { return s.Selector }

func (s *RepStmtContext) SetSelector(v IIdentContext) { s.Selector = v }

func (s *RepStmtContext) REP() antlr.TerminalNode {
	return s.GetToken(QueryREP, 0)
}

func (s *RepStmtContext) NUMBER1() antlr.TerminalNode {
	return s.GetToken(QueryNUMBER1, 0)
}

func (s *RepStmtContext) IN() antlr.TerminalNode {
	return s.GetToken(QueryIN, 0)
}

func (s *RepStmtContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *RepStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RepStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RepStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterRepStmt(s)
	}
}

func (s *RepStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitRepStmt(s)
	}
}

func (s *RepStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitRepStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Query) RepStmt() (localctx IRepStmtContext) {
	this := p
	_ = this

	localctx = NewRepStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, QueryRULE_repStmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(50)
		p.Match(QueryREP)
	}
	{
		p.SetState(51)

		var _m = p.Match(QueryNUMBER1)

		localctx.(*RepStmtContext).Count = _m
	}
	p.SetState(54)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == QueryIN {
		{
			p.SetState(52)
			p.Match(QueryIN)
		}
		{
			p.SetState(53)

			var _x = p.Ident()

			localctx.(*RepStmtContext).Selector = _x
		}

	}

	return localctx
}

// ICbfStmtContext is an interface to support dynamic dispatch.
type ICbfStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetBackupFactor returns the BackupFactor token.
	GetBackupFactor() antlr.Token

	// SetBackupFactor sets the BackupFactor token.
	SetBackupFactor(antlr.Token)

	// IsCbfStmtContext differentiates from other interfaces.
	IsCbfStmtContext()
}

type CbfStmtContext struct {
	*antlr.BaseParserRuleContext
	parser       antlr.Parser
	BackupFactor antlr.Token
}

func NewEmptyCbfStmtContext() *CbfStmtContext {
	var p = new(CbfStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = QueryRULE_cbfStmt
	return p
}

func (*CbfStmtContext) IsCbfStmtContext() {}

func NewCbfStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CbfStmtContext {
	var p = new(CbfStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryRULE_cbfStmt

	return p
}

func (s *CbfStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *CbfStmtContext) GetBackupFactor() antlr.Token { return s.BackupFactor }

func (s *CbfStmtContext) SetBackupFactor(v antlr.Token) { s.BackupFactor = v }

func (s *CbfStmtContext) CBF() antlr.TerminalNode {
	return s.GetToken(QueryCBF, 0)
}

func (s *CbfStmtContext) NUMBER1() antlr.TerminalNode {
	return s.GetToken(QueryNUMBER1, 0)
}

func (s *CbfStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CbfStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CbfStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterCbfStmt(s)
	}
}

func (s *CbfStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitCbfStmt(s)
	}
}

func (s *CbfStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitCbfStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Query) CbfStmt() (localctx ICbfStmtContext) {
	this := p
	_ = this

	localctx = NewCbfStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, QueryRULE_cbfStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(56)
		p.Match(QueryCBF)
	}
	{
		p.SetState(57)

		var _m = p.Match(QueryNUMBER1)

		localctx.(*CbfStmtContext).BackupFactor = _m
	}

	return localctx
}

// ISelectStmtContext is an interface to support dynamic dispatch.
type ISelectStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetCount returns the Count token.
	GetCount() antlr.Token

	// SetCount sets the Count token.
	SetCount(antlr.Token)

	// GetBucket returns the Bucket rule contexts.
	GetBucket() IIdentContext

	// GetFilter returns the Filter rule contexts.
	GetFilter() IIdentWCContext

	// GetName returns the Name rule contexts.
	GetName() IIdentContext

	// SetBucket sets the Bucket rule contexts.
	SetBucket(IIdentContext)

	// SetFilter sets the Filter rule contexts.
	SetFilter(IIdentWCContext)

	// SetName sets the Name rule contexts.
	SetName(IIdentContext)

	// IsSelectStmtContext differentiates from other interfaces.
	IsSelectStmtContext()
}

type SelectStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	Count  antlr.Token
	Bucket IIdentContext
	Filter IIdentWCContext
	Name   IIdentContext
}

func NewEmptySelectStmtContext() *SelectStmtContext {
	var p = new(SelectStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = QueryRULE_selectStmt
	return p
}

func (*SelectStmtContext) IsSelectStmtContext() {}

func NewSelectStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectStmtContext {
	var p = new(SelectStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryRULE_selectStmt

	return p
}

func (s *SelectStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectStmtContext) GetCount() antlr.Token { return s.Count }

func (s *SelectStmtContext) SetCount(v antlr.Token) { s.Count = v }

func (s *SelectStmtContext) GetBucket() IIdentContext { return s.Bucket }

func (s *SelectStmtContext) GetFilter() IIdentWCContext { return s.Filter }

func (s *SelectStmtContext) GetName() IIdentContext { return s.Name }

func (s *SelectStmtContext) SetBucket(v IIdentContext) { s.Bucket = v }

func (s *SelectStmtContext) SetFilter(v IIdentWCContext) { s.Filter = v }

func (s *SelectStmtContext) SetName(v IIdentContext) { s.Name = v }

func (s *SelectStmtContext) SELECT() antlr.TerminalNode {
	return s.GetToken(QuerySELECT, 0)
}

func (s *SelectStmtContext) FROM() antlr.TerminalNode {
	return s.GetToken(QueryFROM, 0)
}

func (s *SelectStmtContext) NUMBER1() antlr.TerminalNode {
	return s.GetToken(QueryNUMBER1, 0)
}

func (s *SelectStmtContext) IdentWC() IIdentWCContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentWCContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentWCContext)
}

func (s *SelectStmtContext) IN() antlr.TerminalNode {
	return s.GetToken(QueryIN, 0)
}

func (s *SelectStmtContext) AS() antlr.TerminalNode {
	return s.GetToken(QueryAS, 0)
}

func (s *SelectStmtContext) AllIdent() []IIdentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IIdentContext); ok {
			len++
		}
	}

	tst := make([]IIdentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IIdentContext); ok {
			tst[i] = t.(IIdentContext)
			i++
		}
	}

	return tst
}

func (s *SelectStmtContext) Ident(i int) IIdentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *SelectStmtContext) Clause() IClauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IClauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IClauseContext)
}

func (s *SelectStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SelectStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterSelectStmt(s)
	}
}

func (s *SelectStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitSelectStmt(s)
	}
}

func (s *SelectStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitSelectStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Query) SelectStmt() (localctx ISelectStmtContext) {
	this := p
	_ = this

	localctx = NewSelectStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, QueryRULE_selectStmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(59)
		p.Match(QuerySELECT)
	}
	{
		p.SetState(60)

		var _m = p.Match(QueryNUMBER1)

		localctx.(*SelectStmtContext).Count = _m
	}
	p.SetState(66)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == QueryIN {
		{
			p.SetState(61)
			p.Match(QueryIN)
		}
		p.SetState(63)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == QueryCLAUSE_SAME || _la == QueryCLAUSE_DISTINCT {
			{
				p.SetState(62)
				p.Clause()
			}

		}
		{
			p.SetState(65)

			var _x = p.Ident()

			localctx.(*SelectStmtContext).Bucket = _x
		}

	}
	{
		p.SetState(68)
		p.Match(QueryFROM)
	}
	{
		p.SetState(69)

		var _x = p.IdentWC()

		localctx.(*SelectStmtContext).Filter = _x
	}
	p.SetState(72)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == QueryAS {
		{
			p.SetState(70)
			p.Match(QueryAS)
		}
		{
			p.SetState(71)

			var _x = p.Ident()

			localctx.(*SelectStmtContext).Name = _x
		}

	}

	return localctx
}

// IClauseContext is an interface to support dynamic dispatch.
type IClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsClauseContext differentiates from other interfaces.
	IsClauseContext()
}

type ClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyClauseContext() *ClauseContext {
	var p = new(ClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = QueryRULE_clause
	return p
}

func (*ClauseContext) IsClauseContext() {}

func NewClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClauseContext {
	var p = new(ClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryRULE_clause

	return p
}

func (s *ClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *ClauseContext) CLAUSE_SAME() antlr.TerminalNode {
	return s.GetToken(QueryCLAUSE_SAME, 0)
}

func (s *ClauseContext) CLAUSE_DISTINCT() antlr.TerminalNode {
	return s.GetToken(QueryCLAUSE_DISTINCT, 0)
}

func (s *ClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterClause(s)
	}
}

func (s *ClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitClause(s)
	}
}

func (s *ClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Query) Clause() (localctx IClauseContext) {
	this := p
	_ = this

	localctx = NewClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, QueryRULE_clause)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(74)
		_la = p.GetTokenStream().LA(1)

		if !(_la == QueryCLAUSE_SAME || _la == QueryCLAUSE_DISTINCT) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IFilterExprContext is an interface to support dynamic dispatch.
type IFilterExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetOp returns the Op token.
	GetOp() antlr.Token

	// SetOp sets the Op token.
	SetOp(antlr.Token)

	// GetF1 returns the F1 rule contexts.
	GetF1() IFilterExprContext

	// GetInner returns the Inner rule contexts.
	GetInner() IFilterExprContext

	// GetF2 returns the F2 rule contexts.
	GetF2() IFilterExprContext

	// SetF1 sets the F1 rule contexts.
	SetF1(IFilterExprContext)

	// SetInner sets the Inner rule contexts.
	SetInner(IFilterExprContext)

	// SetF2 sets the F2 rule contexts.
	SetF2(IFilterExprContext)

	// IsFilterExprContext differentiates from other interfaces.
	IsFilterExprContext()
}

type FilterExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	F1     IFilterExprContext
	Inner  IFilterExprContext
	Op     antlr.Token
	F2     IFilterExprContext
}

func NewEmptyFilterExprContext() *FilterExprContext {
	var p = new(FilterExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = QueryRULE_filterExpr
	return p
}

func (*FilterExprContext) IsFilterExprContext() {}

func NewFilterExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterExprContext {
	var p = new(FilterExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryRULE_filterExpr

	return p
}

func (s *FilterExprContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterExprContext) GetOp() antlr.Token { return s.Op }

func (s *FilterExprContext) SetOp(v antlr.Token) { s.Op = v }

func (s *FilterExprContext) GetF1() IFilterExprContext { return s.F1 }

func (s *FilterExprContext) GetInner() IFilterExprContext { return s.Inner }

func (s *FilterExprContext) GetF2() IFilterExprContext { return s.F2 }

func (s *FilterExprContext) SetF1(v IFilterExprContext) { s.F1 = v }

func (s *FilterExprContext) SetInner(v IFilterExprContext) { s.Inner = v }

func (s *FilterExprContext) SetF2(v IFilterExprContext) { s.F2 = v }

func (s *FilterExprContext) L_PAREN() antlr.TerminalNode {
	return s.GetToken(QueryL_PAREN, 0)
}

func (s *FilterExprContext) R_PAREN() antlr.TerminalNode {
	return s.GetToken(QueryR_PAREN, 0)
}

func (s *FilterExprContext) AllFilterExpr() []IFilterExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFilterExprContext); ok {
			len++
		}
	}

	tst := make([]IFilterExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFilterExprContext); ok {
			tst[i] = t.(IFilterExprContext)
			i++
		}
	}

	return tst
}

func (s *FilterExprContext) FilterExpr(i int) IFilterExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterExprContext)
}

func (s *FilterExprContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *FilterExprContext) AND_OP() antlr.TerminalNode {
	return s.GetToken(QueryAND_OP, 0)
}

func (s *FilterExprContext) OR_OP() antlr.TerminalNode {
	return s.GetToken(QueryOR_OP, 0)
}

func (s *FilterExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterFilterExpr(s)
	}
}

func (s *FilterExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitFilterExpr(s)
	}
}

func (s *FilterExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitFilterExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Query) FilterExpr() (localctx IFilterExprContext) {
	return p.filterExpr(0)
}

func (p *Query) filterExpr(_p int) (localctx IFilterExprContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewFilterExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IFilterExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 10
	p.EnterRecursionRule(localctx, 10, QueryRULE_filterExpr, _p)

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(82)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case QueryL_PAREN:
		{
			p.SetState(77)
			p.Match(QueryL_PAREN)
		}
		{
			p.SetState(78)

			var _x = p.filterExpr(0)

			localctx.(*FilterExprContext).Inner = _x
		}
		{
			p.SetState(79)
			p.Match(QueryR_PAREN)
		}

	case QueryREP, QueryIN, QueryAS, QuerySELECT, QueryFROM, QueryFILTER, QueryAT, QueryIDENT, QuerySTRING:
		{
			p.SetState(81)
			p.Expr()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(92)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(90)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext()) {
			case 1:
				localctx = NewFilterExprContext(p, _parentctx, _parentState)
				localctx.(*FilterExprContext).F1 = _prevctx
				p.PushNewRecursionContext(localctx, _startState, QueryRULE_filterExpr)
				p.SetState(84)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
				}
				{
					p.SetState(85)

					var _m = p.Match(QueryAND_OP)

					localctx.(*FilterExprContext).Op = _m
				}
				{
					p.SetState(86)

					var _x = p.filterExpr(5)

					localctx.(*FilterExprContext).F2 = _x
				}

			case 2:
				localctx = NewFilterExprContext(p, _parentctx, _parentState)
				localctx.(*FilterExprContext).F1 = _prevctx
				p.PushNewRecursionContext(localctx, _startState, QueryRULE_filterExpr)
				p.SetState(87)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
				}
				{
					p.SetState(88)

					var _m = p.Match(QueryOR_OP)

					localctx.(*FilterExprContext).Op = _m
				}
				{
					p.SetState(89)

					var _x = p.filterExpr(4)

					localctx.(*FilterExprContext).F2 = _x
				}

			}

		}
		p.SetState(94)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext())
	}

	return localctx
}

// IFilterStmtContext is an interface to support dynamic dispatch.
type IFilterStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetExpr returns the Expr rule contexts.
	GetExpr() IFilterExprContext

	// GetName returns the Name rule contexts.
	GetName() IIdentContext

	// SetExpr sets the Expr rule contexts.
	SetExpr(IFilterExprContext)

	// SetName sets the Name rule contexts.
	SetName(IIdentContext)

	// IsFilterStmtContext differentiates from other interfaces.
	IsFilterStmtContext()
}

type FilterStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	Expr   IFilterExprContext
	Name   IIdentContext
}

func NewEmptyFilterStmtContext() *FilterStmtContext {
	var p = new(FilterStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = QueryRULE_filterStmt
	return p
}

func (*FilterStmtContext) IsFilterStmtContext() {}

func NewFilterStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterStmtContext {
	var p = new(FilterStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryRULE_filterStmt

	return p
}

func (s *FilterStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterStmtContext) GetExpr() IFilterExprContext { return s.Expr }

func (s *FilterStmtContext) GetName() IIdentContext { return s.Name }

func (s *FilterStmtContext) SetExpr(v IFilterExprContext) { s.Expr = v }

func (s *FilterStmtContext) SetName(v IIdentContext) { s.Name = v }

func (s *FilterStmtContext) FILTER() antlr.TerminalNode {
	return s.GetToken(QueryFILTER, 0)
}

func (s *FilterStmtContext) AS() antlr.TerminalNode {
	return s.GetToken(QueryAS, 0)
}

func (s *FilterStmtContext) FilterExpr() IFilterExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterExprContext)
}

func (s *FilterStmtContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *FilterStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterFilterStmt(s)
	}
}

func (s *FilterStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitFilterStmt(s)
	}
}

func (s *FilterStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitFilterStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Query) FilterStmt() (localctx IFilterStmtContext) {
	this := p
	_ = this

	localctx = NewFilterStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, QueryRULE_filterStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(95)
		p.Match(QueryFILTER)
	}
	{
		p.SetState(96)

		var _x = p.filterExpr(0)

		localctx.(*FilterStmtContext).Expr = _x
	}
	{
		p.SetState(97)
		p.Match(QueryAS)
	}
	{
		p.SetState(98)

		var _x = p.Ident()

		localctx.(*FilterStmtContext).Name = _x
	}

	return localctx
}

// IExprContext is an interface to support dynamic dispatch.
type IExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetFilter returns the Filter rule contexts.
	GetFilter() IIdentContext

	// GetKey returns the Key rule contexts.
	GetKey() IFilterKeyContext

	// GetValue returns the Value rule contexts.
	GetValue() IFilterValueContext

	// SetFilter sets the Filter rule contexts.
	SetFilter(IIdentContext)

	// SetKey sets the Key rule contexts.
	SetKey(IFilterKeyContext)

	// SetValue sets the Value rule contexts.
	SetValue(IFilterValueContext)

	// IsExprContext differentiates from other interfaces.
	IsExprContext()
}

type ExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	Filter IIdentContext
	Key    IFilterKeyContext
	Value  IFilterValueContext
}

func NewEmptyExprContext() *ExprContext {
	var p = new(ExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = QueryRULE_expr
	return p
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryRULE_expr

	return p
}

func (s *ExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprContext) GetFilter() IIdentContext { return s.Filter }

func (s *ExprContext) GetKey() IFilterKeyContext { return s.Key }

func (s *ExprContext) GetValue() IFilterValueContext { return s.Value }

func (s *ExprContext) SetFilter(v IIdentContext) { s.Filter = v }

func (s *ExprContext) SetKey(v IFilterKeyContext) { s.Key = v }

func (s *ExprContext) SetValue(v IFilterValueContext) { s.Value = v }

func (s *ExprContext) AT() antlr.TerminalNode {
	return s.GetToken(QueryAT, 0)
}

func (s *ExprContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *ExprContext) SIMPLE_OP() antlr.TerminalNode {
	return s.GetToken(QuerySIMPLE_OP, 0)
}

func (s *ExprContext) FilterKey() IFilterKeyContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterKeyContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterKeyContext)
}

func (s *ExprContext) FilterValue() IFilterValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterValueContext)
}

func (s *ExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterExpr(s)
	}
}

func (s *ExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitExpr(s)
	}
}

func (s *ExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Query) Expr() (localctx IExprContext) {
	this := p
	_ = this

	localctx = NewExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, QueryRULE_expr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(106)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case QueryAT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(100)
			p.Match(QueryAT)
		}
		{
			p.SetState(101)

			var _x = p.Ident()

			localctx.(*ExprContext).Filter = _x
		}

	case QueryREP, QueryIN, QueryAS, QuerySELECT, QueryFROM, QueryFILTER, QueryIDENT, QuerySTRING:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(102)

			var _x = p.FilterKey()

			localctx.(*ExprContext).Key = _x
		}
		{
			p.SetState(103)
			p.Match(QuerySIMPLE_OP)
		}
		{
			p.SetState(104)

			var _x = p.FilterValue()

			localctx.(*ExprContext).Value = _x
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IFilterKeyContext is an interface to support dynamic dispatch.
type IFilterKeyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFilterKeyContext differentiates from other interfaces.
	IsFilterKeyContext()
}

type FilterKeyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterKeyContext() *FilterKeyContext {
	var p = new(FilterKeyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = QueryRULE_filterKey
	return p
}

func (*FilterKeyContext) IsFilterKeyContext() {}

func NewFilterKeyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterKeyContext {
	var p = new(FilterKeyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryRULE_filterKey

	return p
}

func (s *FilterKeyContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterKeyContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *FilterKeyContext) STRING() antlr.TerminalNode {
	return s.GetToken(QuerySTRING, 0)
}

func (s *FilterKeyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterKeyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterKeyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterFilterKey(s)
	}
}

func (s *FilterKeyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitFilterKey(s)
	}
}

func (s *FilterKeyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitFilterKey(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Query) FilterKey() (localctx IFilterKeyContext) {
	this := p
	_ = this

	localctx = NewFilterKeyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, QueryRULE_filterKey)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(110)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case QueryREP, QueryIN, QueryAS, QuerySELECT, QueryFROM, QueryFILTER, QueryIDENT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(108)
			p.Ident()
		}

	case QuerySTRING:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(109)
			p.Match(QuerySTRING)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IFilterValueContext is an interface to support dynamic dispatch.
type IFilterValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFilterValueContext differentiates from other interfaces.
	IsFilterValueContext()
}

type FilterValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterValueContext() *FilterValueContext {
	var p = new(FilterValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = QueryRULE_filterValue
	return p
}

func (*FilterValueContext) IsFilterValueContext() {}

func NewFilterValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterValueContext {
	var p = new(FilterValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryRULE_filterValue

	return p
}

func (s *FilterValueContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterValueContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *FilterValueContext) Number() INumberContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INumberContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INumberContext)
}

func (s *FilterValueContext) STRING() antlr.TerminalNode {
	return s.GetToken(QuerySTRING, 0)
}

func (s *FilterValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterFilterValue(s)
	}
}

func (s *FilterValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitFilterValue(s)
	}
}

func (s *FilterValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitFilterValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Query) FilterValue() (localctx IFilterValueContext) {
	this := p
	_ = this

	localctx = NewFilterValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, QueryRULE_filterValue)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(115)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case QueryREP, QueryIN, QueryAS, QuerySELECT, QueryFROM, QueryFILTER, QueryIDENT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(112)
			p.Ident()
		}

	case QueryNUMBER1, QueryZERO:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(113)
			p.Number()
		}

	case QuerySTRING:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(114)
			p.Match(QuerySTRING)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// INumberContext is an interface to support dynamic dispatch.
type INumberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNumberContext differentiates from other interfaces.
	IsNumberContext()
}

type NumberContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNumberContext() *NumberContext {
	var p = new(NumberContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = QueryRULE_number
	return p
}

func (*NumberContext) IsNumberContext() {}

func NewNumberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumberContext {
	var p = new(NumberContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryRULE_number

	return p
}

func (s *NumberContext) GetParser() antlr.Parser { return s.parser }

func (s *NumberContext) ZERO() antlr.TerminalNode {
	return s.GetToken(QueryZERO, 0)
}

func (s *NumberContext) NUMBER1() antlr.TerminalNode {
	return s.GetToken(QueryNUMBER1, 0)
}

func (s *NumberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NumberContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterNumber(s)
	}
}

func (s *NumberContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitNumber(s)
	}
}

func (s *NumberContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitNumber(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Query) Number() (localctx INumberContext) {
	this := p
	_ = this

	localctx = NewNumberContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, QueryRULE_number)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(117)
		_la = p.GetTokenStream().LA(1)

		if !(_la == QueryNUMBER1 || _la == QueryZERO) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IKeywordContext is an interface to support dynamic dispatch.
type IKeywordContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsKeywordContext differentiates from other interfaces.
	IsKeywordContext()
}

type KeywordContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyKeywordContext() *KeywordContext {
	var p = new(KeywordContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = QueryRULE_keyword
	return p
}

func (*KeywordContext) IsKeywordContext() {}

func NewKeywordContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *KeywordContext {
	var p = new(KeywordContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryRULE_keyword

	return p
}

func (s *KeywordContext) GetParser() antlr.Parser { return s.parser }

func (s *KeywordContext) REP() antlr.TerminalNode {
	return s.GetToken(QueryREP, 0)
}

func (s *KeywordContext) IN() antlr.TerminalNode {
	return s.GetToken(QueryIN, 0)
}

func (s *KeywordContext) AS() antlr.TerminalNode {
	return s.GetToken(QueryAS, 0)
}

func (s *KeywordContext) SELECT() antlr.TerminalNode {
	return s.GetToken(QuerySELECT, 0)
}

func (s *KeywordContext) FROM() antlr.TerminalNode {
	return s.GetToken(QueryFROM, 0)
}

func (s *KeywordContext) FILTER() antlr.TerminalNode {
	return s.GetToken(QueryFILTER, 0)
}

func (s *KeywordContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *KeywordContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *KeywordContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterKeyword(s)
	}
}

func (s *KeywordContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitKeyword(s)
	}
}

func (s *KeywordContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitKeyword(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Query) Keyword() (localctx IKeywordContext) {
	this := p
	_ = this

	localctx = NewKeywordContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, QueryRULE_keyword)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(119)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1904) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IIdentContext is an interface to support dynamic dispatch.
type IIdentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIdentContext differentiates from other interfaces.
	IsIdentContext()
}

type IdentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdentContext() *IdentContext {
	var p = new(IdentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = QueryRULE_ident
	return p
}

func (*IdentContext) IsIdentContext() {}

func NewIdentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentContext {
	var p = new(IdentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryRULE_ident

	return p
}

func (s *IdentContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentContext) Keyword() IKeywordContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IKeywordContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IKeywordContext)
}

func (s *IdentContext) IDENT() antlr.TerminalNode {
	return s.GetToken(QueryIDENT, 0)
}

func (s *IdentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IdentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterIdent(s)
	}
}

func (s *IdentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitIdent(s)
	}
}

func (s *IdentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitIdent(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Query) Ident() (localctx IIdentContext) {
	this := p
	_ = this

	localctx = NewIdentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, QueryRULE_ident)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(123)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case QueryREP, QueryIN, QueryAS, QuerySELECT, QueryFROM, QueryFILTER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(121)
			p.Keyword()
		}

	case QueryIDENT:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(122)
			p.Match(QueryIDENT)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IIdentWCContext is an interface to support dynamic dispatch.
type IIdentWCContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIdentWCContext differentiates from other interfaces.
	IsIdentWCContext()
}

type IdentWCContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdentWCContext() *IdentWCContext {
	var p = new(IdentWCContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = QueryRULE_identWC
	return p
}

func (*IdentWCContext) IsIdentWCContext() {}

func NewIdentWCContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentWCContext {
	var p = new(IdentWCContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryRULE_identWC

	return p
}

func (s *IdentWCContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentWCContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *IdentWCContext) WILDCARD() antlr.TerminalNode {
	return s.GetToken(QueryWILDCARD, 0)
}

func (s *IdentWCContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentWCContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IdentWCContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterIdentWC(s)
	}
}

func (s *IdentWCContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitIdentWC(s)
	}
}

func (s *IdentWCContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitIdentWC(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Query) IdentWC() (localctx IIdentWCContext) {
	this := p
	_ = this

	localctx = NewIdentWCContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, QueryRULE_identWC)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(127)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case QueryREP, QueryIN, QueryAS, QuerySELECT, QueryFROM, QueryFILTER, QueryIDENT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(125)
			p.Ident()
		}

	case QueryWILDCARD:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(126)
			p.Match(QueryWILDCARD)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

func (p *Query) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 5:
		var t *FilterExprContext = nil
		if localctx != nil {
			t = localctx.(*FilterExprContext)
		}
		return p.FilterExpr_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *Query) FilterExpr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 4)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 3)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
