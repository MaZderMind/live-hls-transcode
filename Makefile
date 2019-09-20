build: clean
	${GOPATH}/bin/packr2
	go build

dependencies:
	go get github.com/gobuffalo/packr/v2/packr2
	curl -fLo ${GOPATH}/bin/air https://raw.githubusercontent.com/cosmtrek/air/master/bin/linux/air && chmod +x ${GOPATH}/bin/air
	glide install
	cd frontend && npm install

run: clean
	${GOPATH}/bin/air

release: clean
	#

clean:
	rm -rf ./tmp ./live-hls-transcode
	${GOPATH}/bin/packr2 clean

.PHONY: build clean dependenciesrun release
