package main

import (
	"fmt"
	"github.com/robertkrimen/otto"
)

func main() {
	vm := otto.New()
	// Run something in the VM
	_,_ = vm.Run(`
		var abc = 2 + 2;
		console.log("The value of abc is " + abc); // 4
	`) // The value of abc is 4
	// Get a value out of the VM
	if value, err := vm.Get("abc"); err == nil {
	   	if valueInt, err := value.ToInteger(); err == nil {
			fmt.Printf("%d,%v\n", valueInt, err) // 4,nil
	   	}
	}
	//Set a number
	_ = vm.Set("def", 11)
	_,_ = vm.Run(`
	   console.log("The value of def is " + def);
	`) // The value of def is 11

	//Set a string
	_ = vm.Set("xyzzy", "Nothing happens.")
	_,_ = vm.Run(`
		console.log(xyzzy.length);
	`) // 16

	// Get the value of an expression
	value, _ := vm.Run("xyzzy.length")
	v, err := value.ToInteger() // int64,error
	fmt.Printf("%d,%v\n", v, err) // 16,nil

	// An error happens
	value, err = vm.Run("abcdefghijlmnopqrstuvwxyz.length")
	if err != nil {
		fmt.Printf("%v\n",err) // ReferenceError: 'abcdefghijlmnopqrstuvwxyz' is not defined
	}

	// Set a Go function
	_ = vm.Set("sayHello", func(call otto.FunctionCall) otto.Value {
		fmt.Printf("Hello, %s.\n", call.Argument(0).String()) // get the first parameter as string
		return otto.Value{}
	})
	//Set a Go function that returns something useful
	_ = vm.Set("twoPlus", func(call otto.FunctionCall) otto.Value {
		right, _ := call.Argument(0).ToInteger()
		result, _ := vm.ToValue(2 + right)
		return result
	})
	// Use the functions in JavaScript
	result, err := vm.Run(`
		sayHello("Xyzzy");      // Hello, Xyzzy.
		sayHello();             // Hello, undefined
		result = twoPlus(2.0); // 4
	`)
	fmt.Printf("%v,%v\n",result,err) // 4,<nil>
}