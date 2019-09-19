build:
	# compile build scss
	${GOPATH}/bin/packr2
	go build

dependencies:
	go get -u github.com/gobuffalo/packr/v2/packr2
	curl -fLo ${GOPATH}/bin/air https://raw.githubusercontent.com/cosmtrek/air/master/bin/linux/air && chmod +x ${GOPATH}/bin/air
	glide install
	# npm install

run:
	${GOPATH}/bin/air

release:
	#

clean:
	rm -rf ./tmp ./live-hls-transcode
	${GOPATH}/bin/packr2 clean
