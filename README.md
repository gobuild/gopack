# gopack
[![Build Status](https://travis-ci.org/gobuild/gopack.svg)](https://travis-ci.org/gobuild/gopack)
[![gorelease](https://dn-gorelease.qbox.me/gorelease-download-blue.svg)](https://gobuild.io/gobuild/gopack)

Tool for [gobuild](https://gobuild.io)

## Features
1. Create `.gopack.yml` config file
2. Build and package build into zip
3. Download and install binary from <https://gobuild.io>

## Install
	go get -v github.com/gobuild/gopack

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

	$ gopack all
	Building linux amd64 -> output/gopack-linux-amd64.zip ...
	Building linux 386 -> output/gopack-linux-386.zip ...
	Building linux arm -> output/gopack-linux-arm.zip ...
	Building darwin amd64 -> output/gopack-darwin-amd64.zip ...
	Building windows amd64 -> output/gopack-windows-amd64.zip ...
	Building windows 386 -> output/gopack-windows-386.zip ...

	$ gopack install gocode
	==> Repository gobuild-official/gocode
	==> Downloading http://dn-gobuild5.qbox.me/gorelease/gobuild-official/gocode/master/darwin-amd64/gocode.zip
	 2.97 MB / 3.04 MB [===============================================>-]  97.75% 0Archive:  /Users/skyblue/.gopack/src/gocode.zip
	 3.04 MB / 3.04 MB [=================================================] 100.00% 0
	  inflating: /Users/skyblue/.gopack/opt/gobuild-official/gocode/README.md  
	  inflating: /Users/skyblue/.gopack/opt/gobuild-official/gocode/LICENSE  
	==> Symlink /Users/skyblue/Documents/godir/bin/gocode
	==> Program [gobuild-official/gocode] installed

see more flags in `gopack -h`

## LICENSE
Under [MIT](LICENSE)
