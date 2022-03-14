package main

import (
	"fanxing/depency"
	"fmt"
	"github.com/webview/webview"
	"github.com/yamnikov-oleg/w32"
	"github.com/yamnikov-oleg/wingo"
	"net"
	"net/http"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

func onclose(w *wingo.Window) bool {
	fmt.Println("whhy")
	w.Hide()
	return false
}
func ontrayclick(w *wingo.Window) {
	w.Show()
}
func ontrayrightclick(w *wingo.Window) {
	traymenu.StartContext(w)

}

var mainview webview.WebView
var mainwindow *wingo.Window

func onresize(w *wingo.Window, xy wingo.Vector) {
	if nil != mainview {
		mainview.SetSize(xy.X, xy.Y, webview.HintNone)

	}
}
func realopenurl(url string) {
	runtime.LockOSThread()
	w := webview.New(true)
	defer w.Destroy()
	w.Navigate(url)
	w.Run()

}
func openurl(url string) int {
	fmt.Println(url)
	go realopenurl(url)
	return 1
}
func onhittest(w *wingo.Window, xy wingo.Vector) (uint32, bool) {

	return w32.HTCLIENT, true
}

type TagnccalcsizeParams struct {
	rgrc  [3]w32.RECT
	lppos *w32.COORD
}

func onnccalcsize(w *wingo.Window, b int, lpncsp uintptr) {
	p := (*TagnccalcsizeParams)(unsafe.Pointer(lpncsp))
	var nTitleHeight int32 = 1
	var nFrameBorerL int32 = 1
	var nFrameBorerR int32 = 1
	var nFrameBorerB int32 = 1
	p.rgrc[0].Top += nTitleHeight
	p.rgrc[0].Left += nFrameBorerL
	p.rgrc[0].Right -= nFrameBorerR
	p.rgrc[0].Bottom -= nFrameBorerB

}
func onlbtndown(w *wingo.Window, xy wingo.Vector) {
	fmt.Println(xy.X)
	fmt.Println(xy.Y)
}
func capthurewindow() {

}
func releasewindow() {

}
func movewindow(x int, y int, w int, h int) {
	hwnd := mainwindow.GetHandle()
	//xint, _ := strconv.Atoi(x)
	//yint, _ := strconv.Atoi(y)
	//wint, _ := strconv.Atoi(width)
	//hint, _ := strconv.Atoi(height)
	fmt.Println(x, y, w, h)
	w32.MoveWindow(hwnd, x, y, w, h, false)
}
func oncreate(w *wingo.Window, url string) {
	runtime.LockOSThread()
	hand := w.GetHandle()
	mainview = webview.NewWindow(true, unsafe.Pointer(&hand))

	//mainview=&wv
	//mainview.Navigate(url + "/assets/index.html")
	mainview.Navigate("http://127.0.0.1:8080/")
	mainview.Bind("open_url", openurl)
	mainview.Bind("capthurewindow", capthurewindow)
	mainview.Bind("releasewindow", releasewindow)
	mainview.Bind("movewindow", movewindow)
	mainview.Run()

	defer mainview.Destroy()
}

var traymenu *wingo.Menu

func myhttp(urlchan chan string) {
	http.Handle("/", http.FileServer(http.FS(depency.Content)))
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	portAddress := listener.Addr().String()
	urlchan <- "http://" + portAddress
	listener.Close()
	http.ListenAndServe(portAddress, nil)
}
func wndProc(hwnd w32.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	oldproc := w32.GetProp(hwnd, "proc")
	var nargs uintptr = 4
	fmt.Println("workd?????")
	if msg == w32.WM_NCHITTEST {
		return w32.TRANSPARENT
	}
	ret, _, _ := syscall.Syscall9(oldproc, nargs, uintptr(hwnd), uintptr(msg), wParam, lParam,
		0,

		0,

		0,

		0,

		0)
	return ret

}
func checkhwnd(w *wingo.Window) {
	var pwin2 w32.HWND = 0
	for {
		time.Sleep(3000000000)
		pwin2 = w32.FindWindowEx(w.GetHandle(), "Chrome_WidgetWin_0")
		if pwin2 != 0 {
			fmt.Println("hand;e::", pwin2)

			break
		} else {
			fmt.Println("what a fuck???")
		}

	}
	if pwin2 != 0 {
		w32.SetProp(pwin2, "proc", w32.GetWindowLongPtr(pwin2, w32.GWLP_WNDPROC))
		fmt.Println("sethook")
		w32.SetWindowLongPtr(pwin2, w32.GWLP_WNDPROC, syscall.NewCallback(wndProc))
	}
}
func main() {
	urlchan := make(chan string)
	go myhttp(urlchan)
	prefix := <-urlchan
	fmt.Println(prefix)
	mainwindow = wingo.NewWindow(true, true)
	size := wingo.Vector{500, 500}
	mainwindow.SetSize(size)
	mainwindow.OnClose = onclose
	mainwindow.OnSizeChanged = onresize
	mainwindow.OnTrayClick = ontrayclick
	//mainwindow.OnHITTEST = onhittest
	//mainwindow.OnLBTNDOWN = onlbtndown
	//mainwindow.OnNCCALCSIZE = onnccalcsize
	mainwindow.OnTrayRightClick = ontrayrightclick

	dw := w32.GetWindowLong(mainwindow.GetHandle(), w32.GWL_STYLE)
	dw = dw & ^w32.WS_CAPTION    //取消标题栏
	dw = dw & ^w32.WS_THICKFRAME //取消拖动改变大小//不取消的话，自绘标题栏上面会有一条白边而且覆盖不了
	w32.SetWindowLong(mainwindow.GetHandle(), w32.GWL_STYLE, uint32(dw))

	//取消边框内的边缘，也就是取消3D效果
	//dw := w32.GetWindowLong(w.GetHandle(), w32.GWL_EXSTYLE)
	//dw = dw & ^w32.WS_EX_DLGMODALFRAME
	//dw = dw & ^w32.WS_EX_CLIENTEDGE
	//dw = dw & ^w32.WS_EX_WINDOWEDGE
	//w32.SetWindowLong(w.GetHandle(), w32.GWL_EXSTYLE, uint32(dw))

	//w.OnCreate=oncreate
	icon := wingo.LoadIcon(101)
	mainwindow.SetIcon(icon)
	mainwindow.AddTrayIcon(icon, "hellofromaiel")
	traymenu = wingo.NewContextMenu()
	exitbtn := traymenu.AppendItemText("退出")
	exitbtn.OnClick = func(item *wingo.MenuItem) {
		mainwindow.Destroy()
	}

	oncreate(mainwindow, prefix)

	mainwindow.Show()

	wingo.Start()
}
