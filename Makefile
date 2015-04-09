MAIN=flood.go
TARGET=vde.test

all:
	${GO} build github.com/kurojishi/vde-testing

install:
	${GO} install github.com/kurojishi/vde-testing

#test:
	#${GO} test -compiler ${GC} github.com/kurojishi/vde-testing

clean:
	rm -f ${TARGET}

