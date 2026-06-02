package main

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"io"
)

func decodeInput(s string) string {
	if s == "" {
		return s
	}

	raw, ok := tryBase64(s)
	if !ok {
		return s
	}

	out, ok := tryInflate(raw)
	if !ok {
		return s
	}
	return string(out)
}

func tryBase64(s string) ([]byte, bool) {
	if b, err := base64.StdEncoding.DecodeString(s); err == nil {
		return b, true
	}
	return nil, false
}

func tryInflate(raw []byte) ([]byte, bool) {
	fr := flate.NewReader(bytes.NewReader(raw))
	defer fr.Close()
	if out, err := io.ReadAll(fr); err == nil {
		return out, true
	}
	return nil, false
}
