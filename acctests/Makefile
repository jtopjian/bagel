PWD:=$(shell pwd)

build:
	go build -o bagel $(GOPATH)/src/github.com/jtopjian/bagel

apt_key: build
	docker run -v $(PWD):$(PWD) -w $(PWD) -it --rm ubuntu:18.04 ./bagel run apt_key.lua --debug

apt_package: build
	docker run -v $(PWD):$(PWD) -w $(PWD) -it --rm ubuntu:18.04 ./bagel run apt_package.lua --debug

apt_ppa: build
	docker run -v $(PWD):$(PWD) -w $(PWD) -it --rm ubuntu:18.04 ./bagel run apt_ppa.lua --debug

apt_source: build
	docker run -v $(PWD):$(PWD) -w $(PWD) -it --rm ubuntu:18.04 ./bagel run apt_source.lua --debug

cron_entry: build
	docker run -v $(PWD):$(PWD) -w $(PWD) -it --rm ubuntu:18.04 ./bagel run cron_entry.lua --debug
