.DEFAULT_GOAL := list
SHELL:=/bin/bash
.ONESHELL:
.PHONY: list new-episode new-prep print-last-rt-link print-next-episode-number print-mp3-tags upload-mp3 deploy

TMPFILE := $(shell mktemp)
makefile_dir := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
notif :="/usr/local/bin/terminal-notifier"
subl := "/Applications/Sublime Text.app/Contents/SharedSupport/bin/subl"

list:
	@docker-compose run --rm publisher --list

new-episode:
	@docker-compose run --rm publisher new-episode | tee ${TMPFILE};
	@${subl} ${makefile_dir}/../`tail -n 1 ${TMPFILE} | tr -d '\r'`;

new-prep:
	@docker-compose run --rm publisher new-prep | tee ${TMPFILE};
	@${subl} ${makefile_dir}/../`tail -n 1 ${TMPFILE} | tr -d '\r'`;

print-last-rt-link:
	@docker-compose run --rm publisher print-last-rt-link

print-next-episode-number:
	@docker-compose run --rm publisher print-next-episode-number

print-mp3-tags:
	@docker-compose run --rm publisher print-mp3-tags $$FILE

upload-mp3:
	@docker-compose run --rm publisher upload-mp3 $$FILE | \
		while read line; do \
			if [[ $$line =~ ^!notif.* ]]; then \
				${notif} -title PodPrc -message "`echo $$line | sed 's/!notif: //'`"; \
			fi \
		done

deploy:
	@docker-compose run --rm publisher deploy