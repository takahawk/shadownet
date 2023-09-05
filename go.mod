module github.com/takahawk/shadownet

go 1.21

replace github.com/takahawk/shadownet/gateway => ./gateway

replace github.com/takahawk/shadownet/downloaders => ./downloaders

replace github.com/takahawk/shadownet/transformers => ./transformers

replace github.com/takahawk/shadownet/resolvers => ./resolvers

replace github.com/takahawk/shadownet/common => ./common

require (
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pborman/getopt v1.1.0 // indirect
	github.com/rs/zerolog v1.30.0 // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)
