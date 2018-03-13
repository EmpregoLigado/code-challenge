//Package proto has the protocol buffer related files
//To generate the service go code you'll require
//protobuf - brew install protobuf
//protoc-gen-go - go get github.com/golang/protobuf/protoc-gen-go
//protoc-gen-twirp - go get github.com/twitchtv/twirp/protoc-gen-twirp
//
//protoc requires protoc-gen-go/protoc-gen-twirp to be found in your $PATH
//be sure to configure a complete path to the binary, like /home/johndoe/go/bin
//protoc will not expand shortcuts like ~/go/bin
//
//protoc --twirp_out=. *.proto
//protoc --go_out=plugins=grpc:. *.proto
//
//None of that is needed to run the application, only if the .proto files are changed.
package proto
