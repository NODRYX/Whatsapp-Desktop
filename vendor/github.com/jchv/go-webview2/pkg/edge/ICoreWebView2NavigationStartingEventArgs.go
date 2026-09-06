package edge

import (
	"unsafe"

	"github.com/jchv/go-webview2/internal/w32"
	"golang.org/x/sys/windows"
)

type _ICoreWebView2NavigationStartingEventArgsVtbl struct {
	_IUnknownVtbl
	GetUri             ComProc
	GetIsUserInitiated ComProc
	GetIsRedirected    ComProc
	GetRequestHeaders  ComProc
	PutCancel          ComProc
	GetCancel          ComProc
	GetNavigationId    ComProc
}

type ICoreWebView2NavigationStartingEventArgs struct {
	vtbl *_ICoreWebView2NavigationStartingEventArgsVtbl
}

func (a *ICoreWebView2NavigationStartingEventArgs) GetUri() (string, error) {
	var pwstr *uint16
	_, _, err := a.vtbl.GetUri.Call(uintptr(unsafe.Pointer(a)), uintptr(unsafe.Pointer(&pwstr)))
	if err != windows.ERROR_SUCCESS {
		return "", err
	}
	defer windows.CoTaskMemFree(unsafe.Pointer(pwstr))
	return w32.Utf16PtrToString(pwstr), nil
}

func (a *ICoreWebView2NavigationStartingEventArgs) PutCancel(cancel bool) error {
	v := 0
	if cancel {
		v = 1
	}
	_, _, err := a.vtbl.PutCancel.Call(uintptr(unsafe.Pointer(a)), uintptr(v))
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (a *ICoreWebView2NavigationStartingEventArgs) GetIsUserInitiated() (bool, error) {
	var b int32
	_, _, err := a.vtbl.GetIsUserInitiated.Call(uintptr(unsafe.Pointer(a)), uintptr(unsafe.Pointer(&b)))
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return b != 0, nil
}
