module github.com/linkbox-group/linkbox-server/user

go 1.23.6

replace github.com/linkbox-group/linkbox-server/rpc-gen => ../../rpc-gen

require (
	github.com/google/wire v0.6.0
	github.com/linkbox-group/linkbox-server/rpc-gen v0.0.0-00010101000000-000000000000
)

require (
	github.com/cloudwego/prutal v0.1.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
