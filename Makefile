# go makefile

program != basename $$(pwd)
latest_release != gh release list --json tagName --jq '.[0].tagName' | tr -d v
version != cat VERSION
install_dir = /usr/local/libexec/smtpd
postinstall = && doas rcctl restart smtpd

build: fmt
	fix go build

fmt: go.sum
	fix go fmt . ./...

go.mod:
	go mod init

go.sum: go.mod
	go mod tidy

install: build
	doas install -m 0755 $(program) $(install_dir)/$(program) $(postinstall)

test:
	fix -- go test -v . ./...

release:
	@gitclean -v -d "git status is dirty"
	echo latest_release=$(latest_release)
	[ "$(latest_release)" != $(version) ] 
	echo gh release create v$(shell cat VERSION) --notes "v$(shell cat VERSION)"

clean:
	rm -f $(program)
	go clean

sterile: clean
	go clean -r
	go clean -cache
	go clean -modcache
	rm -f go.mod go.sum
