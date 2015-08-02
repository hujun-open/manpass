rd /q /s build
rd /q /s dist
rd /q /s output
pyinstaller manpassc.spec
xcopy /Y /E /I C:\Python27\Lib\site-packages\PyNaCl-0.3.0-py2.7-win32.egg\nacl .\dist\manpassc\nacl
copy /Y manpassc.ico .\dist\manpassc\
copy /Y D:\hujun\manpass\src\manpassd\manpassd.exe .\dist\manpassc