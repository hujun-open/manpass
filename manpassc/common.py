#!/usr/bin/env py32
# -*- coding: utf-8 -*-



#-------------------------------------------------------------------------------
# Name:        module1
# Purpose:
#
# Author:      hujun
#
# Created:     12/07/2015
# Copyright:   (c) hujun 2015
# Licence:     <your licence>
#-------------------------------------------------------------------------------
import platform
import os
import os.path
import random
import string
import time
import datetime
import sys
import shutil
import codecs
import json
import wx.lib.newevent
import M2Crypto
import passcrypto
import socket


(ManpassErrEVT, EVT_MANPASS_ERR) = wx.lib.newevent.NewEvent()
(ManpassFatalErrEVT, EVT_MANPASS_FATALERR) = wx.lib.newevent.NewEvent()
(ManpassProgressEVT, EVT_MANPASS_PROGRESS) = wx.lib.newevent.NewEvent()
(ManpassLoadingDone, EVT_MANPASS_LOAD_DONE) = wx.lib.newevent.NewEvent()
(ManpassProgressLabel, EVT_MANPASS_PROGRESS_LABEL) = wx.lib.newevent.NewEvent()




def getConfDir(uname):
    return os.path.join(getRootConfDir(),uname)


def getRootConfDir():
    myos=platform.system()
    if myos=="Windows":
        return os.path.join(os.environ["APPDATA"],"manpass")
    if myos=="Linux" or myos=="Darwin":
        cdir=os.path.join(os.environ["HOME"],".manpass")
        if not os.path.isdir(cdir):
            os.mkdir(cdir,0700)
        return cdir



def getManpassdExeName():
    myos=platform.system()
    if myos=="Windows":
        return os.path.join(cur_file_dir(),"manpassd.exe")
    if myos=="Linux" or myos=="Darwin":
        return os.path.join(cur_file_dir(),"manpassd")


def getNewPort(laddr=None):
    portList=[]
    maxport=8000
    rootconfdir=getRootConfDir()
    if os.path.isdir(rootconfdir):
        userlist=getAllImmediateDir(getRootConfDir())
    else:
        userlist=[]
    for user in userlist:
        uconf=getUserConf(user)
        if "port" in uconf:
            if uconf['port']>maxport:
                maxport=uconf['port']
    if maxport<=65534:
        rport=maxport+10
    else:
        rport= 8000
    laddr="127.0.0.1"
    while checkTCPPort(laddr,rport):
        rport+=1
        if rport==65534:
            return False
    return rport

def getAllImmediateDir(rootdir):
    rlist=[]
    for d in os.listdir(rootdir):
        if os.path.isdir(os.path.join(rootdir,d)):
            rlist.append(d)
    return rlist

def isHidden(kgroup):
    return kgroup[0:2]=="-1"

def HideIt(kgroup,hide=True):
    if hide:
        return kgroup+"-1"
    else:
        return kgroup[2:]


def getUserConf(uname):
    try:
        redirectConf=codecs.open(os.path.join(getConfDir(uname),"redirection.conf"),"r","utf-8")
        confdir=redirectConf.read().strip()
        redirectConf.close()
    except:
        confdir=getConfDir(uname)



    try:
        fp=codecs.open(os.path.join(confdir,"manpass.conf"),"r","utf-8")
        readfList=json.load(fp,"utf-8")
    except Exception as Err:
        return {}
    return readfList

def genPass(passlen=18,number=True,lowercase=True,uppercase=True,
                punction=True,ownset=None,uname=None):

    myset=""
    mypunction=string.punctuation.replace("`","") #not easy to input "`" on iOS
    if number: myset+=string.digits
    if lowercase: myset+=string.lowercase
    if uppercase: myset+=string.uppercase
    if punction: myset+=mypunction
    if ownset!=None: myset+=ownset
    mypass=""
    for i in range(passlen):
        mypass+=random.choice(myset)

    return mypass


def getLocalTime(utcs):
    #utcs is UTC time in string format "%Y-%m-%dT%H:%M:%SZ"
    #return a local time in string format "%Y-%m-%d %H:%M:%S"
    utc_t=time.strptime(utcs,"%Y-%m-%dT%H:%M:%SZ")
    utc_stamp=time.mktime(utc_t)

    if utc_t.tm_isdst==0:
        local_stamp=utc_stamp-time.timezone
    else:
        local_stamp=utc_stamp-time.altzone
    local_t=time.localtime(local_stamp)
    return time.strftime("%Y-%m-%d %H:%M:%S",local_t)


def cur_file_dir():
    #获取脚本路径
    MYOS=platform.system()
    if MYOS == 'Linux':
        path = sys.path[0]
    elif MYOS == 'Windows':
        return os.path.dirname(os.path.abspath(sys.argv[0]))
    else:
        if sys.argv[0].find('/') != -1:
            path = sys.argv[0]
        else:
            path = sys.path[0]
    if isinstance(path,str):
        path=path.decode(sys.getfilesystemencoding())

    #判断为脚本文件还是py2exe编译后的文件，如果是脚本文件，则返回的是脚本的目录，如果是编译后的文件，则返回的是编译后的文件路径
    if os.path.isdir(path):
        return path
    elif os.path.isfile(path):
        return os.path.dirname(path)


def copyConfigFiles(src_dir,dst_dir,uname):
    shutil.copy2(os.path.join(src_dir,"ca.cert"),dst_dir)
    shutil.copy2(os.path.join(src_dir,"ee.cert"),dst_dir)
    shutil.copy2(os.path.join(src_dir,"ca.key"),dst_dir)
    shutil.copy2(os.path.join(src_dir,"ee.key"),dst_dir)
    shutil.copy2(os.path.join(src_dir,uname+u".db"),dst_dir)

def reEncryptCertFiles(confpath,oldpass,newpass):
    #re-encrypt ca.cert and ee.key with new pass
    def old(*args):
        return oldpass
    def new(*args):
        return newpass
    pkey=M2Crypto.RSA.load_key(os.path.join(confpath,"ee.key"),old)
    pkey.save_pem(os.path.join(confpath,"ee.key"),'aes_128_cbc',new)

    caf=open(os.path.join(confpath,"ca.cert"),"r")
    buf=caf.read()
    caf.close()
    clearca=passcrypto.DecryptMeBase32(buf,oldpass)
    caf=open(os.path.join(confpath,"ca.cert"),"w")
    caf.write(passcrypto.EncryptMeBase32(clearca,newpass))
    caf.close()


def checkTCPPort(address, port):
    s=socket.socket()
    try:
        s.connect((address,port))
        s.close()
        return True
    except socket.error,e:
        s.close()
        return False




if __name__ == '__main__':
    print checkTCPPort("127.0.0.1",8030)
