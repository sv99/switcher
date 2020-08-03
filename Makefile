# based on https://habr.com/ru/post/461467/
#
# for watch using watchexec from brew - github.com/watchexec/watchexec
#
.PHONY: all clean data_image help run
.DEFAULT_GOAL := help

PROJECT_NAME=$(shell basename "$(PWD)")
CWD = $(shell pwd)
SERVICE := service
PID="/tmp/.$(PROJECT_NAME).pid"

## switcher: Build binary
$(PROJECT_NAME): index.go
	@-go build -i -o bin/$@ ./cmd/$@/main.go
	@echo end-build $@

## service: Build windows service
$(SERVICE):
	GOOS=windows GOARCH=386 go build -i -o bin/$(PROJECT_NAME)_$@.exe ./cmd/$@
	GOOS=windows GOARCH=amd64 go build -i -o bin/$(PROJECT_NAME)_$@_amd64.exe ./cmd/$@

## clean: Clean build cache and remove bin directory
clean:
	go clean
	go clean -testcache
	rm -rf bin

## generate assets for index file
index.go:
	@-go-bindata -pkg $(PROJECT_NAME) -o index.go -nocompress index.html

## start: Start with watch
start:
	@-bash -c "trap '$(MAKE) stop' EXIT; $(MAKE) watch"

stop:
	@echo stop
	@-touch $(PID)
	@-kill `cat $(PID)` 2> /dev/null || true
	@-rm $(PID)
	@sleep 1

run: stop
	@echo run
	@-$(CWD)/bin/$(PROJECT_NAME) & echo $$! > $(PID)

watch:
	@echo watch
	@-watchexec --exts go \
		-w cmd/ -w . -i videodir/assets.go \
		"make $(PROJECT_NAME) run"

## help: Show commands.
help: Makefile
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
