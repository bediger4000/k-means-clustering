all: genrand genblob km1
	./do7
	./doblob

km1: km1.go
	go build km1.go

genrand: genrand.go
	go build genrand.go
genblob: genrand.go
	go build genblob.go

clean:
	go clean
	-rm -rf clust*
	-rm -rf blob cent randx out
