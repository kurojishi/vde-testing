MAIN=flood.go
TARGET=vde.test

all:
	go build github.com/kurojishi/vde-testing

install:
	go install github.com/kurojishi/vde-testing

#test:
	#${GO} test -compiler ${GC} github.com/kurojishi/vde-go

clean:
	rm -f ${TARGET}

