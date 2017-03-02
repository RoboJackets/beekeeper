
all:
# hack to get go to cooperate
	mkdir -p bin
	export GOPATH=`pwd` && go install smoker/smoker

pretty:
	export GOPATH=`pwd` && go fmt smoker/...

fmt: pretty

clean:
	rm -r bin/

run: all
	./bin/smoker

r: run
