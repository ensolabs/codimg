package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
)

const (
	charWidth    = 9.6
	lineHeight   = 20
	paddingX     = 16
	paddingY     = 16
	fontSize     = 14
	maxInputSize = 64 * 1024
)

var svgEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
)

func codeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	input := r.URL.Query().Get("input")
	if len(input) > maxInputSize {
		http.Error(w, "input too large", http.StatusRequestEntityTooLarge)
		return
	}
	input = decodeInput(input)
	if len(input) > maxInputSize {
		http.Error(w, "input too large", http.StatusRequestEntityTooLarge)
		return
	}

	lang := r.URL.Query().Get("lang")
	tokenLines := Tokenize(input, lang)

	maxLen := 0
	for _, line := range strings.Split(input, "\n") {
		if n := utf8.RuneCountInString(line); n > maxLen {
			maxLen = n
		}
	}

	width := int(float64(maxLen)*charWidth) + paddingX*2
	height := len(tokenLines)*lineHeight + paddingY*2

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	fmt.Fprintf(w, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`, width, height)
	fmt.Fprintf(w, `<rect width="%d" height="%d" rx="8" fill="black"/>`, width, height)
	fmt.Fprintf(w, `<text xml:space="preserve" font-family="monospace" font-size="%d">`, fontSize)

	for i, tokens := range tokenLines {
		y := paddingY + i*lineHeight + (lineHeight+fontSize)/2
		fmt.Fprintf(w, `<tspan x="%d" y="%d">`, paddingX, y)
		for _, tok := range tokens {
			fmt.Fprintf(w, `<tspan fill="%s">%s</tspan>`, tok.Color, svgEscaper.Replace(tok.Text))
		}
		fmt.Fprint(w, `</tspan>`)
	}

	fmt.Fprint(w, `</text></svg>`)
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/code.svg", codeHandler)
	fmt.Println("Listening on :8100")
	log.Fatal(http.ListenAndServe("[::]:8100", nil))
}
