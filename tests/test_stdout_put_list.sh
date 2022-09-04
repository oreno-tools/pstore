export AWS_ACCESS_KEY_ID=foo
export AWS_SECRET_ACCESS_KEY=bar
export AWS_DEFAULT_REGION=us-east-1

go run pstore.go -endpoint=http://127.0.0.1:5000 -del -name=testtest; \
go run pstore.go -endpoint=http://127.0.0.1:5000 -put -name=/test/123/456 -value=testtest; \
go run pstore.go -endpoint=http://127.0.0.1:5000
