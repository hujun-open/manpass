#!/usr/bin/env python
# -*- coding: UTF-8 -*-


#
# This class provides a convient way to setup options for wxpython applications
# it uses a list of dict to represent a list of configurable options, each dict
# includes following keys: "name","type","default"
#

import wx
import wx.propgrid as wxpg
import json
import codecs
import common
# begin wxGlade: dependencies
import gettext
import os.path
# end wxGlade

# begin wxGlade: extracode
# end wxGlade
_ = wx.GetTranslation

class OptionDiag(wx.Dialog):
    def __init__(self, parent,confFile,expectList,uname):
    #confFile is a filename doesn't include path
    #expectList is a list of confGroup, each confGroup is a two element tuple: (group_name, conf_list)
    # conf_list is a list of conf, each conf is a two element tuple(conf_name, conf_dict)
    # each conf_dict inlcude following keys: "desc","value","type"
    # "type" is one of following: string, directory, int, float

        wx.Dialog.__init__(self, parent)
        self.grid=wxpg.PropertyGrid(self,wx.ID_ANY)
        self.uname=uname
        self.confFName=confFile
        self.confPath=self.confFName
        self.confList=expectList
        self.loadConfFile()
        self.refereshView()

        self.__set_properties()
        self.__do_layout()
        self.Centre()
        okb= self.bsizer.GetItem(1).GetWindow()
        canb=self.bsizer.GetItem(2).GetWindow()
        self.Bind(wx.EVT_BUTTON,self.OnOK,okb)
        self.Bind(wx.EVT_BUTTON,self.OnCancel,canb)

        # end wxGlade

    def __set_properties(self):
        # begin wxGlade: MainPannel.__set_properties
        self.SetTitle(_("Generate a password"))
        #self.SetWindowStyle(wx.BORDER_DEFAULT)
        # end wxGlade

    def __do_layout(self):
        # begin wxGlade: MainPannel.__do_layout
        sizer_v = wx.BoxSizer(wx.VERTICAL)
        sizer_v.Add(self.grid, 0, wx.ALL|wx.ALIGN_CENTER_HORIZONTAL|wx.EXPAND, 5)
        self.bsizer=self.CreateButtonSizer(wx.OK|wx.CANCEL)
        sizer_v.Add(self.bsizer, 0, wx.ALL|wx.ALIGN_CENTER_HORIZONTAL, 5)
        self.SetSizer(sizer_v)
        sizer_v.Fit(self)
        self.Layout()
        # end wxGlade


    def loadConfFile(self):
        #load a redirect.conf from default conf directory first, if it exisits,
        #it contains the real directory path to the main configuration file
        #if it doesn't exist, then system will load main config from default dir
        try:
            redirectConf=codecs.open(os.path.join(common.getConfDir(self.uname),"redirection.conf"),"r","utf-8")
            confdir=redirectConf.read()
            redirectConf.close()
        except:
            confdir=common.getConfDir(self.uname)

        self.confPath=os.path.join(confdir,self.confFName)

        try:
            fp=codecs.open(self.confPath,"r","utf-8")
            readfList=json.load(fp,"utf-8")
        except Exception as Err:
            readfList=[]
        for cgroup in self.confList:
            for conf in cgroup[1]:
                if conf[0] in readfList:
                    conf[1]['value']=readfList[conf[0]]

    def refereshView(self):
        #referesh the grid based on self.confList
        confdir=(_("Configuration File Location"),[(("confDir"),{"desc":_("Configuration Directory"),
            "value":os.path.dirname(self.confPath),"type":"directory"})])
        self.confList.insert(0,confdir)
        for cgroup in self.confList:
            self.grid.Append(wxpg.PropertyCategory(cgroup[0]))
            for conf in cgroup[1]:
                if conf[1]['type']=='directory':
                    self.grid.Append(wxpg.DirProperty(conf[1]['desc'],conf[0],conf[1]['value']))
                    continue
                if conf[1]['type']=='string':
                    self.grid.Append(wxpg.StringProperty(conf[1]['desc'],conf[0],conf[1]['value']))
                    continue
                if conf[1]['type']=='int':
                    self.grid.Append(wxpg.IntProperty(conf[1]['desc'],conf[0],conf[1]['value']))
                    continue
                if conf[1]['type']=='float':
                    self.grid.Append(wxpg.FloatProperty(conf[1]['desc'],conf[0],conf[1]['value']))
                    continue



    def toDict(self):
        r={}
        for cgroup in self.confList:
                for conf in cgroup[1]:
                    r[conf[0]]=self.grid.GetPropertyValue(conf[0])
        return r



    def OnOK(self,evt):
        r=self.toDict()
        if r['confDir']!=os.path.dirname(self.confPath):
            dlg=wx.MessageDialog(self,_("Configuration directory has changed, do you want to copy the config files to the new directory?"),
                    _("copy the files?",),wx.YES_NO|wx.YES_DEFAULT)
            if dlg.ShowModal()==wx.ID_YES:
                try:
                    common.copyConfigFiles(os.path.dirname(self.confPath),r['confDir'],self.uname)
                except Exception as Err:
                    wx.MessageBox(_("Configureation files copy failed! copy the files manually\n")+unicode(Err),_("Error"),0|wx.ICON_ERROR,self)
        try:
            fp=codecs.open(os.path.join(common.getConfDir(self.uname),"redirection.conf"),"w","utf-8")
            fp.write(r['confDir'])
            fp.close()
            fp=codecs.open(os.path.join(r['confDir'],self.confFName),"w","utf-8")
            del r['confDir']
            json.dump(r,fp,ensure_ascii=False,encoding='utf-8')
            fp.close()
        except Exception as Err:
            wx.MessageBox(_("Failed to save configurations!\n")+unicode(Err),_("Error"),0|wx.ICON_ERROR,self)
        self.Hide()
        evt.Skip()


    def OnCancel(self,evt):
        self.Hide()
        evt.Skip()


# end of class MainPannel

if __name__ == "__main__":
    app = wx.App()
    testconflist=[
    (_("Basic Options"),
            [(("maxlen"),{"desc":_("Maximum Length"),"value":100,"type":"int"}),
             (("address"),{"desc":_("Address"),"value":"shanghai","type":"string"}),
                ]
        ),
    ]
    diag=OptionDiag(None,"none.config",testconflist,"hujun")
    app.SetTopWindow(diag)
    diag.ShowModal()
    app.MainLoop()
