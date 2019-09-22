build: clean
	${GOPATH}/bin/packr2
	go build

dependencies:
	go get github.com/gobuffalo/packr/v2/packr2
	cd frontend && npm install

run: clean
	go get github.com/codegangsta/gin
	ROOT_DIR=/video/ ${GOPATH}/bin/gin --all --port 8048

release: clean
	#

clean:
	rm -rf ./tmp ./live-hls-transcode
	${GOPATH}/bin/packr2 clean

.PHONY: build clean dependenciesrun release
