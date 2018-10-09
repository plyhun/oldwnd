package resources

import (
	_ "github.com/jteeuwen/go-bindata"
)

//go:generate $GOPATH/bin/go-bindata -nomemcopy=true -pkg=resources shaders
