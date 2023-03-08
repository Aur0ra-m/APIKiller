package hook

import "net/http"

type RequestHook interface {
	HookBefore(*http.Request) // hook before initiating http request
	HookAfter(*http.Request)  // hook after finishing http request
}

// hook sample
/*

type RequestHook interface {
	HookBefore(*http.Request) // hook before initiating http request
	HookAfter(*http.Request)  // hook after finishing http request
}

type AddHeaderHook struct {
}

func (a AddHeaderHook) HookBefore(request *http.Request) {

}

func (a AddHeaderHook) HookAfter(request *http.Request) {

}

// Hook this is exported, and this name must be set Hook
var Hook AddHeaderHook

*/
