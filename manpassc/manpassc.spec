# -*- mode: python -*-
a = Analysis(['manpassc.py'],
             pathex=['D:\\hujun\\manpass\\src\\manpassc'],
             hiddenimports=['nacl','M2Crypto'],
             hookspath=None,
             runtime_hooks=None)
pyz = PYZ(a.pure)
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
