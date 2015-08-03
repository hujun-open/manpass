#!/usr/bin/env python
# -*- coding: UTF-8 -*-
#
# generated by wxGlade 0.7.0 on Mon Jul 13 05:41:37 2015
#

import wx
import common
import os.path

# begin wxGlade: dependencies
import gettext
# end wxGlade
import genPassDiag
import shlex
import subprocess
import time
import threading
import Queue
import sys
import newUserDiag
import common
import traceback


# begin wxGlade: extracode
# end wxGlade
_ = wx.GetTranslation




class ChangeMasterPassDiag(wx.Dialog):
    def __init__(self, parent):
        # begin wxGlade: MainPannel.__init__
        style =  wx.CAPTION|wx.SYSTEM_MENU
        wx.Dialog.__init__(self,parent,style=style)
        self.text_ctrl_currentpass = wx.TextCtrl(self, wx.ID_ANY, "",size=(200,-1),style=wx.TE_PASSWORD)
        self.label_currentpass = wx.StaticText(self, wx.ID_ANY,label=_("Current Password:"),style=wx.ALIGN_LEFT)
        self.text_ctrl_upass1 = wx.TextCtrl(self, wx.ID_ANY, "",size=(200,-1),style=wx.TE_PASSWORD)
        self.label_upass1 = wx.StaticText(self, wx.ID_ANY,label=_("New Password:"),style=wx.ALIGN_RIGHT)
        self.text_ctrl_upass2 = wx.TextCtrl(self, wx.ID_ANY, "",size=(200,-1),style=wx.TE_PASSWORD)
        self.label_upass2 = wx.StaticText(self, wx.ID_ANY,label=_("Type Again:"),style=wx.ALIGN_RIGHT)
        self.pbar=wx.Gauge(self)
        self.label_pbar=wx.StaticText(self,label="label-1",style=wx.ALIGN_CENTER)

        self.__set_properties()
        self.__do_layout()


        okb= self.bsizer.GetItem(1).GetWindow()
        canb=self.bsizer.GetItem(2).GetWindow()
        self.Bind(wx.EVT_BUTTON,self.OnOK,okb)
        self.Bind(wx.EVT_BUTTON,self.OnCancel,canb)
        self.Bind(common.EVT_MANPASS_ERR,self.OnErr)
        self.Bind(common.EVT_MANPASS_PROGRESS,self.OnProgress)
        self.Bind(common.EVT_MANPASS_PROGRESS_LABEL,self.UpdateProgressLabel)
        self.Bind(common.EVT_MANPASS_LOAD_DONE,self.Done)

        self.Centre()
        self.text_ctrl_currentpass.SetFocus()
        # end wxGlade

    def __set_properties(self):
        # begin wxGlade: MainPannel.__set_properties
        self.SetTitle(_("Change master password"))
        #self.SetWindowStyle(wx.BORDER_DEFAULT)
        # end wxGlade

    def __do_layout(self):
        # begin wxGlade: MainPannel.__do_layout
        sizer_all=wx.BoxSizer(wx.VERTICAL)
        sizer_v = wx.FlexGridSizer(3,2,5,5)

        sizer_v.Add(self.label_currentpass,0,wx.FIXED_MINSIZE|wx.ALIGN_RIGHT|wx.TOP,5)
        sizer_v.Add(self.text_ctrl_currentpass,3,wx.TOP|wx.EXPAND,5)
        sizer_v.Add(self.label_upass1,0,wx.FIXED_MINSIZE|wx.ALIGN_RIGHT|wx.TOP,5)
        sizer_v.Add(self.text_ctrl_upass1,3,wx.TOP|wx.EXPAND,5)
        sizer_v.Add(self.label_upass2,0,wx.FIXED_MINSIZE|wx.ALIGN_RIGHT|wx.TOP,5)
        sizer_v.Add(self.text_ctrl_upass2,3,wx.TOP|wx.EXPAND,5)
        sizer_all.Add(sizer_v,0, wx.ALL|wx.ALIGN_CENTER_HORIZONTAL, 5)

        sizer_h2=wx.BoxSizer(wx.HORIZONTAL)
        sizer_h2.Add(self.label_pbar,1,wx.ALIGN_CENTER_HORIZONTAL|wx.ALIGN_CENTER_VERTICAL|wx.TOP,0)
        sizer_all.Add(sizer_h2,1,wx.ALIGN_CENTER_HORIZONTAL|wx.ALIGN_CENTER_VERTICAL|wx.TOP,0)

        sizer_h3=wx.BoxSizer(wx.HORIZONTAL)
        sizer_h3.Add(self.pbar,1,wx.ALIGN_CENTER_HORIZONTAL|wx.ALIGN_CENTER_VERTICAL|wx.TOP,0)
        sizer_all.Add(sizer_h3,1,wx.ALIGN_CENTER_HORIZONTAL|wx.ALIGN_CENTER_VERTICAL|wx.TOP|wx.BOTTOM,0)

        self.pbar.Hide()
        self.label_pbar.Hide()

        self.bsizer=self.CreateButtonSizer(wx.OK|wx.CANCEL)
        sizer_all.Add(self.bsizer,0, wx.ALL|wx.ALIGN_CENTER_HORIZONTAL, 5)
        self.SetSizer(sizer_all)
        sizer_all.Fit(self)
        self.Layout()
        # end wxGlade


    def disableMe(self):

        self.text_ctrl_currentpass.Disable()
        self.text_ctrl_upass1.Disable()
        self.text_ctrl_upass2.Disable()
        okb= self.bsizer.GetItem(1).GetWindow()
        canb=self.bsizer.GetItem(2).GetWindow()
        self.orig_label=okb.GetLabel()
        okb.Disable()
        canb.Disable()

    def enableMe(self):
        self.text_ctrl_currentpass.Enable()
        self.text_ctrl_upass1.Enable()
        self.text_ctrl_upass2.Enable()
        okb= self.bsizer.GetItem(1).GetWindow()
        canb=self.bsizer.GetItem(2).GetWindow()
        okb.SetLabel(self.orig_label)
        okb.Enable()
        canb.Enable()


    def OnCancel(self,evt):
        evt.Skip()

    def OnOK(self,evt):
        b=evt.GetEventObject()
        self.disableMe()
        b.SetLabel(_("Updating"))
        self.Update()
        self.currentpass=self.text_ctrl_currentpass.GetValue()

        if self.currentpass != self.GetParent().GetParent().upass:
            wx.MessageBox(_("Wrong Password!"),_("Error"),0|wx.ICON_ERROR,self)
            self.enableMe()
            return

        self.pass1=self.text_ctrl_upass1.GetValue()
        pass2=self.text_ctrl_upass2.GetValue()
        if self.pass1!=pass2:
            wx.MessageBox(_("Password of two typing doesn't match!"),_("Error"),0|wx.ICON_ERROR,self)
            self.enableMe()
            return
        try:
            newUserDiag.goodPass(self.pass1)
        except Exception as Err:
            wx.MessageBox(_("Not a good password, choose a different password\n")+unicode(Err),_("Error"),0|wx.ICON_ERROR,self)
            self.enableMe()
            return

        def replacethem(self):
            try:
                pevt=common.ManpassProgressLabel(Label=_("Getting the latest list"))
                wx.PostEvent(self,pevt)
                listctrl=self.GetParent()
                mlist=listctrl.apc.getAllMetaId()
                rlist=[]
                pevt=common.ManpassProgressLabel(Label=_("Getting history passwords"))
                wx.PostEvent(self,pevt)
                i=0
                for m in mlist:
                    i+=1
                    pevt=common.ManpassProgressEVT(Range=len(mlist),Pos=i)
                    wx.PostEvent(self,pevt)
                    rlist+=listctrl.apc.getAllRecodsForMeta(mid=m,win=None)
                pevt=common.ManpassProgressLabel(Label=_("Re-Encrypting files"))
                wx.PostEvent(self,pevt)
                common.reEncryptCertFiles(listctrl.GetParent().confDict['confDir'],self.currentpass.encode("utf-8"),self.pass1.encode('utf-8'))
                pevt=common.ManpassProgressLabel(Label=_("Re-Encrypting passwords"))
                wx.PostEvent(self,pevt)
                listctrl.apc.replaceAll(rlist,self.pass1.encode('utf-8'),self)
                listctrl.GetParent().upass=self.pass1.encode('utf-8')
                listctrl.apc.masterpass=self.pass1.encode('utf-8')
                listctrl.GetParent().apiclient.masterpass=self.pass1.encode('utf-8')
                devt=common.ManpassLoadingDone(Type="final")
                wx.PostEvent(self,devt)


            except Exception as Err:
                errevt=common.ManpassErrEVT(Value=unicode(Err))
                wx.PostEvent(self,errevt)
                traceback.print_exc(Err)
        self.showProgress()
        t=threading.Thread(target=replacethem,args=(self,))
        t.daemon=True
        t.start()

    def OnErr(self,evt):
        wx.MessageBox(_("Updating records failed!\n")+evt.Value,_("Error"),0|wx.ICON_ERROR,self)
        self.enableMe()

    def showProgress(self):
        self.label_currentpass.Hide()
        self.label_upass1.Hide()
        self.label_upass2.Hide()
        self.text_ctrl_currentpass.Hide()
        self.text_ctrl_upass1.Hide()
        self.text_ctrl_upass2.Hide()
        self.label_pbar.Show(True)
        self.pbar.Show(True)
        self.Layout()

    def OnProgress(self,evt):
        self.pbar.SetRange(evt.Range)
        self.pbar.SetValue(evt.Pos)


    def UpdateProgressLabel(self,evt):
        self.label_pbar.SetLabel(evt.Label)
        self.Layout()


    def Done(self,evt):
        if hasattr(evt,"Type"):
            if evt.Type=='final':
                wx.MessageBox(_("Master password has been changed successfully!"),_("Done"))
                self.Destroy()
        evt.Skip()

# end of class MainPannel

if __name__ == "__main__":
    app = wx.App()
    diag=ChangeMasterPassDiag(None)
    app.SetTopWindow(diag)
    diag.showProgress()
    diag.ShowModal()
    app.MainLoop()
