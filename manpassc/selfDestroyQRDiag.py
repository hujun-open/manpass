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
from PIL import Image
import qrcode


_ = wx.GetTranslation

def WxImageToWxBitmap( myWxImage ) :
    return myWxImage.ConvertToBitmap()

def WxBitmapToPilImage( myBitmap ) :
    return WxImageToPilImage( WxBitmapToWxImage( myBitmap ) )

def WxBitmapToWxImage( myBitmap ) :
    return wx.ImageFromBitmap( myBitmap )


def PilImageToWxBitmap( myPilImage ) :
    return WxImageToWxBitmap( PilImageToWxImage( myPilImage ) )

def PilImageToWxImage( myPilImage ):
    myWxImage = wx.EmptyImage( myPilImage.size[0], myPilImage.size[1] )
    myWxImage.SetData( myPilImage.convert( 'RGB' ).tostring() )
    return myWxImage


class SelfDestroyQRDiag(wx.Dialog):
    def __init__(self,parent,upass,timeout):

        wx.Dialog.__init__(self,parent,title=_("Password QR Code"))
        myqr=qrcode.QRCode()
        myqr.add_data(upass)
        myqr.make(fit=True)
        img = myqr.make_image()
        passbitmap=PilImageToWxBitmap(img)

        self.imgctrl=wx.StaticBitmap(self,-1,passbitmap)
        self.timer_close=wx.Timer(self,wx.NewId())
        self.timer_close.Start(timeout*1000,wx.TIMER_ONE_SHOT)
        self.__do_layout()
        self.Bind(wx.EVT_TIMER,self.OnTimer,self.timer_close)
        okb=self.bsizer.GetAffirmativeButton()
        self.Bind(wx.EVT_BUTTON,self.OnOK,okb)



    def __do_layout(self):
        # begin wxGlade: MainPannel.__do_layout
        sizer_all=wx.BoxSizer(wx.VERTICAL)
        sizer_all.Add(self.imgctrl,1,wx.FIXED_MINSIZE|wx.ALIGN_CENTER|wx.ALL,10)
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
    dlg=SelfDestroyQRDiag(None,"mypassword",5)
    app.SetTopWindow(dlg)
    dlg.ShowModal()
    dlg.Destroy()
    app.MainLoop()
