all: hello dump_log

hello:
	go build -o ./bin/hello ./tests/hello

dump_log:
	go build -o ./bin/dump_log ./tests/dump_log


# build:
# 	go build -o ./bin/hello ./tests/hello

clean:
	rm ./bin/hello
	rm ./bin/dump_log

