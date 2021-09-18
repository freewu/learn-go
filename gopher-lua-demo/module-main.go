package main

import (
	"github.com/yuin/gopher-lua"
	"gopher-lua-demo/module"
)

func main() {
	L := lua.NewState()
	defer L.Close()
	// Creating a module by Go
	L.PreloadModule("module", module.Loader)
	if err := L.DoFile("module-main.lua"); err != nil {
		panic(err)
	}
}