rd /q /s build
rd /q /s dist
rd /q /s output
cd ..\manpassd
go build manpassd.go
cd ..\manpassc
xcopy /y ..\manpassd\manpassd.exe .
pyinstaller --onedir --additional-hooks-dir=C:\Python27\Lib\site-packages\PyInstaller\hooks\ --hidden-import=nacl --hidden-import=cffi --noupx -w -i manpassc.ico manpassc.py
copy /Y manpassc.ico .\dist\manpassc\
copy /Y manpassd.exe .\dist\manpassc\
copy /Y msvcr120.dll .\dist\manpassc\