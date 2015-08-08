# -*- mode: python -*-

block_cipher = None


a = Analysis(['manpassc.py'],
             pathex=['D:\\hujun\\Dropbox\\manpass\\src\\manpassc'],
             hiddenimports=['nacl', 'cffi'],
             hookspath=['C:\\Python27\\Lib\\site-packages\\PyInstaller\\hooks\\'],
             runtime_hooks=None,
             excludes=None,
             cipher=block_cipher)
pyz = PYZ(a.pure,
             cipher=block_cipher)
exe = EXE(pyz,
          a.scripts,
          exclude_binaries=True,
          name='manpassc.exe',
          debug=False,
          strip=None,
          upx=False,
          console=False , icon='manpassc.ico')
coll = COLLECT(exe,
               a.binaries,
               a.zipfiles,
               a.datas,
               strip=None,
               upx=False,
               name='manpassc')
