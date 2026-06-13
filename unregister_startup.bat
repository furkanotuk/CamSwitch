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

powershell -NoProfile -Command "Unregister-ScheduledTask -TaskName 'CamSwitch' -Confirm:$false"

if %errorLevel% == 0 (
    echo BaÅarÄ±yla baÅlangÄ±Ã§tan kaldÄ±rÄ±ldÄ±.
) else (
    echo GÃ¶rev bulunamadÄ± veya kaldÄ±rÄ±lÄ±rken hata oluÅtu.
)
pause
