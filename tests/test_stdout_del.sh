export AWS_ACCESS_KEY_ID=foo
export AWS_SECRET_ACCESS_KEY=bar
export AWS_DEFAULT_REGION=us-east-1

echo y | go run pstore.go -endpoint=http://127.0.0.1:5000 -del -name=test123; go run pstore.go -endpoint=http://127.0.0.1:5000
