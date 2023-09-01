module github.com/takahawk/shadownet

go 1.21

replace github.com/takahawk/shadownet/gateway => ./gateway

replace github.com/takahawk/shadownet/downloaders => ./downloaders

replace github.com/takahawk/shadownet/transformers => ./transformers

replace github.com/takahawk/shadownet/resolvers => ./resolvers

replace github.com/takahawk/shadownet/common => ./common

require (
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/pborman/getopt v1.1.0 // indirect
)
