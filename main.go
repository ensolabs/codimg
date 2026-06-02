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

	lines := strings.Split(input, "\n")

	maxLen := 0
	for _, line := range lines {
		if n := utf8.RuneCountInString(line); n > maxLen {
			maxLen = n
		}
	}

	width := int(float64(maxLen)*charWidth) + paddingX*2
	height := len(lines)*lineHeight + paddingY*2

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	fmt.Fprintf(w, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`, width, height)
	fmt.Fprintf(w, `<rect width="%d" height="%d" rx="8" fill="black"/>`, width, height)
	fmt.Fprintf(w, `<text xml:space="preserve" font-family="monospace" font-size="%d" fill="white">`, fontSize)

	for i, line := range lines {
		y := paddingY + i*lineHeight + (lineHeight+fontSize)/2
		fmt.Fprintf(w, `<tspan x="%d" y="%d">%s</tspan>`, paddingX, y, svgEscaper.Replace(line))
	}

	fmt.Fprint(w, `</text></svg>`)
}

func main() {
	http.HandleFunc("/code.svg", codeHandler)
	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
