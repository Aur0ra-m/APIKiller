package hook

import "net/http"

type RequestHook interface {
	HookBefore(*http.Request) // hook before initiating http request
	HookAfter(*http.Request)  // hook after finishing http request
}
