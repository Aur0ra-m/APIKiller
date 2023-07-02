package hook

var Hooks []RequestHook

//
// RegisterHooks
//  @Description: append http request hook to modify request data
//  @param requestHook
//
func RegisterHooks(requestHook RequestHook) {
	Hooks = append(Hooks, requestHook)
}
