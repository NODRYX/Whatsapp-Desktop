package edge

type _ICoreWebView2NewWindowRequestedEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2NewWindowRequestedEventHandler struct {
	vtbl *_ICoreWebView2NewWindowRequestedEventHandlerVtbl
	impl _ICoreWebView2NewWindowRequestedEventHandlerImpl
}

func _ICoreWebView2NewWindowRequestedEventHandler_QueryInterface(this *ICoreWebView2NewWindowRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2NewWindowRequestedEventHandler_AddRef(this *ICoreWebView2NewWindowRequestedEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2NewWindowRequestedEventHandler_Release(this *ICoreWebView2NewWindowRequestedEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2NewWindowRequestedEventHandler_Invoke(this *ICoreWebView2NewWindowRequestedEventHandler, sender *ICoreWebView2, args *ICoreWebView2NewWindowRequestedEventArgs) uintptr {
	return this.impl.NewWindowRequested(sender, args)
}

type _ICoreWebView2NewWindowRequestedEventHandlerImpl interface {
	_IUnknownImpl
	NewWindowRequested(sender *ICoreWebView2, args *ICoreWebView2NewWindowRequestedEventArgs) uintptr
}

var _ICoreWebView2NewWindowRequestedEventHandlerFn = _ICoreWebView2NewWindowRequestedEventHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2NewWindowRequestedEventHandler_QueryInterface),
		NewComProc(_ICoreWebView2NewWindowRequestedEventHandler_AddRef),
		NewComProc(_ICoreWebView2NewWindowRequestedEventHandler_Release),
	},
	NewComProc(_ICoreWebView2NewWindowRequestedEventHandler_Invoke),
}

func newICoreWebView2NewWindowRequestedEventHandler(impl _ICoreWebView2NewWindowRequestedEventHandlerImpl) *ICoreWebView2NewWindowRequestedEventHandler {
	return &ICoreWebView2NewWindowRequestedEventHandler{
		vtbl: &_ICoreWebView2NewWindowRequestedEventHandlerFn,
		impl: impl,
	}
}
