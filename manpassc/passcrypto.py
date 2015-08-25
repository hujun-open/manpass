import pyscrypt
import nacl
import nacl.secret
import nacl.utils
import nacl.encoding
import nacl.hash

KeySize   = 32
NonceSize = 24
SaltSize  = 16
Scrypt_N  = 1024
Scrypt_r  = 1
Scrypt_p  = 1

def GenerateSalt():
    return nacl.utils.random(SaltSize)

def GenerateEncKey(passwd, salt):
    enc_key=pyscrypt.hash(password=passwd,
        salt = salt,
        N=Scrypt_N,
        r=Scrypt_r,
        p=Scrypt_p,
        dkLen=KeySize)
    return enc_key

def GenerateNonce():
    return nacl.utils.random(NonceSize)


def EncryptWithoutSalt(msg,skey):
    nonce=GenerateNonce()
    box = nacl.secret.SecretBox(skey)
    cipher_txt=box.encrypt(msg,nonce)
    return cipher_txt

def DecryptWithoutSalt(cipher,skey):
    box = nacl.secret.SecretBox(skey)
    return box.decrypt(cipher)

def EncryptMe(msg, passwd):
    salt=GenerateSalt()
    skey=GenerateEncKey(passwd,salt)
    rstr=EncryptWithoutSalt(msg,skey)
    return salt+rstr

def DecryptMe(cipher,passwd):
    salt=cipher[:SaltSize]
    skey=GenerateEncKey(passwd,salt)
    return DecryptWithoutSalt(cipher[SaltSize:],skey)

def EncryptMeBase32(msg,passwd):
    cipher=EncryptMe(msg,passwd)
    enc=nacl.encoding.Base32Encoder()
    return enc.encode(cipher)

def EncryptWithoutSaltBase32(msg,skey,salt):
    cipher=EncryptWithoutSalt(msg,skey)
    enc=nacl.encoding.Base32Encoder()
    return enc.encode(salt+cipher)

def DecryptMeBase32(cipher,passwd):
    enc=nacl.encoding.Base32Encoder()
    return DecryptMe(enc.decode(cipher),passwd)



def HashMsg(msg,skey):
    h=nacl.hash.sha256(skey+msg)
    return h+msg

def VerifyHash(msg,skey):
    if len(msg)<=64:
        return False
    h=msg[0:64]
    nh=nacl.hash.sha256(skey+msg[64:])
    if nh==h:
        return msg[64:]
    else:
        return False



if __name__ == '__main__':
    s=HashMsg("xixixi","z123")
    print VerifyHash(s+"1","z123")
##    clear="hello world!"
##    passwd="alu123"
##    s=EncryptMeBase32(clear,passwd)
##    if DecryptMeBase32(s,passwd) != clear:
##        print "error!"
##    else:
##        print "ok"
##    inf=open("d:\\temp\\tls\\root.cer","rb")
##    clear=inf.read()
##    inf.close()
##    outf=open("d:\\temp\\tls\\root.encrypted","w")
##    outf.write(EncryptMeBase32(clear,"alu123"))
##    outf.close()
