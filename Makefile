build:
	go build pick-aws-profile.go

install: build
	ln -sf $$PWD/pick-aws-profile /usr/local/bin
