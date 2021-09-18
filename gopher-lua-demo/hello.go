package main

import (
	"fmt"
	"github.com/yuin/gopher-lua"
)

func main() {
	// Run scripts content in the VM.
	//L := lua.NewState()
	L := lua.NewState(lua.Options{
		// Registry
		RegistrySize: 1024 * 20, // this is the initial size of the registry
		RegistryMaxSize: 1024 * 80, // this is the maximum size that the registry can grow to. If set to `0` (the default) then the registry will not auto grow
		RegistryGrowStep: 32, // this is how much to step up the registry by each time it runs out of space. The default is `32`.
		// Callstack
		CallStackSize: 120, // this is the maximum callstack size of this LState
		MinimizeStackMemory: true, // Defaults to `false` if not specified. If set, the callstack will auto grow and shrink as needed up to a max of `CallStackSize`. If not set, the callstack will be fixed at `CallStackSize`.
	})

	defer L.Close()
	if err := L.DoString(`print("hello")`); err != nil {
		fmt.Printf("%v\n",err)
	}
	// Run scripts file in the VM.
	if err := L.DoFile("hello.lua"); err != nil {
		fmt.Printf("%v\n",err)
	}
	// Run scripts file in the VM.
	if err := L.DoFile("hello1.lua"); err != nil {
		fmt.Printf("%v\n",err) // open hello1.lua: The system cannot find the file specified.
	}

	// Register a lua function
	L.SetGlobal("double", L.NewFunction(func (L *lua.LState) int {
		lv := L.ToInt(1)  /* get first argument */
		L.Push(lua.LNumber(lv * 2)) /* push result */
		return 1 /* number of results */
	}))
	// Call the custom lua function
	if err := L.DoString(`print(double(3))`); err != nil {
		fmt.Printf("Call the custom lua function: %v\n",err)
	}

	// Calling Lua from Go
	if err := L.CallByParam(lua.P{ Fn: L.GetGlobal("double"), NRet: 1, Protect: true, }, lua.LNumber(10)); err != nil {
		fmt.Printf("Calling Lua from Go: %v\n",err)
	}
	ret := L.Get(-1) // returned value
	fmt.Printf("type: %v, value: %v",ret.Type(), ret.String()) // number,20
	L.Pop(1)  // remove received value
}

/*
Type name	Go type			Type() value	Constants
LNilType	(constants)		LTNil			LNil
LBool		(constants)		LTBool			LTrue, LFalse
LNumber		float64			LTNumber		-
LString		string			LTString		-
LFunction	struct pointer	LTFunction		-
LUserData	struct pointer	LTUserData		-
LState		struct pointer	LTThread		-
LTable		struct pointer	LTTable			-
LChannel	chan 			LValue			LTChannel	-
 */
