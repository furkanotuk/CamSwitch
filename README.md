# CamSwitch 📷

Choose your language / Dil seçin:
* [Türkçe (TR)](#türkçe)
* [English (EN)](#english)

---

## Türkçe

CamSwitch, Windows bilgisayarınızda bağlı olan kameraları sistem çekmecesinden (system tray - sağ alt köşe) tek tıkla kolayca etkinleştirip devre dışı bırakmanızı (açıp kapatmanızı) sağlayan pratik bir araçtır.

### Özellikler

* **Otomatik Algılama:** Bilgisayara bağlı olan kameraları otomatik olarak tarar ve listeler (ekstra ayar gerektirmez).
* **Durum Göstergesi:** Kameranız aktifken yeşil, tüm kameralar devre dışı veya kamera bağlı değilken kırmızı simge gösterilir.
* **Kolay Aç/Kapat:** Sistem çekmecesindeki simgeye sağ tıklayıp kameranızı seçerek durumunu değiştirebilirsiniz.
* **Başlangıçta Çalıştırma:** Bilgisayar açıldığında otomatik olarak arka planda başlaması için hazır betikler içerir.

### Nasıl Kullanılır?

1. **Uygulamayı Çalıştırın:** `CamSwitch.exe` dosyasına çift tıklayarak uygulamayı başlatın.
   * *Not:* Kameraları donanımsal düzeyde yönetmek için Windows yönetici izinleri gereklidir. Uygulama açılırken otomatik olarak yönetici yetkisi (UAC) isteyecektir.
2. **Kamerayı Yönetin:** Sağ alttaki sistem çekmecesinde bulunan kamera simgesine sağ tıklayın. Listelenen kameralardan kapatmak veya açmak istediğiniz kameranın üzerine tıklayın.

### Başlangıçta Otomatik Çalıştırma (Opsiyonel)

Uygulamanın Windows her başladığında otomatik olarak arka planda çalışmasını istiyorsanız:

1. **register_startup.bat** dosyasına sağ tıklayıp **"Yönetici olarak çalıştır"** seçeneğini seçin.
2. Bu işlem, uygulamayı Windows Görev Zamanlayıcı'ya yüksek yetkilerle otomatik başlatılacak şekilde ekler.
3. Otomatik başlatmayı kaldırmak isterseniz **unregister_startup.bat** dosyasına sağ tıklayıp **"Yönetici olarak çalıştır"** diyebilirsiniz.

### Geliştiriciler İçin Derleme (Build)

Projeyi kaynak kodundan kendiniz derlemek isterseniz:

```bash
go build -ldflags "-H windowsgui" -o CamSwitch.exe main.go
```

*(Not: `-H windowsgui` bayrağı, uygulamanın arka planda çalışırken gereksiz bir komut satırı (CMD) penceresi açmasını engeller.)*

---

## English

CamSwitch is a lightweight utility that allows you to easily enable or disable (turn on/off) your connected webcams directly from the Windows system tray (bottom-right corner) with a single click.

### Features

* **Auto-Detection:** Automatically scans and lists all connected cameras (no manual configuration required).
* **Status Indicators:** Show a green icon when a camera is active/enabled and a red icon when all cameras are disabled or none is found.
* **Quick Toggle:** Right-click the system tray icon, select a camera, and click to toggle its status.
* **Run at Startup:** Includes scripts to register the application to run automatically on Windows boot.

### How to Use

1. **Run the App:** Double-click `CamSwitch.exe` to start the application.
   * *Note:* Administrator privileges are required to toggle hardware devices. The application will automatically prompt for UAC elevation on startup.
2. **Toggle Cameras:** Right-click the camera icon in your system tray. Click on any listed camera to toggle its state (Enabled/Disabled).

### Run at Startup (Optional)

If you want the application to start automatically in the background whenever Windows boots:

1. Right-click **register_startup.bat** and select **"Run as administrator"**.
2. This creates a Windows Scheduled Task configured to run the application with highest privileges at logon.
3. To remove the startup task, right-click **unregister_startup.bat** and select **"Run as administrator"**.

### Building from Source

To compile the application from source code:

```bash
go build -ldflags "-H windowsgui" -o CamSwitch.exe main.go
```

*(Note: The `-H windowsgui` flag prevents an empty command prompt window from popping up when the application is running in the background.)*
