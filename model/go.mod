module github.com/linkbox-group/linkbox-server/model

go 1.23.6

replace github.com/linkbox-group/linkbox-server/rpc-gen => ../rpc-gen

require (
	github.com/google/uuid v1.6.0
	github.com/linkbox-group/linkbox-server/rpc-gen v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.3
	gorm.io/gorm v1.25.12
)

require (
	github.com/cloudwego/prutal v0.1.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
