package execCtx

import (
	"github.com/dop251/goja"
)

type JsRuntime struct {
	*goja.Runtime
}

func NewJsRuntime(jsCode string) *JsRuntime {
	vm := goja.New()

	_, err := vm.RunString(jsCode)
	if err != nil {
		return nil
	}
	rt := &JsRuntime{vm}

	return rt
}

func (rt *JsRuntime) DecodeM3u8(content string) []byte {
	decodeM3u8, _ := goja.AssertFunction(rt.Get("decodeM3u8"))
	res, _ := decodeM3u8(goja.Undefined(), rt.ToValue(content))
	return []byte(res.String())
}

func (rt *JsRuntime) DecodeKey(key string) []byte {
	decodeKey, _ := goja.AssertFunction(rt.Get("decodeKey"))
	res, _ := decodeKey(goja.Undefined(), rt.ToValue(key))
	return []byte(res.String())
}
