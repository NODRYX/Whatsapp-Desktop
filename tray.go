//go:build windows

package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	shell32          = windows.NewLazySystemDLL("shell32.dll")
	procShellNotify  = shell32.NewProc("Shell_NotifyIconW")
	procLoadIcon     = user32.NewProc("LoadIconW")
)

const (
	nimAdd       = 0x00000000
	nimModify    = 0x00000001
	nimDelete    = 0x00000002
	nifMessage   = 0x00000001
	nifIcon      = 0x00000002
	nifTip       = 0x00000004
	wmTrayIcon   = 0x8001 // WM_APP+1
	wmLButtonUp  = 0x0202
	wmLButtonDblClk = 0x0203
	wmRButtonUp  = 0x0205
)

type notifyIconData struct {
	CbSize           uint32
	HWnd             windows.Handle
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	HIcon            windows.Handle
	SzTip            [128]uint16
	DwState          uint32
	DwStateMask      uint32
	SzInfo           [256]uint16
	UVersion         uint32
	SzInfoTitle      [64]uint16
	DwInfoFlags      uint32
	GuidItem         windows.GUID
	HBalloonIcon     windows.Handle
}

func addTrayIcon(hwnd uintptr) bool {
	var hInstance windows.Handle
	_ = windows.GetModuleHandleEx(0, nil, &hInstance)
	// Load icon ID 2 from resources (same as window icon)
	icon, _, _ := procLoadIcon.Call(uintptr(hInstance), uintptr(2))
	if icon == 0 {
		icon, _, _ = user32.NewProc("LoadIconW").Call(0, uintptr(32512)) // IDI_APPLICATION
	}
	var tip [128]uint16
	copy(tip[:], windows.StringToUTF16("WhatsApp Desktop"))
	var nid notifyIconData
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.HWnd = windows.Handle(hwnd)
	nid.UID = 1
	nid.UFlags = nifMessage | nifIcon | nifTip
	nid.UCallbackMessage = wmTrayIcon
	nid.HIcon = windows.Handle(icon)
	nid.SzTip = tip
	ret, _, _ := procShellNotify.Call(uintptr(nimAdd), uintptr(unsafe.Pointer(&nid)))
	return ret != 0
}

func removeTrayIcon(hwnd uintptr) {
	var nid notifyIconData
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.HWnd = windows.Handle(hwnd)
	nid.UID = 1
	_, _, _ = procShellNotify.Call(uintptr(nimDelete), uintptr(unsafe.Pointer(&nid)))
}
