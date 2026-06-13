# CamSwitch 📷

CamSwitch, Windows bilgisayarınızda bağlı olan kameraları sistem çekmecesinden (system tray - sağ alt köşe) tek tıkla kolayca etkinleştirip devre dışı bırakmanızı (açıp kapatmanızı) sağlayan pratik bir araçtır.

## Özellikler

* **Otomatik Algılama:** Bilgisayara bağlı olan kameraları otomatik olarak tarar ve listeler (ekstra ayar gerektirmez).
* **Durum Göstergesi:** Kameranız aktifken yeşil, tüm kameralar devre dışıyken veya kamera bağlı değilken kırmızı simge gösterilir.
* **Kolay Aç/Kapat:** Sistem çekmecesindeki simgeye sağ tıklayıp kameranızı seçerek durumunu değiştirebilirsiniz.

## Nasıl Kullanılır?

1. **Uygulamayı Çalıştırın:** `CamSwitch.exe` dosyasına çift tıklayarak uygulamayı başlatın.
   * *Not:* Kameraları donanımsal düzeyde yönetmek için Windows yönetici izinleri gereklidir. Uygulama açılırken otomatik olarak yönetici yetkisi (UAC) isteyecektir.
2. **Kamerayı Yönetin:** Sağ alttaki sistem çekmecesinde bulunan kamera simgesine sağ tıklayın. Listelenen kameralardan kapatmak veya açmak istediğiniz kameranın üzerine tıklayın.

## Başlangıçta Otomatik Çalıştırma (Opsiyonel)

Uygulamanın Windows her başladığında otomatik olarak arka planda çalışmasını istiyorsanız:

1. **register_startup.bat** dosyasına sağ tıklayıp **"Yönetici olarak çalıştır"** seçeneğini seçin.
2. Bu işlem, uygulamayı Windows Görev Zamanlayıcı'ya yüksek yetkilerle otomatik başlatılacak şekilde ekler.
3. Otomatik başlatmayı kaldırmak isterseniz **unregister_startup.bat** dosyasına sağ tıklayıp **"Yönetici olarak çalıştır"** diyebilirsiniz.

## Geliştiriciler İçin Derleme (Build)

Projeyi kaynak kodundan kendiniz derlemek isterseniz, Go (Golang) yüklü bilgisayarınızda aşağıdaki komutu çalıştırabilirsiniz:

```bash
go build -ldflags "-H windowsgui" -o CamSwitch.exe main.go
```

*(Not: `-H windowsgui` bayrağı, uygulamanın arka planda çalışırken gereksiz bir komut satırı (CMD) penceresi açmasını engeller.)*
