build:
	# compile build scss
	go build

dependencies:
	glide install
	# npm install

run:
	./vendor/github.com/cosmtrek/air/bin/linux/air

release:
	#
