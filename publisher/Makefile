.DEFAULT_GOAL := help
SHELL:=/bin/bash
.ONESHELL:
.PHONY: new prep print-mp3-tags upload-mp3 deploy build

TMPFILE := $(shell mktemp)
makefile_dir := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
subl := "/Applications/Sublime Text.app/Contents/SharedSupport/bin/subl"

B=$(shell git rev-parse --abbrev-ref HEAD)
BRANCH=$(subst /,-,$(B))
GITREV=$(shell git describe --abbrev=7 --always --tags)
REV=$(GITREV)-$(BRANCH)-$(shell date +%Y%m%d-%H:%M:%S)

help:
	@echo 'available commands: new, prep, upload-mp3, deploy'

# generate new episode post markdown file and open it using SublimeText 
new:
	@docker-compose run --rm -it publisher new | tee ${TMPFILE};
	@${subl} ${makefile_dir}/../hugo/`tail -n 1 ${TMPFILE} | tr -d '\r'`;

# generate new prep post markdown file and open it using SublimeText 
prep:
	@docker-compose run --rm -it publisher prep | tee ${TMPFILE};
	@${subl} ${makefile_dir}/../hugo/`tail -n 1 ${TMPFILE} | tr -d '\r'`;

# set necessary tags to episode mp3 file and upload it to master and all nodes
proc:
	@docker-compose run --rm -it publisher proc --episode="$$EPISODE" --dbg

# deploy new podcast episode page to https://radio-t.com and regenerate site
deploy:
	@docker-compose run --rm -it publisher deploy

build:
	@echo $(REV)
	@docker build -t publisher --build-arg REV=$(REV) .