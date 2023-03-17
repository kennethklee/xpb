package main

import (
	"github.com/kennethklee/xpb/xpocketbase/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		panic(err)
	}
}
