package edge

import (
	"unsafe"

	"github.com/jchv/go-webview2/internal/w32"
	"golang.org/x/sys/windows"
)

type _ICoreWebView2NewWindowRequestedEventArgsVtbl struct {
	_IUnknownVtbl
	GetUri             ComProc
	GetNewWindow       ComProc
	PutNewWindow       ComProc
	GetDeferral        ComProc
	GetWindowFeatures  ComProc
	GetIsUserInitiated ComProc
	GetHandled         ComProc
	PutHandled         ComProc
}

type ICoreWebView2NewWindowRequestedEventArgs struct {
	vtbl *_ICoreWebView2NewWindowRequestedEventArgsVtbl
}

func (a *ICoreWebView2NewWindowRequestedEventArgs) GetUri() (string, error) {
	var pwstr *uint16
	_, _, err := a.vtbl.GetUri.Call(uintptr(unsafe.Pointer(a)), uintptr(unsafe.Pointer(&pwstr)))
	if err != windows.ERROR_SUCCESS {
		return "", err
	}
	defer windows.CoTaskMemFree(unsafe.Pointer(pwstr))
	return w32.Utf16PtrToString(pwstr), nil
}

func (a *ICoreWebView2NewWindowRequestedEventArgs) PutHandled(handled bool) error {
	v := 0
	if handled {
		v = 1
	}
	_, _, err := a.vtbl.PutHandled.Call(uintptr(unsafe.Pointer(a)), uintptr(v))
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (a *ICoreWebView2NewWindowRequestedEventArgs) GetIsUserInitiated() (bool, error) {
	var b int32
	_, _, err := a.vtbl.GetIsUserInitiated.Call(uintptr(unsafe.Pointer(a)), uintptr(unsafe.Pointer(&b)))
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return b != 0, nil
}
