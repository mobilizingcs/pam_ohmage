set -e
go get -d -v
go build -buildmode=c-shared -o ./bin/pam_ohmage.so
echo "Build complete"