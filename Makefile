ELM := $(shell command -v elm 2>/dev/null)
ELM_CMD := $(if $(ELM),elm,npx --yes elm)

.PHONY: all elm go clean

all: elm go

elm:
	$(ELM_CMD) make src/Main.elm --optimize --output=static/main.js

go:
	GOOS=linux GOARCH=amd64 go build -o codimg .

clean:
	rm -f static/main.js codimg
