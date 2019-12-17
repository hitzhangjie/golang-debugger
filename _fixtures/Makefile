FLAGS=-gcflags="-N -l"
GOGC=go
CMD=build
BUILD=${GOGC} ${CMD} ${FLAGS}
all:	*.go
	${BUILD} continuetestprog.go
	${BUILD} integrationprog.go
	${BUILD} livetestprog.go
	${BUILD} testnextprog.go
	${BUILD} testprog.go
	${BUILD} testvariables.go

clean:
	rm -f continuetestprog
	rm -f integrationprog
	rm -f livetestprog
	rm -f testnextprog
	rm -f testprog
	rm -f testvariables
	
