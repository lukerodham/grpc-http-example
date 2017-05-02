# /bin/bash

PROTOPATH=$GOPATH/src/github.com/lukerodham/grpc-http-example/proto
FILES=$PROTOPATH/*.proto

for f in $FILES
do
  file=$(basename "$f")
  package="${file%.*}"

  echo "--------------"
  echo "Processing $package $f \n"

  mkdir -p $package

  protoc -I $PROTOPATH --go_out=plugins=micro:./../../../../ $f

  protoc-go-inject-tag -input=`pwd`/$package/$package.pb.go
done
