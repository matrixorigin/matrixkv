ROOT_DIR = $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/
LD_FLAGS = -ldflags "-w -s"

GOOS 		= linux
DIST_DIR 	= $(ROOT_DIR)dist/

.PHONY: dist_dir
dist_dir: ; $(info ======== prepare distribute dir:)
	mkdir -p $(DIST_DIR)
	@rm -rf $(DIST_DIR)tinykv

.PHONY: tinykv
tinykv: dist_dir; $(info ======== compiled tinykv binary)
	env GOOS=$(GOOS) go build -o $(DIST_DIR)tinykv $(LD_FLAGS) $(ROOT_DIR)cmd/*.go

.PHONY: docker
docker: ; $(info ======== compiled tinykv docker)
	docker build -t tinykv -f Dockerfile .

.DEFAULT_GOAL := tinykv