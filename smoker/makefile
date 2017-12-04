
# hack to get go to cooperate
all:
	@export GOPATH=`pwd` && go install smoker/smoker

pretty:
	export GOPATH=`pwd` && go fmt smoker/...
	export GOPATH=`pwd` && go fmt varroa/...

# Run to generate enums (you need go stringer for this to be installed?!?!?!)
stringer:
	export GOPATH=`pwd` && go generate smoker/backends

fmt: pretty

clean:
	rm -r bin/

run: all
	@./bin/smoker

run-rlwrap: all
	@rlwrap ./bin/smoker

r: run
rr: run-rlwrap

varroa:
	@export GOPATH=`pwd` && go install varroa/varroa
varroa-run: varroa
	@./bin/varroa
vr: varroa-run
v: varroa


win:
	export GOPATH=`pwd` && GOOS=windows GOARCH=amd64 go build -o bin/smoker.exe smoker/smoker

linux:
	export GOPATH=`pwd` && GOOS=linux GOARCH=amd64 go build -o bin/smoker smoker/smoker

mac:
	export GOPATH=`pwd` && GOOS=darwin GOARCH=amd64 go build -o bin/smoker_mac smoker/smoker

netbsd:
	export GOPATH=`pwd` && GOOS=netbsd GOARCH=amd64 go build -o bin/smoker_netbsd smoker/smoker

openbsd:
	export GOPATH=`pwd` && GOOS=openbsd GOARCH=amd64 go build -o bin/smoker_openbsd smoker/smoker

freebsd:
	export GOPATH=`pwd` && GOOS=freebsd GOARCH=amd64 go build -o bin/smoker_freebsd smoker/smoker

varroa_dist:
	export GOPATH=`pwd` && GOOS=linux GOARCH=amd64 go build -o bin/varroa_linux varroa/varroa

# TODO support plan9 (bug golang color library to support it)
# plan9:
# 	export GOPATH=`pwd` && GOOS=plan9 GOARCH=amd64 go build -o bin/smoker_plan9 smoker/smoker

multirelease: linux freebsd netbsd openbsd win mac

dist: multirelease varroa_dist
