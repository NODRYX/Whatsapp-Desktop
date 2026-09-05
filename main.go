package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/go-toast/toast"
	"github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"
)

var (
	kernel32        = windows.NewLazySystemDLL("kernel32.dll")
	user32          = windows.NewLazySystemDLL("user32.dll")
	dwmapi          = windows.NewLazySystemDLL("dwmapi.dll")
	procCreateMutex = kernel32.NewProc("CreateMutexW")
	procFindWindow  = user32.NewProc("FindWindowW")
	procSetFgWindow = user32.NewProc("SetForegroundWindow")
	procShowNormal  = user32.NewProc("ShowWindow")
	procDwmSetAttr  = dwmapi.NewProc("DwmSetWindowAttribute")
	// geometry helpers
	procGetWindowRect = user32.NewProc("GetWindowRect")
	procIsZoomed      = user32.NewProc("IsZoomed")
	procSetWindowPos  = user32.NewProc("SetWindowPos")
)

const (
	windowTitle = "WhatsApp Desktop"
	appURL      = "https://web.whatsapp.com"
	mutexName   = "WhatsAppDesktopSingleInstanceMutex"

	// DWM Window Attributes for Dark Theme
	DWMWA_USE_IMMERSIVE_DARK_MODE_BEFORE_20H1 = 19
	DWMWA_USE_IMMERSIVE_DARK_MODE             = 20
	DWMWA_CAPTION_COLOR                      = 35
	DWMWA_TEXT_COLOR                         = 36

	// Window state file version
	windowStateVersion = 1
)

type windowState struct {
	Version     int  `json:"version"`
	Width       int  `json:"width"`
	Height      int  `json:"height"`
	X           int  `json:"x"`
	Y           int  `json:"y"`
	Maximized   bool `json:"maximized"`
	CloseToTray bool `json:"closeToTray,omitempty"`
}

type w32Rect struct {
	Left, Top, Right, Bottom int32
}

func setDarkWindowFrame(hwnd uintptr, dark bool) {
	// immersive dark mode toggle (Win10 20H1+ & Win11)
	darkMode := int32(0)
	if dark {
		darkMode = 1
	}
	procDwmSetAttr.Call(
		hwnd,
		uintptr(DWMWA_USE_IMMERSIVE_DARK_MODE),
		uintptr(unsafe.Pointer(&darkMode)),
		unsafe.Sizeof(darkMode),
	)
	procDwmSetAttr.Call(
		hwnd,
		uintptr(DWMWA_USE_IMMERSIVE_DARK_MODE_BEFORE_20H1),
		uintptr(unsafe.Pointer(&darkMode)),
		unsafe.Sizeof(darkMode),
	)

	var captionColor uint32
	var textColor uint32
	if dark {
		captionColor = uint32(0x00211B11) // 0x00BBGGRR WhatsApp Dark Header 17,27,33
		textColor = uint32(0x00FFFFFF)
	} else {
		captionColor = uint32(0x00FFFFFF) // white caption for light mode
		textColor = uint32(0x00000000)
	}
	procDwmSetAttr.Call(
		hwnd,
		uintptr(DWMWA_CAPTION_COLOR),
		uintptr(unsafe.Pointer(&captionColor)),
		unsafe.Sizeof(captionColor),
	)
	procDwmSetAttr.Call(
		hwnd,
		uintptr(DWMWA_TEXT_COLOR),
		uintptr(unsafe.Pointer(&textColor)),
		unsafe.Sizeof(textColor),
	)
}

func isDarkModeEnabled() bool {
	// Registry: HKCU\Software\Microsoft\Windows\CurrentVersion\Themes\Personalize\AppsUseLightTheme (0=dark)
	var key windows.Handle
	subkey, _ := windows.UTF16PtrFromString(`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`)
	if err := windows.RegOpenKeyEx(windows.HKEY_CURRENT_USER, subkey, 0, windows.KEY_QUERY_VALUE, &key); err != nil {
		return true
	}
	defer windows.RegCloseKey(key)
	name, _ := windows.UTF16PtrFromString(`AppsUseLightTheme`)
	var typ uint32
	var data uint32
	var n uint32 = 4
	if err := windows.RegQueryValueEx(key, name, nil, &typ, (*byte)(unsafe.Pointer(&data)), &n); err != nil {
		return true
	}
	return data == 0
}

func checkSingleInstance() (uintptr, bool) {
	namePtr, _ := syscall.UTF16PtrFromString(mutexName)
	handle, _, err := procCreateMutex.Call(0, 1, uintptr(unsafe.Pointer(namePtr)))
	if err == windows.ERROR_ALREADY_EXISTS {
		titlePtr, _ := syscall.UTF16PtrFromString(windowTitle)
		hwnd, _, _ := procFindWindow.Call(0, uintptr(unsafe.Pointer(titlePtr)))
		if hwnd != 0 {
			procShowNormal.Call(hwnd, 9) // SW_RESTORE
			procSetFgWindow.Call(hwnd)
		}
		return handle, false
	}
	return handle, true
}

func getUserDataDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			configDir = "."
		}
	}
	dir := filepath.Join(configDir, "WhatsAppDesktopLight", "UserData")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

func getWindowStatePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			configDir = "."
		}
	}
	dir := filepath.Join(configDir, "WhatsAppDesktopLight")
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "window.json")
}

func loadWindowState() (*windowState, bool) {
	path := getWindowStatePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	var s windowState
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, false
	}
	if s.Version != windowStateVersion {
		return nil, false
	}
	if s.Width < 400 || s.Width > 5000 || s.Height < 300 || s.Height > 5000 {
		return nil, false
	}
	return &s, true
}

func saveWindowState(hwnd uintptr) {
	if hwnd == 0 {
		return
	}
	var rect w32Rect
	ret, _, _ := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&rect)))
	if ret == 0 {
		return
	}
	width := int(rect.Right - rect.Left)
	height := int(rect.Bottom - rect.Top)
	if width < 400 || height < 300 {
		return
	}
	maximized, _, _ := procIsZoomed.Call(hwnd)
	s := windowState{
		Version:   windowStateVersion,
		Width:     width,
		Height:    height,
		X:         int(rect.Left),
		Y:         int(rect.Top),
		Maximized: maximized != 0,
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}
	path := getWindowStatePath()
	tmp := path + ".tmp"
	_ = os.WriteFile(tmp, data, 0644)
	_ = os.Rename(tmp, path)
}

func applyWindowGeometry(hwnd uintptr, s *windowState) {
	if s == nil || hwnd == 0 {
		return
	}
	// HWND_TOP = 0, SWP_NOZORDER|SWP_NOACTIVATE
	const swpNoZOrder = 0x0004
	const swpNoActivate = 0x0010
	_, _, _ = procSetWindowPos.Call(hwnd, 0,
		uintptr(s.X), uintptr(s.Y),
		uintptr(s.Width), uintptr(s.Height),
		uintptr(swpNoZOrder|swpNoActivate),
	)
	if s.Maximized {
		procShowNormal.Call(hwnd, 3) // SW_MAXIMIZE
	}
}

func showNativeNotification(title, message, iconPath string) {
	notification := toast.Notification{
		AppID:   "WhatsApp Desktop",
		Title:   title,
		Message: message,
		Icon:    iconPath,
	}
	_ = notification.Push()
}

func main() {
	_, isSingle := checkSingleInstance()
	if !isSingle {
		os.Exit(0)
	}
	userDataDir := getUserDataDir()
	executablePath, _ := os.Executable()
	exeDir := filepath.Dir(executablePath)
	// Toast icon: prefer logo.png (converted from logo.avif), fall back to
	// icon.ico, and finally no icon if neither is deployed next to the exe.
	iconFullPath := filepath.Join(exeDir, "logo.png")
	if _, err := os.Stat(iconFullPath); err != nil {
		iconFullPath = filepath.Join(exeDir, "icon.ico")
		if _, err := os.Stat(iconFullPath); err != nil {
			iconFullPath = ""
		}
	}

	// Load persisted geometry
	savedState, hasSaved := loadWindowState()
	initW := 1100
	initH := 750
	center := true
	if hasSaved {
		initW = savedState.Width
		initH = savedState.Height
		center = false
		if initW < 800 {
			initW = 800
		}
		if initH < 600 {
			initH = 600
		}
	}

	opts := webview2.WebViewOptions{
		Window:    nil,
		Debug:     false,
		DataPath:  userDataDir,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  windowTitle,
			Width:  uint(initW),
			Height: uint(initH),
			IconId: 2,
			Center: center,
		},
	}

	w := webview2.NewWithOptions(opts)
	if w == nil {
		log.Fatalln("Gagal inisialisasi WebView2")
	}
	hwnd := uintptr(w.Window())
	setDarkWindowFrame(hwnd, isDarkModeEnabled())
	if hasSaved {
		applyWindowGeometry(hwnd, savedState)
	}
	// Tray icon (additive, double-click to restore)
	_ = addTrayIcon(hwnd)

	defer func() {
		saveWindowState(hwnd)
		removeTrayIcon(hwnd)
		w.Destroy()
	}()

	w.SetTitle(windowTitle)
	w.SetSize(initW, initH, webview2.HintNone)

	// Bind native notification bridge
	_ = w.Bind("sendNativeNotification", func(title, body string) {
		go showNativeNotification(title, body, iconFullPath)
	})
	// Bind for dynamic title (unread count)
	_ = w.Bind("updateWindowTitle", func(title string) {
		t := strings.TrimSpace(title)
		if t == "" {
			t = windowTitle
		}
		// Avoid overly long titles; WhatsApp title includes unread count like "(3) WhatsApp"
		if len(t) > 120 {
			t = t[:120]
		}
		w.SetTitle(t)
	})

	// Bridge page notifications to Windows toasts. No User-Agent spoofing:
	// the WebView2 runtime reports its own current Edge UA (header and JS),
	// which stays up to date unlike a frozen string.
	initScript := `
		// Native Notification bridge for Windows Desktop Toast
		(function() {
			try {
				if (window.Notification && window.Notification.permission === 'granted') {
					var OrigNotification = window.Notification;
					window.Notification = function(title, options) {
						try {
							var body = (options && options.body) || '';
							if (window.sendNativeNotification) {
								window.sendNativeNotification(title, body);
							}
						} catch (e) {}
						return new OrigNotification(title, options);
					};
					window.Notification.prototype = OrigNotification.prototype;
					window.Notification.permission = OrigNotification.permission;
					window.Notification.requestPermission = OrigNotification.requestPermission.bind(OrigNotification);
				} else {
					window.Notification = function(title, options) {
						options = options || {};
						var body = options.body || '';
						if (window.sendNativeNotification) {
							window.sendNativeNotification(title, body);
						}
						this.title = title;
						this.close = function() {};
					};
					window.Notification.permission = 'granted';
					window.Notification.requestPermission = function(callback) {
						var p = Promise.resolve('granted');
						if (typeof callback === 'function') {
							callback('granted');
						}
						return p;
					};
				}
			} catch (e) {
				window.Notification = function(title, options) {
					options = options || {};
					var body = options.body || '';
					if (window.sendNativeNotification) {
						window.sendNativeNotification(title, body);
					}
					this.title = title;
					this.close = function() {};
				};
				window.Notification.permission = 'granted';
				window.Notification.requestPermission = function(callback) {
					var p = Promise.resolve('granted');
					if (typeof callback === 'function') {
						callback('granted');
					}
					return p;
				};
			}
			// Dynamic title observer for unread badge -> native title bar
			try {
				var lastTitle = document.title;
				function pushTitle(t) {
					if (t !== lastTitle && window.updateWindowTitle) {
						lastTitle = t;
						try { window.updateWindowTitle(t); } catch(e) {}
					}
				}
				// Observe title element
				var titleEl = document.querySelector('title');
				if (titleEl) {
					new MutationObserver(function() { pushTitle(document.title); }).observe(titleEl, {childList:true});
				}
				// Fallback poll every 1s for SPA title changes
				setInterval(function(){ pushTitle(document.title); }, 1000);
				// Initial
				if (window.updateWindowTitle) {
					try { window.updateWindowTitle(document.title); } catch(e) {}
				}
			} catch(e) {}
		})();
	`

	w.Init(initScript)
	w.Navigate(appURL)
	w.Run()
}
