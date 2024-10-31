# filter-rspamd-class  makefile

filter = filter-rspamd-class

fmt:
	go fmt

build:
	go build

install: build
	doas install -m 0555 $(filter) /usr/local/libexec/smtpd/$(filter) && doas rcctl restart smtpd

test:
	go test -v
