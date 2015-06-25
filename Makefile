MAIN=flood.go
TARGET=vde.test

all:
	go build github.com/kurojishi/vdetesting/util
	go build github.com/kurojishi/vdetesting
	go build github.com/kurojishi/vdetesting/example


install:
	go install github.com/kurojishi/vdetesting/utils
	go install github.com/kurojishi/vdetesting
	go install github.com/kurojishi/vdetesting/example

#test:
	#${GO} test -compiler ${GC} github.com/kurojishi/vde-go

clean:
	rm -f ${TARGET}

