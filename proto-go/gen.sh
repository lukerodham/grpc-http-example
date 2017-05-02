# /bin/bash

PROTOPATH=$GOPATH/src/github.com/sipsynergy/accounts-proto
FILES=$PROTOPATH/*.proto

for f in $FILES
do
  file=$(basename "$f")
  package="${file%.*}"

  echo "Processing $package $f \n"

  mkdir -p $package

  protoc -I $PROTOPATH --go_out=plugins=micro:../../../ $f

  protoc-go-inject-tag -input=`pwd`/$package/$package.pb.go
done

#protoc --go_out=plugins=grpc:./ ./proto/user.proto
#protoc-go-inject-tag -input=./proto/user.pb.go

#protoc --go_out=plugins=grpc:./ ./proto/organisation.proto
#protoc-go-inject-tag -input=./proto/organisation.pb.go
