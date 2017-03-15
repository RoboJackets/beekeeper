
all:
# hack to get go to cooperate
	@export GOPATH=`pwd` && go install smoker/smoker

pretty:
	export GOPATH=`pwd` && go fmt smoker/...

fmt: pretty

clean:
	rm -r bin/

run: all
	@./bin/smoker

r: run

win:
	export GOPATH=`pwd` && GOOS=windows GOARCH=amd64 go build -o bin/smoker.exe smoker/smoker

linux:
	export GOPATH=`pwd` && GOOS=linux GOARCH=amd64 go build -o bin/smoker smoker/smoker

mac:
	export GOPATH=`pwd` && GOOS=darwin GOARCH=amd64 go build -o bin/smoker_mac smoker/smoker

multirelease: linux win mac
