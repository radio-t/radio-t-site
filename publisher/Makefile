.DEFAULT_GOAL := help
SHELL:=/bin/bash
.ONESHELL:
.PHONY: new-show-post prep-show-post proc-mp3 deploy-site build

TMPFILE := $(shell mktemp)
makefile_dir := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
subl := "/Applications/Sublime Text.app/Contents/SharedSupport/bin/subl"

B=$(shell git rev-parse --abbrev-ref HEAD)
BRANCH=$(subst /,-,$(B))
GITREV=$(shell git describe --abbrev=7 --always --tags)
REV=$(GITREV)-$(BRANCH)-$(shell date +%Y%m%d-%H:%M:%S)

help:
	@echo 'available commands: new-show-post, prep-show-post, proc-mp3, deploy-site'

# generate new episode post markdown file and open it using SublimeText 
new-episode:
	@docker-compose run --rm -it publisher new | tee ${TMPFILE};
	@${subl} ${makefile_dir}/../hugo/`tail -n 1 ${TMPFILE} | tr -d '\r'`;

# generate new prep post markdown file and open it using SublimeText 
prep-show-post:
	@docker-compose run --rm -it publisher prep | tee ${TMPFILE};
	@${subl} ${makefile_dir}/../hugo/`tail -n 1 ${TMPFILE} | tr -d '\r'`;

# show tags for mp3 file
show-tags:
	@docker-compose run --rm -it publisher tags --file="$$FILE" --dbg

proc-mp3:
	@docker-compose run --rm -it publisher proc --file="$$FILE"

# deploy new podcast episode page to https://radio-t.com and regenerate site
deploy-site:
	@docker-compose run --rm -it publisher deploy

build:
	@echo $(REV)
	@docker build -t radio-t/publisher --build-arg REV=$(REV) .