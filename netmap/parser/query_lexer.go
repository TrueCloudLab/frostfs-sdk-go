// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type QueryLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var querylexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	channelNames           []string
	modeNames              []string
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func querylexerLexerInit() {
	staticData := &querylexerLexerStaticData
	staticData.channelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.modeNames = []string{
		"DEFAULT_MODE",
	}
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
		"AND_OP", "OR_OP", "SIMPLE_OP", "REP", "IN", "AS", "CBF", "SELECT",
		"FROM", "FILTER", "WILDCARD", "CLAUSE_SAME", "CLAUSE_DISTINCT", "L_PAREN",
		"R_PAREN", "AT", "IDENT", "Digit", "Nondigit", "NUMBER1", "ZERO", "STRING",
		"ESC", "UNICODE", "HEX", "SAFECODEPOINTSINGLE", "SAFECODEPOINTDOUBLE",
		"WS",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 21, 198, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7,
		20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25,
		2, 26, 7, 26, 2, 27, 7, 27, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1,
		2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 3,
		2, 77, 8, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5,
		1, 6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8,
		1, 8, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10,
		1, 10, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 12, 1, 12, 1, 12, 1, 12, 1,
		12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1, 14, 1, 14, 1, 15, 1, 15,
		1, 16, 1, 16, 1, 16, 5, 16, 137, 8, 16, 10, 16, 12, 16, 140, 9, 16, 1,
		17, 1, 17, 1, 18, 1, 18, 1, 19, 1, 19, 5, 19, 148, 8, 19, 10, 19, 12, 19,
		151, 9, 19, 1, 20, 1, 20, 1, 21, 1, 21, 1, 21, 5, 21, 158, 8, 21, 10, 21,
		12, 21, 161, 9, 21, 1, 21, 1, 21, 1, 21, 1, 21, 5, 21, 167, 8, 21, 10,
		21, 12, 21, 170, 9, 21, 1, 21, 3, 21, 173, 8, 21, 1, 22, 1, 22, 1, 22,
		3, 22, 178, 8, 22, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 24, 1,
		24, 1, 25, 1, 25, 1, 26, 1, 26, 1, 27, 4, 27, 193, 8, 27, 11, 27, 12, 27,
		194, 1, 27, 1, 27, 0, 0, 28, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7,
		15, 8, 17, 9, 19, 10, 21, 11, 23, 12, 25, 13, 27, 14, 29, 15, 31, 16, 33,
		17, 35, 0, 37, 0, 39, 18, 41, 19, 43, 20, 45, 0, 47, 0, 49, 0, 51, 0, 53,
		0, 55, 21, 1, 0, 8, 1, 0, 48, 57, 3, 0, 65, 90, 95, 95, 97, 122, 1, 0,
		49, 57, 9, 0, 34, 34, 39, 39, 47, 47, 92, 92, 98, 98, 102, 102, 110, 110,
		114, 114, 116, 116, 3, 0, 48, 57, 65, 70, 97, 102, 3, 0, 0, 31, 39, 39,
		92, 92, 3, 0, 0, 31, 34, 34, 92, 92, 3, 0, 9, 10, 13, 13, 32, 32, 205,
		0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0,
		0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0,
		0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0,
		0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 1, 0, 0, 0, 0, 29, 1, 0, 0, 0, 0, 31, 1,
		0, 0, 0, 0, 33, 1, 0, 0, 0, 0, 39, 1, 0, 0, 0, 0, 41, 1, 0, 0, 0, 0, 43,
		1, 0, 0, 0, 0, 55, 1, 0, 0, 0, 1, 57, 1, 0, 0, 0, 3, 61, 1, 0, 0, 0, 5,
		76, 1, 0, 0, 0, 7, 78, 1, 0, 0, 0, 9, 82, 1, 0, 0, 0, 11, 85, 1, 0, 0,
		0, 13, 88, 1, 0, 0, 0, 15, 92, 1, 0, 0, 0, 17, 99, 1, 0, 0, 0, 19, 104,
		1, 0, 0, 0, 21, 111, 1, 0, 0, 0, 23, 113, 1, 0, 0, 0, 25, 118, 1, 0, 0,
		0, 27, 127, 1, 0, 0, 0, 29, 129, 1, 0, 0, 0, 31, 131, 1, 0, 0, 0, 33, 133,
		1, 0, 0, 0, 35, 141, 1, 0, 0, 0, 37, 143, 1, 0, 0, 0, 39, 145, 1, 0, 0,
		0, 41, 152, 1, 0, 0, 0, 43, 172, 1, 0, 0, 0, 45, 174, 1, 0, 0, 0, 47, 179,
		1, 0, 0, 0, 49, 185, 1, 0, 0, 0, 51, 187, 1, 0, 0, 0, 53, 189, 1, 0, 0,
		0, 55, 192, 1, 0, 0, 0, 57, 58, 5, 65, 0, 0, 58, 59, 5, 78, 0, 0, 59, 60,
		5, 68, 0, 0, 60, 2, 1, 0, 0, 0, 61, 62, 5, 79, 0, 0, 62, 63, 5, 82, 0,
		0, 63, 4, 1, 0, 0, 0, 64, 65, 5, 69, 0, 0, 65, 77, 5, 81, 0, 0, 66, 67,
		5, 78, 0, 0, 67, 77, 5, 69, 0, 0, 68, 69, 5, 71, 0, 0, 69, 77, 5, 69, 0,
		0, 70, 71, 5, 71, 0, 0, 71, 77, 5, 84, 0, 0, 72, 73, 5, 76, 0, 0, 73, 77,
		5, 84, 0, 0, 74, 75, 5, 76, 0, 0, 75, 77, 5, 69, 0, 0, 76, 64, 1, 0, 0,
		0, 76, 66, 1, 0, 0, 0, 76, 68, 1, 0, 0, 0, 76, 70, 1, 0, 0, 0, 76, 72,
		1, 0, 0, 0, 76, 74, 1, 0, 0, 0, 77, 6, 1, 0, 0, 0, 78, 79, 5, 82, 0, 0,
		79, 80, 5, 69, 0, 0, 80, 81, 5, 80, 0, 0, 81, 8, 1, 0, 0, 0, 82, 83, 5,
		73, 0, 0, 83, 84, 5, 78, 0, 0, 84, 10, 1, 0, 0, 0, 85, 86, 5, 65, 0, 0,
		86, 87, 5, 83, 0, 0, 87, 12, 1, 0, 0, 0, 88, 89, 5, 67, 0, 0, 89, 90, 5,
		66, 0, 0, 90, 91, 5, 70, 0, 0, 91, 14, 1, 0, 0, 0, 92, 93, 5, 83, 0, 0,
		93, 94, 5, 69, 0, 0, 94, 95, 5, 76, 0, 0, 95, 96, 5, 69, 0, 0, 96, 97,
		5, 67, 0, 0, 97, 98, 5, 84, 0, 0, 98, 16, 1, 0, 0, 0, 99, 100, 5, 70, 0,
		0, 100, 101, 5, 82, 0, 0, 101, 102, 5, 79, 0, 0, 102, 103, 5, 77, 0, 0,
		103, 18, 1, 0, 0, 0, 104, 105, 5, 70, 0, 0, 105, 106, 5, 73, 0, 0, 106,
		107, 5, 76, 0, 0, 107, 108, 5, 84, 0, 0, 108, 109, 5, 69, 0, 0, 109, 110,
		5, 82, 0, 0, 110, 20, 1, 0, 0, 0, 111, 112, 5, 42, 0, 0, 112, 22, 1, 0,
		0, 0, 113, 114, 5, 83, 0, 0, 114, 115, 5, 65, 0, 0, 115, 116, 5, 77, 0,
		0, 116, 117, 5, 69, 0, 0, 117, 24, 1, 0, 0, 0, 118, 119, 5, 68, 0, 0, 119,
		120, 5, 73, 0, 0, 120, 121, 5, 83, 0, 0, 121, 122, 5, 84, 0, 0, 122, 123,
		5, 73, 0, 0, 123, 124, 5, 78, 0, 0, 124, 125, 5, 67, 0, 0, 125, 126, 5,
		84, 0, 0, 126, 26, 1, 0, 0, 0, 127, 128, 5, 40, 0, 0, 128, 28, 1, 0, 0,
		0, 129, 130, 5, 41, 0, 0, 130, 30, 1, 0, 0, 0, 131, 132, 5, 64, 0, 0, 132,
		32, 1, 0, 0, 0, 133, 138, 3, 37, 18, 0, 134, 137, 3, 35, 17, 0, 135, 137,
		3, 37, 18, 0, 136, 134, 1, 0, 0, 0, 136, 135, 1, 0, 0, 0, 137, 140, 1,
		0, 0, 0, 138, 136, 1, 0, 0, 0, 138, 139, 1, 0, 0, 0, 139, 34, 1, 0, 0,
		0, 140, 138, 1, 0, 0, 0, 141, 142, 7, 0, 0, 0, 142, 36, 1, 0, 0, 0, 143,
		144, 7, 1, 0, 0, 144, 38, 1, 0, 0, 0, 145, 149, 7, 2, 0, 0, 146, 148, 3,
		35, 17, 0, 147, 146, 1, 0, 0, 0, 148, 151, 1, 0, 0, 0, 149, 147, 1, 0,
		0, 0, 149, 150, 1, 0, 0, 0, 150, 40, 1, 0, 0, 0, 151, 149, 1, 0, 0, 0,
		152, 153, 5, 48, 0, 0, 153, 42, 1, 0, 0, 0, 154, 159, 5, 34, 0, 0, 155,
		158, 3, 45, 22, 0, 156, 158, 3, 53, 26, 0, 157, 155, 1, 0, 0, 0, 157, 156,
		1, 0, 0, 0, 158, 161, 1, 0, 0, 0, 159, 157, 1, 0, 0, 0, 159, 160, 1, 0,
		0, 0, 160, 162, 1, 0, 0, 0, 161, 159, 1, 0, 0, 0, 162, 173, 5, 34, 0, 0,
		163, 168, 5, 39, 0, 0, 164, 167, 3, 45, 22, 0, 165, 167, 3, 51, 25, 0,
		166, 164, 1, 0, 0, 0, 166, 165, 1, 0, 0, 0, 167, 170, 1, 0, 0, 0, 168,
		166, 1, 0, 0, 0, 168, 169, 1, 0, 0, 0, 169, 171, 1, 0, 0, 0, 170, 168,
		1, 0, 0, 0, 171, 173, 5, 39, 0, 0, 172, 154, 1, 0, 0, 0, 172, 163, 1, 0,
		0, 0, 173, 44, 1, 0, 0, 0, 174, 177, 5, 92, 0, 0, 175, 178, 7, 3, 0, 0,
		176, 178, 3, 47, 23, 0, 177, 175, 1, 0, 0, 0, 177, 176, 1, 0, 0, 0, 178,
		46, 1, 0, 0, 0, 179, 180, 5, 117, 0, 0, 180, 181, 3, 49, 24, 0, 181, 182,
		3, 49, 24, 0, 182, 183, 3, 49, 24, 0, 183, 184, 3, 49, 24, 0, 184, 48,
		1, 0, 0, 0, 185, 186, 7, 4, 0, 0, 186, 50, 1, 0, 0, 0, 187, 188, 8, 5,
		0, 0, 188, 52, 1, 0, 0, 0, 189, 190, 8, 6, 0, 0, 190, 54, 1, 0, 0, 0, 191,
		193, 7, 7, 0, 0, 192, 191, 1, 0, 0, 0, 193, 194, 1, 0, 0, 0, 194, 192,
		1, 0, 0, 0, 194, 195, 1, 0, 0, 0, 195, 196, 1, 0, 0, 0, 196, 197, 6, 27,
		0, 0, 197, 56, 1, 0, 0, 0, 12, 0, 76, 136, 138, 149, 157, 159, 166, 168,
		172, 177, 194, 1, 6, 0, 0,
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

// QueryLexerInit initializes any static state used to implement QueryLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewQueryLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func QueryLexerInit() {
	staticData := &querylexerLexerStaticData
	staticData.once.Do(querylexerLexerInit)
}

// NewQueryLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewQueryLexer(input antlr.CharStream) *QueryLexer {
	QueryLexerInit()
	l := new(QueryLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &querylexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	l.channelNames = staticData.channelNames
	l.modeNames = staticData.modeNames
	l.RuleNames = staticData.ruleNames
	l.LiteralNames = staticData.literalNames
	l.SymbolicNames = staticData.symbolicNames
	l.GrammarFileName = "QueryLexer.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// QueryLexer tokens.
const (
	QueryLexerAND_OP          = 1
	QueryLexerOR_OP           = 2
	QueryLexerSIMPLE_OP       = 3
	QueryLexerREP             = 4
	QueryLexerIN              = 5
	QueryLexerAS              = 6
	QueryLexerCBF             = 7
	QueryLexerSELECT          = 8
	QueryLexerFROM            = 9
	QueryLexerFILTER          = 10
	QueryLexerWILDCARD        = 11
	QueryLexerCLAUSE_SAME     = 12
	QueryLexerCLAUSE_DISTINCT = 13
	QueryLexerL_PAREN         = 14
	QueryLexerR_PAREN         = 15
	QueryLexerAT              = 16
	QueryLexerIDENT           = 17
	QueryLexerNUMBER1         = 18
	QueryLexerZERO            = 19
	QueryLexerSTRING          = 20
	QueryLexerWS              = 21
)
