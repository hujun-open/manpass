#-------------------------------------------------------------------------------
# Name:        module1
# Purpose:
#
# Author:      hujun
#
# Created:     10/08/2015
# Copyright:   (c) hujun 2015
# Licence:     <your licence>
#-------------------------------------------------------------------------------
import wx
import platform



_ = wx.GetTranslation

class SelfDestroyDiag(wx.Dialog):
    def __init__(self,parent,msg,timeout):

        wx.Dialog.__init__(self,parent,title=_("Password"))
        self.label_passwd=wx.StaticText(self, wx.ID_ANY,label=msg,style=wx.ALIGN_LEFT )
        self.timer_close=wx.Timer(self,wx.NewId())
        myos=platform.system()
        if myos=="Windows":
            fontf="Consolas"
        elif myos=="Linux":
            fontf="DejaVu Sans Mono"
        myfont=wx.Font(20,wx.MODERN,wx.FONTSTYLE_NORMAL,wx.FONTWEIGHT_NORMAL,
                            faceName=fontf)
        self.label_passwd.SetFont(myfont)
        self.timer_close.Start(timeout*1000,wx.TIMER_ONE_SHOT)
        self.__do_layout()
        self.Bind(wx.EVT_TIMER,self.OnTimer,self.timer_close)
        okb=self.bsizer.GetAffirmativeButton()
        self.Bind(wx.EVT_BUTTON,self.OnOK,okb)



    def __do_layout(self):
        # begin wxGlade: MainPannel.__do_layout
        sizer_all=wx.BoxSizer(wx.VERTICAL)
        sizer_all.Add(self.label_passwd,1,wx.FIXED_MINSIZE|wx.ALIGN_CENTER|wx.ALL,10)
        self.bsizer=self.CreateButtonSizer(wx.OK)
        sizer_all.Add(self.bsizer,0, wx.ALL|wx.ALIGN_CENTER_HORIZONTAL, 5)
        self.SetSizer(sizer_all)
        sizer_all.Fit(self)
        self.Layout()


    def OnTimer(self,evt):
        self.Close()

    def OnOK(self,evt):
        self.timer_close.Stop()
        evt.Skip()



if __name__ == '__main__':
    app=wx.App()
    dlg=SelfDestroyDiag(None,"0O0O`'wW",5)
    app.SetTopWindow(dlg)
    dlg.ShowModal()
    dlg.Destroy()
    app.MainLoop()
