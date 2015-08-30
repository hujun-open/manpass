# Manpass
Manpass is an opensource, secure password manager, it has following features:
*	Support Windows/Linux/OSX
* Only need to remember one masterpassword
*	Use proven highly secure crypto algorithm/implementation
*	Double-click to copy password to clip-board, and automatically clear it from clipboard after 20 seconds.
*	Show password as QRCode image for mobile device use.
*	Instant search among saved credentials (support Pinyin)
*	Password history save the past passwords
*	Passwords are stored in encrypted form within database, so that the database could reside in cloud storage (e.g. Dropbox) for redundancy and multi-clients synchronization.
*	Support multiple Users
# Usage
see help.htm or "Help" in popup menu

# To Build/Run Source Code
manpassc is written in python, manpassd is written in Go.
You need following enviroment:
* Python 2.7.10
  * [wxpython 3.0.2.0](http://www.wxpython.org/)
  * [pynacl 0.3.0](https://github.com/pyca/pynacl)
  * [pyscrypt 1.6.2](https://github.com/ricmoo/pyscrypt)
  * [M2Crypto 0.22.3](https://github.com/martinpaljak/M2Crypto)
  * [Pillow 2.9.0](https://github.com/python-pillow/Pillow)
  * [qrcode 5.1](https://github.com/lincolnloop/python-qrcode)
* Go 1.4.2
  * [go-sqlite3](https://github.com/mattn/go-sqlite3)



# Information For Geek
Manpass has client/server architecture, manpassd is the server which basically is a daemon for credential database, all credentials are stored in manpassd's DB; manpassc is the GUI client talk to manpassd via RESTful API over HTTPS.

* All credentials are stored in encrypted form in manpassd's DB.

* Encryption/decryption is done by manpassc, via [NaCl libray](http://nacl.cr.yp.to/)

* Each credential is encrypted by a key derived from masterpassword via [scrypt algorithm](https://en.wikipedia.org/wiki/Scrypt). the derived key is unique for each credential.
* manpassc and manpassd will authenticate each other via certificate, the CA/EE certificate are generated upon new user creation, each user has its own CA/EE certificates.
the RSA key and CA certificate are encrypted by user's masterpassword.
