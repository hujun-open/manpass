# -*- mode: python -*-

block_cipher = None


a = Analysis(['manpassc.py'],
             pathex=['/Users/hujun/Dropbox/manpass/src/manpassc'],
             hiddenimports=['M2Crypto', 'nacl', 'cffi'],
             hookspath=['/usr/local/lib/python2.7/site-packages/PyInstaller/hooks'],
             runtime_hooks=None,
             excludes=None,
             cipher=block_cipher)
pyz = PYZ(a.pure,
             cipher=block_cipher)
exe = EXE(pyz,
          a.scripts,
          exclude_binaries=True,
          name='manpassc',
          debug=True,
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
app = BUNDLE(coll,
             name='manpassc.app',
             icon='manpassc.ico',
             bundle_identifier=None)
