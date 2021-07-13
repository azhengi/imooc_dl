package execEnv

import (
	"fmt"
	"io/ioutil"

	"github.com/dop251/goja"
)

type JsRuntime struct {
	*goja.Runtime
}

var JsRt = &JsRuntime{}

func NewJsRuntime(codepath string) *JsRuntime {

	jsBytes, err := ioutil.ReadFile(codepath)

	vm := goja.New()

	_, err = vm.RunString(string(jsBytes))
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return nil
	}
	JsRt = &JsRuntime{}
	JsRt.Runtime = vm

	return JsRt
}

func (rt *JsRuntime) DecryptInfo(info, e string) []byte {

	handleDecrypt, ok := goja.AssertFunction(rt.Get("handleDecrypt"))
	if !ok {
		panic("handleDecrypt not a function")
	}
	res, err := handleDecrypt(goja.Undefined(), rt.ToValue(info), rt.ToValue(e))
	if err != nil {
		panic("handleDecrypt function exec fail")
	}

	return []byte(res.String())
}

func (rt *JsRuntime) EncryptPassword(pw string) []byte {
	handleCryptPassword, ok := goja.AssertFunction(rt.Get("handleCryptPassword"))
	if !ok {
		panic("handleCryptPassword not a function")
	}
	res, err := handleCryptPassword(goja.Undefined(), rt.ToValue(pw))
	if err != nil {
		panic(err)
	}

	return []byte(res.String())
}
