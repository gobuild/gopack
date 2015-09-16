# gopack
[![Build Status](https://travis-ci.org/gorelease/gopack.svg)](https://travis-ci.org/gorelease/gopack)
[![gorelease](https://dn-gorelease.qbox.me/gorelease-download-blue.svg)](http://gorelease.herokuapp.com/gorelease/gopack)

Tool for [gorelease](https://github.com/gorelease/gorelease)

## Install
	go get -v github.com/gorelease/gopack

## Usage
	$ gopack init
	# create .gopack.yml config file

	$ gopack pack
	# build go code and package README.md LICENSE ... to a zip file
	[golang-sh]$ bash -c go get -v
	[golang-sh]$ bash -c go install
	2015/09/16 23:30:35 [Info] pack.go:183 zip add file: gopack
	2015/09/16 23:30:35 [Info] pack.go:183 zip add file: README.md
	2015/09/16 23:30:35 [Info] pack.go:183 zip add file: LICENSE
	2015/09/16 23:30:35 [Info] pack.go:188 finish archive file

	$ unzip -t gopack.zip
	Archive:  gopack.zip
		testing: gopack                   OK
		testing: README.md                OK
		testing: LICENSE                  OK
	No errors detected in compressed data of gopack.zip.

see more flags in `gopack -h`

## LICENSE
Under [MIT](LICENSE)
