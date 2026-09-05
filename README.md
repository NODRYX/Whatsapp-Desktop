# WhatsApp Webview Desktop — Fork

> **🇮🇩 Fork dari [Adytm404/whatsapp-web.view](https://github.com/Adytm404/whatsapp-web.view)** — fork terawat dengan optimasi icon & performa WebView2.
> **🇬🇧 Fork of [Adytm404/whatsapp-web.view](https://github.com/Adytm404/whatsapp-web.view)** — maintained fork with icon & WebView2 performance tuning.

<p align="center">
  <a href="#indonesia">🇮🇩 Indonesia</a> &nbsp;|&nbsp; <a href="#english">🇬🇧 English</a> &nbsp;|&nbsp; <a href="#perbedaan-dari-upstream--whats-changed">🔀 Perbedaan / Changes</a>
</p>

<p align="center">
  <sub>
    <b>Upstream:</b> <a href="https://github.com/Adytm404/whatsapp-web.view">Adytm404/whatsapp-web.view</a> &nbsp;•&nbsp;
    <b>Fork:</b> <a href="https://github.com/NODRYX/Whatsapp-Desktop">NODRYX/Whatsapp-Desktop</a> &nbsp;•&nbsp;
    <b>Base:</b> <code>508d8ad</code> &nbsp;•&nbsp;
    <b>Fork commit:</b> <code>5bb1f05</code> — <code>feat: logo.avif as master icon plus WebView2 low-resource tuning</code>
  </sub>
</p>

---

<a id="indonesia"></a>
## 🇮🇩 Indonesia

### Tentang Fork

Proyek ini adalah **fork** dari [Adytm404/whatsapp-web.view](https://github.com/Adytm404/whatsapp-web.view). Semua kredit untuk ide awal, struktur `main.go`, dan wrapper WebView2 tetap milik upstream. Fork ini hanya menambahkan **satu commit pembeda** (`5bb1f05`) di atas `508d8ad` untuk optimasi branding dan konsumsi resource — tanpa mengubah alur login, penyimpanan sesi, atau kompatibilitas WhatsApp Web.

> Jika kamu mencari versi original tanpa modifikasi, gunakan upstream: https://github.com/Adytm404/whatsapp-web.view

### Fitur

- Window WebView2 native Windows (bukan Electron) — `main.go:130-207`
- Sesi WhatsApp persisten lintas restart (cookies, LocalStorage, IndexedDB, service-worker di profil khusus)
- Title bar & frame gelap (DWM immersive dark mode)
- Bridge notifikasi toast Windows
- Akses kamera & mikrofon untuk voice/video call
- Proteksi single-instance (mutex `WhatsAppDesktopSingleInstanceMutex`)
- Dukungan High-DPI
- Executable native kecil
- **Baru di fork:** Ikon multi-resolusi dari `logo.avif` (16/32/48/64/128/256) + `logo.png` 256px untuk toast — `gen_icon.py:14-74`
- **Baru di fork:** Tuning hemat resource — menonaktifkan 6 fitur WebView2 yang tidak dipakai + filter request telemetry — `vendor/github.com/jchv/go-webview2/webview.go:132-207`

### Persyaratan

- Windows 10 atau lebih baru
- Microsoft Edge WebView2 Runtime (akan diunduh otomatis jika belum ada)
- Akun WhatsApp yang sudah pairing dengan WhatsApp Web

### Unduhan

- **Upstream (original):** `WhatsApp.exe` dari [Adytm404/whatsapp-web.view Releases](https://github.com/Adytm404/whatsapp-web.view/releases)
- **Fork ini:** build dari branch `main` repo ini (`NODRYX/Whatsapp-Desktop`). Jalankan exe, scan QR, izinkan kamera/mikrofon saat diminta Windows — sesi tersimpan otomatis.

### Data Sesi

Data profil disimpan di:

```text
%APPDATA%\WhatsAppDesktopLight\UserData
```

Jangan hapus folder ini jika sesi login ingin tetap tersimpan. Menutup aplikasi tidak menghapus data sesi.

### Build dari Source

Butuh Go + compiler C untuk Windows (mis. `mingw`), dan **opsional `ffmpeg`** jika ingin regenerate ikon.

```powershell
# 1. (Opsional) regenerate ikon dari logo.avif — butuh ffmpeg di PATH
python gen_icon.py
# -> menghasilkan logo.png + icon.ico (6 ukuran, PNG-compressed)

# 2. Embed manifest + ikon ke resource Windows
rsrc -manifest app.manifest -ico icon.ico -o rsrc.syso

# 3. Build
go mod download
go build -ldflags="-H windowsgui -s -w" -o WhatsApp.exe .
```

Tanpa `ffmpeg`, `logo.png`/`icon.ico` yang sudah ada di repo tetap dipakai untuk build.

### File Proyek

- `main.go` — window WebView2, profil persisten, dark frame, notifikasi (`main.go:26-207`)
- `logo.avif` — master logo 740x740, tidak dipakai langsung, hanya source untuk downscale
- `logo.png` — ikon toast 256px hasil generate dari `logo.avif` (ditaruh sebelah exe)
- `gen_icon.py` — regenerate `logo.png` + `icon.ico` multi-size dari `logo.avif` (butuh ffmpeg, lanczos downscale)
- `app.manifest` — manifest Windows DPI & aplikasi
- `resource.rc` — definisi resource ikon & manifest Windows
- `icon.ico` — ikon aplikasi multi-size 16–256 (generated, entri PNG-compressed)
- `rsrc.syso` — resource ter-embed (manifest + ikon)

### Privasi

Aplikasi memuat WhatsApp Web secara langsung. Data chat & status autentikasi ditangani oleh WhatsApp Web dan disimpan lokal di profil WebView2 di atas. Proyek ini tidak terafiliasi dengan WhatsApp atau Meta.

---

<a id="english"></a>
## 🇬🇧 English

### About This Fork

This project is a **fork** of [Adytm404/whatsapp-web.view](https://github.com/Adytm404/whatsapp-web.view). All credit for the original idea, `main.go` structure, and WebView2 wrapper goes to upstream. This fork adds **a single divergent commit** (`5bb1f05`) on top of `508d8ad` for branding and resource-usage improvements — without changing login flow, session storage, or WhatsApp Web compatibility.

> For the unmodified original, use upstream: https://github.com/Adytm404/whatsapp-web.view

### Features

- Native Windows WebView2 window (no Electron) — `main.go:130-207`
- Persistent WhatsApp session across restarts (cookies, LocalStorage, IndexedDB, service workers in a dedicated profile)
- Dark Windows title bar & frame (DWM immersive dark mode)
- Windows toast notification bridge
- Camera & microphone access for voice/video calls
- Single-instance protection (mutex `WhatsAppDesktopSingleInstanceMutex`)
- High-DPI support
- Small native executable
- **New in fork:** Multi-resolution icons from `logo.avif` (16/32/48/64/128/256) + 256px `logo.png` for toasts — `gen_icon.py:14-74`
- **New in fork:** Low-resource tuning — 6 unused WebView2 features disabled + telemetry request filter — `vendor/github.com/jchv/go-webview2/webview.go:132-207`

### Requirements

- Windows 10 or newer
- Microsoft Edge WebView2 Runtime (auto-downloaded if missing)
- WhatsApp account paired with WhatsApp Web

### Download

- **Upstream (original):** `WhatsApp.exe` from [Adytm404/whatsapp-web.view Releases](https://github.com/Adytm404/whatsapp-web.view/releases)
- **This fork:** build from `main` of this repo (`NODRYX/Whatsapp-Desktop`). Run the exe, scan the QR code, allow camera/microphone when prompted — session is saved automatically.

### Session Data

Profile data is stored at:

```text
%APPDATA%\WhatsAppDesktopLight\UserData
```

Do not delete this folder if you want to keep the login session. Closing the app does not clear session data.

### Build From Source

Requires Go + a Windows C compiler (e.g. `mingw`), and **optionally `ffmpeg`** to regenerate icons.

```powershell
# 1. (Optional) regenerate icons from logo.avif — requires ffmpeg on PATH
python gen_icon.py
# -> outputs logo.png + icon.ico (6 sizes, PNG-compressed)

# 2. Embed manifest + icon into Windows resource
rsrc -manifest app.manifest -ico icon.ico -o rsrc.syso

# 3. Build
go mod download
go build -ldflags="-H windowsgui -s -w" -o WhatsApp.exe .
```

Without `ffmpeg`, the committed `logo.png`/`icon.ico` are used as-is.

### Project Files

- `main.go` — WebView2 window, persistent profile, dark frame, notifications (`main.go:26-207`)
- `logo.avif` — 740×740 master logo, never used directly, only downscaled
- `logo.png` — 256px toast icon generated from `logo.avif` (shipped next to the exe)
- `gen_icon.py` — regenerates `logo.png` + multi-size `icon.ico` from `logo.avif` (requires ffmpeg, lanczos downscale)
- `app.manifest` — Windows DPI & application manifest
- `resource.rc` — Windows icon & manifest resource definitions
- `icon.ico` — application icon, multi-size 16–256 (generated, PNG-compressed entries)
- `rsrc.syso` — embedded resource (manifest + icon)

### Privacy

This app loads WhatsApp Web directly. Chat data and authentication state are handled by WhatsApp Web and stored locally in the WebView2 profile above. This project is not affiliated with WhatsApp or Meta.

---

<a id="perbedaan-dari-upstream--whats-changed"></a>
## 🔀 Perbedaan dari Upstream / What's Changed

> **ID:** Satu-satunya pembeda fork ini vs upstream adalah commit [`5bb1f05` `feat: logo.avif as master icon plus WebView2 low-resource tuning`](https://github.com/NODRYX/Whatsapp-Desktop/commit/5bb1f05) di atas `508d8ad`. Ringkasan di bawah dikelompokkan per tema agar mudah di-review.
> **EN:** The only difference between this fork and upstream is commit [`5bb1f05` `feat: logo.avif as master icon plus WebView2 low-resource tuning`](https://github.com/NODRYX/Whatsapp-Desktop/commit/5bb1f05) on top of `508d8ad`. Grouped by theme for easy review.

| # | Area | 🇮🇩 Indonesia (apa yang berubah) | 🇬🇧 English (what changed) | File |
|---|------|--------------------------------|---------------------------|------|
| 1 | **Icon pipeline** | `gen_icon.py` lama generate lingkaran hijau 32×32 prosedural (struct/BMP). Sekarang `logo.avif` 740×740 jadi single source of truth; `gen_icon.py:33-48` pakai `ffmpeg -vf scale=SIZE:SIZE:flags=lanczos` downscale ke `SIZES=[16,32,48,64,128,256]`; hasil `logo.png` 256px untuk toast + `icon.ico` 38 KB multi-size (entri PNG-compressed, Vista+) — sebelumnya 4 KB single-size. `rsrc.syso` di-regenerate. | Old `gen_icon.py` procedurally generated a 32×32 green circle (struct/BMP). Now `logo.avif` 740×740 is the single source; `gen_icon.py:33-48` uses `ffmpeg -vf scale=SIZE:SIZE:flags=lanczos` to downscale to `SIZES=[16,32,48,64,128,256]`; outputs 256px `logo.png` for toasts + 38 KB multi-size `icon.ico` (PNG-compressed entries) vs 4 KB single-size before. `rsrc.syso` regenerated. | `gen_icon.py`, `logo.avif` (baru/new), `logo.png` (baru/new), `icon.ico`, `rsrc.syso` |
| 2 | **Toast icon fallback** | Sebelumnya `iconFullPath = exeDir/icon.ico` fixed. Sekarang prefer `logo.png` → fallback `icon.ico` → `""` jika tidak ada (`main.go:118-128`, cek `os.Stat`). | Previously fixed `iconFullPath = exeDir/icon.ico`. Now prefers `logo.png` → falls back to `icon.ico` → `""` if neither exists (`main.go:118-128`, `os.Stat` check). | `main.go:118-128` |
| 3 | **Hapus User-Agent spoofing** | Hapus `const userAgent = "Chrome/133..."` (`main.go:27` di upstream) + `Object.defineProperty(navigator, 'userAgent'/'appVersion')`. WebView2 sekarang melaporkan Edge UA native (header & JS) sehingga selalu up-to-date, tidak frozen. Komentar di `main.go:161-163` menjelaskan alasan. | Removed `const userAgent = "Chrome/133..."` (`main.go:27` in upstream) + `Object.defineProperty(navigator, 'userAgent'/'appVersion')`. WebView2 now reports its native Edge UA (header & JS) so it stays current instead of frozen. See comment `main.go:161-163`. | `main.go:27`, `main.go:161-163` |
| 4 | **Notification bridge** | Polyfill lama `window.Notification = function(...)` sederhana (`onclick/onclose/... = null`). Sekarang: jika `Notification.permission==='granted'` bungkus `OrigNotification` (delegate `prototype` & `requestPermission.bind`), kirim `sendNativeNotification` juga; fallback polyfill tambah `this.close=function(){}` + `requestPermission` return `Promise.resolve('granted')` + support callback (`main.go:164-202`). | Old polyfill was a simple `window.Notification = function(...)` (`onclick/onclose/... = null`). Now: if `Notification.permission==='granted'` wraps `OrigNotification` (delegates `prototype` & `requestPermission.bind`) and also calls `sendNativeNotification`; fallback polyfill adds `this.close=function(){}` + `requestPermission` returns `Promise.resolve('granted')` with callback support (`main.go:164-202`). | `main.go:164-202` |
| 5 | **Low-resource tuning** | `webview.go:132-207` — `NewWithOptions` kini panggil `applyLowResourceSettings()` yang disable 6 fitur tidak dipakai via `PutIsStatusBarEnabled(false)`, `PutAreDefaultScriptDialogsEnabled`, `PutIsPinchZoomEnabled`, `PutIsSwipeNavigationEnabled`, `PutIsPasswordAutosaveEnabled`, `PutIsGeneralAutofillEnabled` (butuh 4 method baru di `ICoreWebViewSettings.go:343-395`) + `setupResourceFilter()` yang block `crashlog/analytics/telemetry/csp-report/.webmanifest` via `CreateWebResourceResponse(204)` + `AddWebResourceRequestedFilter("*", COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)` (allow-by-default, chat/media/login tidak terpengaruh). | `webview.go:132-207` — `NewWithOptions` now calls `applyLowResourceSettings()` disabling 6 unused features via `PutIsStatusBarEnabled(false)`, `PutAreDefaultScriptDialogsEnabled`, `PutIsPinchZoomEnabled`, `PutIsSwipeNavigationEnabled`, `PutIsPasswordAutosaveEnabled`, `PutIsGeneralAutofillEnabled` (4 new methods in `ICoreWebViewSettings.go:343-395`) + `setupResourceFilter()` blocking `crashlog/analytics/telemetry/csp-report/.webmanifest` via `CreateWebResourceResponse(204)` + `AddWebResourceRequestedFilter("*", COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)` (allow-by-default, chat/media/login unaffected). | `vendor/github.com/jchv/go-webview2/webview.go:132-207`, `vendor/github.com/jchv/go-webview2/pkg/edge/ICoreWebViewSettings.go:343-395` |
| 6 | **Stabilitas** | `chromium.go:281-284` — `WebResourceRequested` sebelumnya `log.Fatal(err)` (crash app jika COM error). Sekarang `log.Printf` dan return 0 (request lanjut tanpa filter). | `chromium.go:281-284` — `WebResourceRequested` previously `log.Fatal(err)` (crashed app on COM error). Now `log.Printf` and returns 0 (request continues unfiltered). | `vendor/github.com/jchv/go-webview2/pkg/edge/chromium.go:281-284` |

**Cara verifikasi / How to verify:**

```powershell
git log --oneline origin/main..HEAD
# -> 5bb1f05 feat: logo.avif as master icon plus WebView2 low-resource tuning

git diff origin/main..HEAD --stat
# -> 10 files changed, 241 insertions(+), 68 deletions(-)

git show 5bb1f05 --stat
```

### Lisensi / License

**ID:** Belum ada lisensi yang dideklarasikan di upstream maupun fork ini. Hak cipta tetap milik kontributor masing-masing. Jika upstream menambahkan lisensi, fork ini akan mengikuti.

**EN:** No license has been declared yet in upstream or this fork. Copyright remains with respective contributors. If upstream adds a license, this fork will follow.

---

<p align="center">
  <sub>
    🇮🇩 Dibuat dengan Go + WebView2 &nbsp;|&nbsp; 🇬🇧 Built with Go + WebView2 &nbsp;•&nbsp;
    <a href="#indonesia">🇮🇩 Indonesia</a> · <a href="#english">🇬🇧 English</a> · <a href="#perbedaan-dari-upstream--whats-changed">🔀 Changes</a>
  </sub>
</p>
