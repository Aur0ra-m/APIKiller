package hook

import "net/http"

type RequestHook interface {
	HookBefore(*http.Request) // hooks before initiating http request
	HookAfter(*http.Request)  // hooks after finishing http request
}
