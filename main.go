package main

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows"
)

//go:embed camera_on.png
var cameraOnPng []byte

//go:embed camera_off.png
var cameraOffPng []byte

var (
	icoOnBytes  []byte
	icoOffBytes []byte
)

type Camera struct {
	FriendlyName string `json:"FriendlyName"`
	InstanceId   string `json:"InstanceId"`
	Status       string `json:"Status"`
}

const maxCams = 10

var (
	camItems    [maxCams]*systray.MenuItem
	camInfos    [maxCams]Camera
	mNoCamera   *systray.MenuItem
	refreshChan chan struct{}
	mu          sync.Mutex
)

func init() {
	icoOnBytes = pngToIco(cameraOnPng)
	icoOffBytes = pngToIco(cameraOffPng)
}

func setupLogging() *os.File {
	exe, err := os.Executable()
	if err != nil {
		return nil
	}
	dir := filepath.Dir(exe)
	logPath := filepath.Join(dir, "error.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(logFile)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		return logFile
	}
	return nil
}

func main() {
	logFile := setupLogging()
	if logFile != nil {
		defer logFile.Close()
	}

	log.Println("Uygulama başlatılıyor...")

	// Exe'yi yönetici olarak çalışmaya zorla (izin isteme penceresi başlangıçta bir kez açılır)
	runMeAsAdmin()

	log.Println("Yönetici yetkileri doğrulandı. Systray başlatılıyor...")
	systray.Run(onReady, onExit)
}

func isAdmin() bool {
	token := windows.GetCurrentProcessToken()
	return token.IsElevated()
}

func runMeAsAdmin() {
	if isAdmin() {
		return
	}

	log.Println("Yönetici yetkisi yok, yeniden başlatılıyor...")

	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("Uygulama yolu bulunamadı: %v", err)
	}

	dir := filepath.Dir(exe)

	verbPtr, _ := syscall.UTF16PtrFromString("runas")
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(dir)
	argPtr, _ := syscall.UTF16PtrFromString("")

	shell32 := syscall.NewLazyDLL("shell32.dll")
	shellExecuteW := shell32.NewProc("ShellExecuteW")

	ret, _, err := shellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(verbPtr)),
		uintptr(unsafe.Pointer(exePtr)),
		uintptr(unsafe.Pointer(argPtr)),
		uintptr(unsafe.Pointer(cwdPtr)),
		uintptr(1), // SW_SHOWNORMAL
	)

	// HINSTANCE <= 32 means error
	if ret <= 32 {
		messageBox("CamSwitch Hatası", fmt.Sprintf("Uygulama yönetici yetkisiyle başlatılamadı.\nHata Kodu: %d\nDetay: %v", ret, err), 0x10) // MB_ICONERROR
		log.Fatalf("Yönetici olarak yeniden başlatılamadı (Hata Kodu %d): %v", ret, err)
	}

	log.Println("Yönetici olarak yeniden başlatma tetiklendi, bu işlem sonlandırılıyor.")
	os.Exit(0)
}

func onReady() {
	systray.SetTitle("CamSwitch")
	systray.SetIcon(icoOffBytes)

	// Create camera menu slots (max 10 cameras supported dynamically)
	for i := 0; i < maxCams; i++ {
		camItems[i] = systray.AddMenuItem("", "")
		camItems[i].Hide()
	}

	mNoCamera = systray.AddMenuItem("Kamera Bulunamadı", "Sistemde bağlı kamera bulunamadı")
	mNoCamera.Disable()

	systray.AddSeparator()

	mRefresh := systray.AddMenuItem("Yenile", "Kameraları yeniden tarar")
	mQuit := systray.AddMenuItem("Çıkış", "Uygulamayı kapatır")

	refreshChan = make(chan struct{}, 1)

	// Handle clicks for each camera menu item
	for i := 0; i < maxCams; i++ {
		go func(idx int) {
			for range camItems[idx].ClickedCh {
				mu.Lock()
				cam := camInfos[idx]
				mu.Unlock()

				if cam.InstanceId == "" {
					continue
				}

				isEnabled := cam.Status == "OK"
				log.Printf("Kamera durumu değiştiriliyor: %s (Etkinleştir: %v)", cam.FriendlyName, !isEnabled)

				// Run toggle command
				err := toggleCamera(cam.InstanceId, !isEnabled)
				if err != nil {
					log.Printf("Kamera değiştirme hatası: %v", err)
					messageBox("CamSwitch Hatası", fmt.Sprintf("Kamera durumu değiştirilemedi:\n\n%v", err), 0x10) // MB_ICONERROR
				}

				// Trigger refresh immediately
				triggerRefresh()
			}
		}(i)
	}

	// Handle refresh button clicks
	go func() {
		for range mRefresh.ClickedCh {
			triggerRefresh()
		}
	}()

	// Handle quit button clicks
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	// Periodic refresh & manual refresh runner
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		// Initial scan
		triggerRefresh()

		for {
			select {
			case <-refreshChan:
				doRefresh()
			case <-ticker.C:
				doRefresh()
			}
		}
	}()
}

func onExit() {
	// Clean up if necessary
}

func triggerRefresh() {
	select {
	case refreshChan <- struct{}{}:
	default:
		// Queue is full, nothing to do
	}
}

func doRefresh() {
	cameras, err := getCameras()
	if err != nil {
		log.Printf("Kameralar taranırken hata oluştu: %v", err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	anyEnabled := false
	for _, cam := range cameras {
		if cam.Status == "OK" {
			anyEnabled = true
			break
		}
	}

	if len(cameras) == 0 {
		systray.SetIcon(icoOffBytes)
		systray.SetTooltip("Kamera Bulunamadı")
		mNoCamera.Show()
	} else {
		mNoCamera.Hide()
		if anyEnabled {
			systray.SetIcon(icoOnBytes)
			systray.SetTooltip("Kamera Etkin")
		} else {
			systray.SetIcon(icoOffBytes)
			systray.SetTooltip("Kamera Devre Dışı")
		}
	}

	for i := 0; i < maxCams; i++ {
		if i < len(cameras) {
			cam := cameras[i]
			camInfos[i] = cam

			statusStr := "Devre Dışı"
			if cam.Status == "OK" {
				statusStr = "Etkin"
				camItems[i].Check()
			} else {
				camItems[i].Uncheck()
			}

			camItems[i].SetTitle(fmt.Sprintf("%s (%s)", cam.FriendlyName, statusStr))
			camItems[i].SetTooltip(fmt.Sprintf("%s kamerasını aç/kapat", cam.FriendlyName))
			camItems[i].Show()
		} else {
			camInfos[i] = Camera{}
			camItems[i].Hide()
		}
	}
}

func getCameras() ([]Camera, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-Command", "@(Get-PnpDevice -Class Camera | Where-Object { $_.Present -eq $true } | Select-Object FriendlyName, InstanceId, Status) | ConvertTo-Json")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // Kameraları tararken de konsol ekranı yanıp sönmesin
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return []Camera{}, nil
	}

	// Try array first
	var cameras []Camera
	if err := json.Unmarshal([]byte(trimmed), &cameras); err == nil {
		return cameras, nil
	}

	// Try single object
	var single Camera
	if err := json.Unmarshal([]byte(trimmed), &single); err == nil {
		return []Camera{single}, nil
	}

	return nil, fmt.Errorf("JSON ayrıştırma hatası: %s", trimmed)
}

func toggleCamera(instanceID string, enable bool) error {
	cmdStr := fmt.Sprintf("Disable-PnpDevice -InstanceId '%s' -Confirm:$false", instanceID)
	if enable {
		cmdStr = fmt.Sprintf("Enable-PnpDevice -InstanceId '%s' -Confirm:$false", instanceID)
	}

	log.Printf("PowerShell komutu çalıştırılıyor: %s", cmdStr)
	b64Cmd := encodePowerShellCommand(cmdStr)
	cmd := exec.Command("powershell", "-NoProfile", "-EncodedCommand", b64Cmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // Konsol ekranı tamamen gizlensin
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("PowerShell exit error: %v, output: %s", err, string(output))
	}

	outStr := string(output)
	if len(outStr) > 0 {
		log.Printf("PowerShell çıktısı: %s", outStr)
	}

	// PowerShell standard cmdlets might exit with code 0 even on failure. Check output text.
	if strings.Contains(outStr, "hata") || strings.Contains(outStr, "Error") || strings.Contains(outStr, "Failure") || strings.Contains(outStr, "failed") {
		return fmt.Errorf("PowerShell işlem hatası algılandı: %s", outStr)
	}

	return nil
}

func encodePowerShellCommand(cmd string) string {
	utf16Buf := utf16.Encode([]rune(cmd))
	byteBuf := make([]byte, len(utf16Buf)*2)
	for i, r := range utf16Buf {
		byteBuf[i*2] = byte(r)
		byteBuf[i*2+1] = byte(r >> 8)
	}
	return base64.StdEncoding.EncodeToString(byteBuf)
}

func pngToIco(pngBytes []byte) []byte {
	ico := make([]byte, 22+len(pngBytes))
	// Header
	ico[0] = 0x00
	ico[1] = 0x00
	ico[2] = 0x01
	ico[3] = 0x00
	ico[4] = 0x01
	ico[5] = 0x00

	// Directory entry
	ico[6] = 0x00  // Width
	ico[7] = 0x00  // Height
	ico[8] = 0x00  // Color count
	ico[9] = 0x00  // Reserved
	ico[10] = 0x01 // Color planes
	ico[11] = 0x00
	ico[12] = 0x20 // Bits per pixel (32)
	ico[13] = 0x00

	size := uint32(len(pngBytes))
	ico[14] = byte(size)
	ico[15] = byte(size >> 8)
	ico[16] = byte(size >> 16)
	ico[17] = byte(size >> 24)

	offset := uint32(22)
	ico[18] = byte(offset)
	ico[19] = byte(offset >> 8)
	ico[20] = byte(offset >> 16)
	ico[21] = byte(offset >> 24)

	copy(ico[22:], pngBytes)
	return ico
}

func messageBox(title, text string, style uintptr) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	textPtr, _ := syscall.UTF16PtrFromString(text)
	user32 := syscall.NewLazyDLL("user32.dll")
	messageBoxW := user32.NewProc("MessageBoxW")
	messageBoxW.Call(0, uintptr(unsafe.Pointer(textPtr)), uintptr(unsafe.Pointer(titlePtr)), style)
}
