@echo off
chcp 65001 >nul
:: Check for admin rights
net session >nul 2>&1
if %errorLevel% == 0 (
    echo YÃ¶netici yetkileri doÄrulandÄ±.
) else (
    echo HATA: LÃ¼tfen bu dosyaya saÄ tÄ±klayÄ±p "YÃ¶netici olarak Ã§alÄ±ÅtÄ±r" seÃ§eneÄini seÃ§in.
    pause
    exit /b
)

set "EXE_PATH=%~dp0CamSwitch.exe"
echo CamSwitch yolu: %EXE_PATH%

powershell -NoProfile -Command "Register-ScheduledTask -TaskName 'CamSwitch' -Action (New-ScheduledTaskAction -Execute '%EXE_PATH%') -Trigger (New-ScheduledTaskTrigger -AtLogOn) -Principal (New-ScheduledTaskPrincipal -UserId \"$env:USERDOMAIN\$env:USERNAME\" -RunLevel Highest) -Settings (New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries) -Force"

if %errorLevel% == 0 (
    echo BaÅarÄ±yla baÅlangÄ±ca eklendi! BilgisayarÄ±nÄ±zÄ± her aÃ§tÄ±ÄÄ±nÄ±zda uygulama otomatik olarak baÅlayacaktÄ±r.
) else (
    echo Bir hata oluÅtu.
)
pause
