.DEFAULT_GOAL := help
SHELL:=/bin/bash
.ONESHELL:
.PHONY: new prep print-mp3-tags upload-mp3 deploy build

TMPFILE := $(shell mktemp)
makefile_dir := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
notif := "/usr/local/bin/terminal-notifier"
subl := "/Applications/Sublime Text.app/Contents/SharedSupport/bin/subl"

B=$(shell git rev-parse --abbrev-ref HEAD)
BRANCH=$(subst /,-,$(B))
GITREV=$(shell git describe --abbrev=7 --always --tags)
REV=$(GITREV)-$(BRANCH)-$(shell date +%Y%m%d-%H:%M:%S)

help:
	@echo 'Available commands: new, prep, print-mp3-tags, upload-mp3, deploy'


# generate new episode post markdown file and open it using SublimeText 
new:
	@docker-compose run --rm publisher new | tee ${TMPFILE};
	@${subl} ${makefile_dir}/../hugo/`tail -n 1 ${TMPFILE} | tr -d '\r'`;

# generate new prep post markdown file and open it using SublimeText 
prep:
	@docker-compose run --rm publisher prep | tee ${TMPFILE};
	@${subl} ${makefile_dir}/../hugo/`tail -n 1 ${TMPFILE} | tr -d '\r'`;

# print mp3 tags from file
print-mp3-tags:
	@docker-compose run --rm --entrypoint mp3tags publisher print-tags "$$EPISODE";

# set necessary tags to episode mp3 file and upload it, notifying via MacOS local notifications on every step
upload-mp3:
	@docker-compose run --rm publisher upload --episode="$$EPISODE" --dbg

# deploy new podcast episode page to https://radio-t.com
deploy:
	@docker-compose run --rm publisher deploy

build:
	@echo $(REV)
	@docker build -t publisher --build-arg REV=$(REV) .