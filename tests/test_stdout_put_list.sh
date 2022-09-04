go run pstore.go -endpoint=http://127.0.0.1:5000 -del -name=testtest; \
go run pstore.go -endpoint=http://127.0.0.1:5000 -put -name=/test/123/456 -value=testtest; \
go run pstore.go -endpoint=http://127.0.0.1:5000
