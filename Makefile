MJS_GO_EXECUTABLE ?= go

build:
	${MJS_GO_EXECUTABLE} build -o mjs main.go
