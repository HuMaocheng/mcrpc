package main

import (
	"fmt"
	"reflect"
)

func main(){
	testServer := new(Server)
	testServer.funcs = make(map[string]reflect.Value)
	testServer.Register("getMax", getMax)
	v := testServer.funcs["getMax"]

	arg := make([]reflect.Value, 2)
	arg[0] = reflect.ValueOf(1)
	arg[1] = reflect.ValueOf(2)
	out := v.Call(arg)
	fmt.Println(out[0])

}


func getMax(a, b int) (int){
	if a > b {
		return a
	} else {
		return b
	}
}
