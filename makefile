all: hello dump_log system_monitor

# 定义hello的目标文件
hello: ./bin/hello
# 目标文件生成规则
./bin/hello:
	go build -o ./bin/hello ./tests/hello
# 给出清理规则
clean_hello:
	rm -f ./bin/hello

dump_log: ./bin/dump_log
./bin/dump_log:
	go build -o ./bin/dump_log ./tests/dump_log
clean_dump_log:
	rm -f ./bin/dump_log

system_monitor: ./bin/system_monitor
./bin/system_monitor:
	go build -o ./bin/system_monitor ./tests/system_monitor
clean_system_monitor:
	rm -f ./bin/system_monitor

# build:
# 	go build -o ./bin/hello ./tests/hello

# .PHONY 定义虚拟目标？
.PHONY: clean clean_hello clean_dump_log clean_system_monitor
clean: clean_hello clean_dump_log clean_system_monitor




