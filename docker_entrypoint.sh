set -e go get -d -v
set -e go build -buildmode=c-shared -o ./bin/pam_ohmage.so
echo "Build complete"