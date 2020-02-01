.DEFAULT_GOAL := list
SHELL:=/bin/bash
.ONESHELL:
.PHONY: list new-episode new-prep print-last-rt-link print-next-episode-number print-mp3-tags upload-mp3 deploy

TMPFILE := $(shell mktemp)
makefile_dir := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
notif := "/usr/local/bin/terminal-notifier"
subl := "/Applications/Sublime Text.app/Contents/SharedSupport/bin/subl"

list:
	@docker-compose run --rm publisher --list

# generate new episode post markdown file and open it using SublimeText 
new-episode:
	@docker-compose run --rm --entrypoint /bin/sh publisher scripts/make_new_episode.sh | tee ${TMPFILE};
	@${subl} ${makefile_dir}/../hugo/`tail -n 1 ${TMPFILE} | tr -d '\r'`;

# generate new prep post markdown file and open it using SublimeText 
new-prep:
	@docker-compose run --rm --entrypoint /bin/sh publisher scripts/make_new_prep.sh | tee ${TMPFILE};
	@${subl} ${makefile_dir}/../hugo/`tail -n 1 ${TMPFILE} | tr -d '\r'`;

# print url to last podcast episode on https://radio-t.com
print-last-rt-link:
	@docker-compose run --rm publisher print-last-rt-link

# print next episode number using info from https://radio-t.com
print-next-episode-number:
	@docker-compose run --rm publisher print-next-episode-number

# print mp3 tags from file
print-mp3-tags:
	@FDIR=`dirname "$$FILE" | xargs -0 basename` && FNAME=`basename "$$FILE"` && \
		docker-compose run --rm publisher print-mp3-tags "$$FDIR/$$FNAME";

# set necessary tags to episode mp3 file and upload it, notifying via MacOS local notifications on every step
upload-mp3:
	@FDIR=`dirname "$$FILE" | xargs -0 basename` && FNAME=`basename "$$FILE"` && \
	docker-compose run --rm --entrypoint /bin/sh publisher scripts/upload_mp3.sh "$$FDIR/$$FNAME" | \
		while read line; do \
			if [[ $$line =~ ^!notif.* ]]; then \
				${notif} -title PodPrc -message "`echo $$line | sed 's/!notif: //'`"; \
			else \
				echo $$line; \
			fi \
		done

# deploy new podcast episode page to https://radio-t.com
deploy:
	@docker-compose run --rm --entrypoint /bin/sh publisher scripts/deploy.sh