# filter-rspamd-class  makefile

filter = filter-rspamd-class

build: fmt
	fix go build

fmt: go.sum
	fix go fmt . ./...

go.mod:
	go mod init

go.sum: go.mod
	go mod tidy

install: build
	doas install -m 0555 $(filter) /usr/local/libexec/smtpd/$(filter) && doas rcctl restart smtpd

test:
	fix -- go test -v . ./...

release:
	bump
	gh release create v$(shell cat VERSION) --notes "v$(shell cat VERSION)"

clean:
	rm -f $(filter)
	go clean

sterile: clean
	rm -f go.mod go.sum
