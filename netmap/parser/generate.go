package parser

// ANTLR can be downloaded from https://www.antlr.org/download/antlr-4.11.1-complete.jar
//go:generate java -Xmx500M -cp "./antlr-4.11.1-complete.jar:$CLASSPATH" org.antlr.v4.Tool -Dlanguage=Go -visitor QueryLexer.g4 Query.g4
