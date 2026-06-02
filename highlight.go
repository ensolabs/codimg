package main

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

const (
	fallbackColor = "#ffffff"
	themeName     = "gruvbox"
)

// Token is a single styled run of text within a line.
type Token struct {
	Text  string
	Color string // hex, e.g. "#f92672"
}

// Tokenize splits source code into lines of colored tokens.
// lang is a language name (e.g. "go", "python"); pass "" to auto-detect.
// Falls back gracefully: unknown language → plain text, unset color → white.
func Tokenize(source, lang string) [][]Token {
	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Analyse(source)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get(themeName)
	if style == nil {
		style = styles.Fallback
	}

	iterator, err := lexer.Tokenise(nil, source)
	if err != nil {
		return plainLines(source)
	}

	chromaLines := chroma.SplitTokensIntoLines(iterator.Tokens())
	result := make([][]Token, len(chromaLines))
	for i, chromaLine := range chromaLines {
		result[i] = make([]Token, 0, len(chromaLine))
		for _, ct := range chromaLine {
			result[i] = append(result[i], Token{
				Text:  ct.Value,
				Color: tokenColor(style, ct.Type),
			})
		}
	}
	return result
}

// tokenColor resolves the foreground hex color for a token type from the style.
// Falls back to fallbackColor if the style entry has no foreground set.
func tokenColor(style *chroma.Style, tt chroma.TokenType) string {
	entry := style.Get(tt)
	if entry.Colour.IsSet() {
		return entry.Colour.String()
	}
	return fallbackColor
}

// plainLines wraps raw source into [][]Token with no color (fallback white).
func plainLines(source string) [][]Token {
	lines := splitLines(source)
	result := make([][]Token, len(lines))
	for i, line := range lines {
		result[i] = []Token{{Text: line, Color: fallbackColor}}
	}
	return result
}

// splitLines splits source on newlines, preserving empty lines but dropping a
// single trailing empty line to match Chroma's SplitTokensIntoLines behaviour.
func splitLines(source string) []string {
	out := []string{}
	start := 0
	for i := 0; i < len(source); i++ {
		if source[i] == '\n' {
			out = append(out, source[start:i])
			start = i + 1
		}
	}
	tail := source[start:]
	if len(out) > 0 || tail != "" {
		out = append(out, tail)
	}
	if len(out) > 0 && out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}
	return out
}
