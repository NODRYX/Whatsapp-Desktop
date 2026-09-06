package edge

type _ICoreWebView2NavigationStartingEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2NavigationStartingEventHandler struct {
	vtbl *_ICoreWebView2NavigationStartingEventHandlerVtbl
	impl _ICoreWebView2NavigationStartingEventHandlerImpl
}

func _ICoreWebView2NavigationStartingEventHandler_QueryInterface(this *ICoreWebView2NavigationStartingEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2NavigationStartingEventHandler_AddRef(this *ICoreWebView2NavigationStartingEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2NavigationStartingEventHandler_Release(this *ICoreWebView2NavigationStartingEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2NavigationStartingEventHandler_Invoke(this *ICoreWebView2NavigationStartingEventHandler, sender *ICoreWebView2, args *ICoreWebView2NavigationStartingEventArgs) uintptr {
	return this.impl.NavigationStarting(sender, args)
}

type _ICoreWebView2NavigationStartingEventHandlerImpl interface {
	_IUnknownImpl
	NavigationStarting(sender *ICoreWebView2, args *ICoreWebView2NavigationStartingEventArgs) uintptr
}

var _ICoreWebView2NavigationStartingEventHandlerFn = _ICoreWebView2NavigationStartingEventHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2NavigationStartingEventHandler_QueryInterface),
		NewComProc(_ICoreWebView2NavigationStartingEventHandler_AddRef),
		NewComProc(_ICoreWebView2NavigationStartingEventHandler_Release),
	},
	NewComProc(_ICoreWebView2NavigationStartingEventHandler_Invoke),
}

func newICoreWebView2NavigationStartingEventHandler(impl _ICoreWebView2NavigationStartingEventHandlerImpl) *ICoreWebView2NavigationStartingEventHandler {
	return &ICoreWebView2NavigationStartingEventHandler{
		vtbl: &_ICoreWebView2NavigationStartingEventHandlerFn,
		impl: impl,
	}
}
