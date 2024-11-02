# filter-rspamd-class  makefile

filter = filter-rspamd-class

fmt:
	fix go fmt . ./...

build: fmt
	fix go build

install: build
	doas install -m 0555 $(filter) /usr/local/libexec/smtpd/$(filter) && doas rcctl restart smtpd

test:
	fix -- go test -v . ./...
