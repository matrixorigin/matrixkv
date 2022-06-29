ROOT_DIR = $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/
LD_FLAGS = -ldflags "-w -s"

GOOS 		= linux
DIST_DIR 	= $(ROOT_DIR)dist/

.PHONY: dist_dir
dist_dir: ; $(info ======== prepare distribute dir:)
	mkdir -p $(DIST_DIR)
	@rm -rf $(DIST_DIR)matrixkv

.PHONY: matrixkv
matrixkv: dist_dir; $(info ======== compiled matrixkv binary)
	env GOOS=$(GOOS) go build -mod=vendor -o $(DIST_DIR)matrixkv $(LD_FLAGS) $(ROOT_DIR)cmd/*.go

.PHONY: docker
docker: ; $(info ======== compiled matrixkv docker)
	docker build -t matrixkv -f Dockerfile .

.DEFAULT_GOAL := matrixkv