module github.com/takahawk/shadownet

go 1.21.0

replace github.com/takahawk/shadownet/downloaders => ./downloaders

replace github.com/takahawk/shadownet/encryptors => ./encryptors

replace github.com/takahawk/shadownet/transformers => ./transformers

replace github.com/takahawk/shadownet/resolvers => ./resolvers

replace github.com/takahawk/shadownet/common => ./common

require github.com/pborman/getopt v1.1.0 // indirect
