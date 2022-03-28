package main

import (
	"fanxing/depency"
	"fanxing/rsa"
	"fmt"
	"github.com/jchv/go-webview2"
	"github.com/jchv/go-webview2/pkg/edge"
	"github.com/kirinlabs/HttpRequest"
	"github.com/yamnikov-oleg/w32"
	"github.com/yamnikov-oleg/wingo"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
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

var mainview webview2.WebView
var gameview webview2.WebView
var mainwindow *wingo.Window

func onresize(w *wingo.Window, xy wingo.Vector) {
	//if nil != mainview {
	//
	//	mainview.Resize()
	//}
}
func realopenurl(url string) {
	runtime.LockOSThread()
	//w := mainview.New(true)
	//defer w.Destroy()
	//w.Navigate(url)
	//w.Run()

}

func openurl(url string) int {
	fmt.Println(url)
	go realopenurl(url)
	return 1
}

//func setcookie(cookie string) {
//	fmt.Println(cookie)
//	mainview.Dispatch(func() { mainview.Eval("addacountcookie('" + cookie + "');") }) //addacountcookie
//}
func closelogin() {
	fmt.Println("close login window!!!")
}
func realshowgameview(dir string) {
	runtime.LockOSThread()
	fmt.Println("datadir:", path.Join(depency.Rootdir, "tmplg", dir))

	tmpv := webview2.NewWithOptions(webview2.WebViewOptions{
		DataPath:  path.Join(depency.Rootdir, "tmplg", dir),
		Debug:     true,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title: "Minimal webview example",
		},
	})

	tmpv.SetSize(410, 490, 1)

	url := "https://fanxing.kugou.com/cterm/edge/game_tower_defense/m/views/index.html?roomId=3174234&amp;_=16473101&amp;_qframe_scope=1647310108234_2"
	tmpv.Navigate(url)
	tmpv.Run()

}
func showgameview(dir string) {
	go realshowgameview(dir)
}
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
func deleteunuseddir(dirs string) {
	fmt.Println("dirs:::", dirs)
	ds := strings.Split(dirs, "_")
	fmt.Println("dirs:::", ds)

	dir, _ := ioutil.ReadDir(path.Join(depency.Rootdir, "tmplg"))

	for _, d := range dir {
		if !Contains(ds, d.Name()) {
			os.RemoveAll(path.Join(depency.Rootdir, "tmplg", d.Name()))
		}
	}
}
func realopenloginurl(url string) {
	runtime.LockOSThread()
	rand.Seed(time.Now().UnixNano())
	var subdir string
	for {
		subdir = strconv.Itoa(rand.Intn(1000000))
		if ok, _ := depency.PathExists(path.Join(depency.Rootdir, "tmplg", subdir)); !ok {
			break
		}
	}
	os.MkdirAll(path.Join(depency.Rootdir, "tmplg", subdir), 0x777)
	tmpview := webview2.NewWithOptions(webview2.WebViewOptions{
		DataPath:  path.Join(depency.Rootdir, "tmplg", subdir),
		Debug:     true,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title: "Minimal webview example",
		},
	})
	tmpview.SetSize(648, 293, 1)
	tmpview.Bind("setcookie", func(cookie string) {

		mainview.Dispatch(func() { mainview.Eval("addacountcookie('" + cookie + "','" + subdir + "');") }) //addacountcookie
		tmpview.Destroy()
		tmpview.Terminate()

		tmpexp, _ := regexp.Compile("KugooID=(\\d+)")
		tmpmatch := tmpexp.FindStringSubmatch(cookie)
		fmt.Println("getkugouid", tmpmatch)
		if len(tmpmatch) >= 2 {
			//os.Rename(path.Join(depency.Rootdir, "tmplg", subdir), path.Join(depency.Rootdir, "tmplg", tmpmatch[1]+"_data"))
		}
	})
	//tmpview.Bind("closelogin", func() { tmpview.Terminate();
	//
	//})
	tmpview.Navigate(url)
	tmpview.Init("window.onload=function(){ let ck=document.cookie; if(!Kg.Cookie.read('_fxNickName')){showLogin()}else{window.setcookie(ck)} }")
	tmpview.Run()
	mainview.Dispatch(func() { mainview.Eval("window.jsdeleteunseddir()") })
	//w := mainview.New(true)
	//defer tmpview.Destroy()
	//w.Navigate(url)
	//w.Run()

}

func openloginurl(url string) int {
	fmt.Println(url)
	go realopenloginurl(url)
	return 1
}
func onhittest(w *wingo.Window, xy wingo.Vector) (uint32, bool) {

	return w32.HTCLIENT, true
}

type TagnccalcsizeParams struct {
	rgrc  [3]w32.RECT
	lppos *w32.COORD
}

func gorequest(url string, cookie string) interface{} {
	req := HttpRequest.NewRequest().Debug(true)
	req.SetHeaders(map[string]string{
		"Cookie":     cookie,
		"Connection": "keep-alive",
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36",
	})
	fmt.Println("url", url)
	fmt.Println("cookie", cookie)
	resp, err := req.Get(url) //
	if err != nil {
		fmt.Println(err)
		return false
	}
	json, ok := resp.Export()
	if ok != nil {
		//fmt.Println("error happend")
		return false
	} else {
		fmt.Println("requst:", json)
		return json
	}
}

func tokenrequest(url, token string) string {
	req := HttpRequest.NewRequest()
	fmt.Println("token", token)
	req.SetHeaders(map[string]string{
		"Authorization": "Bearer " + token,
		"Connection":    "keep-alive",
		"user-agent":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36",
	})

	resp, err := req.Get(url) //
	if err != nil {
		return `{"code": "-1", "data": ""}`
	}
	body, ok := resp.Body()

	if ok != nil {

		return `{"code": "-1", "data": ""}`
	} else {
		if Rsascript {
			tmp := rsa.DecodeByte(body)
			return string(tmp)
		} else {
			return string(body)
		}

	}
}
func gouserinfo(token string) string {
	var tmpinfo string
	if Rsascript {

		tmpinfo = tokenrequest(Domain+"api/user/getuserinfo", token)
	} else {
		tmpinfo = tokenrequest("http://127.0.0.1:8000/api/user/getuserinfo", token)
	}
	fmt.Println("myuserinfo::", tmpinfo)
	return tmpinfo
}
func gopost(url string, data string, cookie string) interface{} {
	fmt.Println("url", url)
	fmt.Println("dat:", data)
	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"Cookie":     cookie,
		"Connection": "keep-alive",
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36",
	})

	resp, err := req.Post(url, data) //
	if err != nil {
		return false
	}
	json, ok := resp.Export()
	if ok != nil {
		fmt.Println("error happend")
		return false
	} else {
		fmt.Println(json)
		return json
	}
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
	w32.MoveWindow(hwnd, x, y, w, h, false)
}

var Rsascript = false
var Domain = "http://127.0.0.1:8000/"

func getscript(url string) ([]byte, bool) {
	_filename := path.Base(url)
	filename := _filename[0 : len(_filename)-3]
	if ret, _ := depency.PathExists(path.Join(depency.Rootdir, "data", filename)); ret {
		tmp, _ := ioutil.ReadFile(path.Join(depency.Rootdir, "data", filename))
		if Rsascript {
			return rsa.DecodeByte(tmp), true
		} else {
			return tmp, true
		}
	}

	req := HttpRequest.NewRequest()
	resp, err := req.Get(url) //
	if err != nil {
		return []byte(""), false
	}
	result, ok := resp.Body()
	if ok != nil {
		//fmt.Println("error happend")
		return []byte(""), false
	} else {
		//fmt.Println(json)
		ioutil.WriteFile(path.Join(depency.Rootdir, "data", filename), result, 0x777)
		if Rsascript {
			return rsa.DecodeByte(result), true
		} else {
			return result, true
		}
		//return result, true
	}
}
func oncreate(w *wingo.Window, url string) {
	runtime.LockOSThread()

	mainview = webview2.NewWithOptions(webview2.WebViewOptions{
		DataPath:  path.Join(depency.Rootdir, "mainv"),
		Debug:     true,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:        "Minimal webview example",
			Parentwindow: uintptr(w.GetHandle()),
		},
	})
	if Rsascript {
		mainview.AddWebResourceRequestedFilter("*.js", 6)
		mainview.SetWebResourceRequestedCallback(func(request *edge.ICoreWebView2WebResourceRequest, args *edge.ICoreWebView2WebResourceRequestedEventArgs) {

			//tmp, _ := args.GetResponse()
			tmprequest, _ := args.GetRequest()
			url, _ := tmprequest.GetUri()
			content, _ := getscript(url)
			env := mainview.WebviewEnvironment()
			//content := []byte("console.log(123123)")
			tmps := "content-length: $strlen\r\ncontent-type: text/javascript"
			str := strings.Replace(tmps, "$strlen", strconv.Itoa(len(content)), 1)
			response, _ := env.CreateWebResourceResponse(content, 200, "OK", str)
			args.PutResponse(response)
			//tmp.PutStatusCode()

			//fmt.Println(args)

		})
	}
	mainview.Navigate(Domain)

	mainview.Bind("open_url", openurl)
	mainview.Bind("open_loginurl", openloginurl)
	mainview.Bind("gorequest", gorequest)
	mainview.Bind("gopost", gopost)
	mainview.Bind("deleteunuseddir", deleteunuseddir)

	mainview.Bind("gouserinfo", gouserinfo)
	//mainview.Bind("capthurewindow", capthurewindow)
	//mainview.Bind("releasewindow", releasewindow)
	mainview.Bind("movewindow", movewindow)
	mainview.Bind("switchgameuser", showgameview)
	mainview.Bind("closewindow", func() {
		mainview.Destroy()
		mainview.Terminate()
		w.Destroy()
		wingo.Exit()
	})
	mainview.Bind("minuswindow", func() {
		w.Hide()
	})
	mainview.Run()
	//mainview.
	//mainview.Run()

	//defer mainview.Destroy()
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
	//depency.Webviewdll
	os.Setenv("WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", "--disable-web-security --allow-insecure-localhost")
	urlchan := make(chan string)
	go myhttp(urlchan)
	prefix := <-urlchan
	fmt.Println(prefix)
	mainwindow = wingo.NewWindow(true, true)
	size := wingo.Vector{600, 800}
	mainwindow.SetSize(size)
	mainwindow.OnClose = onclose
	mainwindow.OnSizeChanged = onresize
	mainwindow.OnTrayClick = ontrayclick
	//mainwindow.OnHITTEST = onhittest
	//mainwindow.OnLBTNDOWN = onlbtndown
	//mainwindow.OnNCCALCSIZE = onnccalcsize
	mainwindow.OnTrayRightClick = ontrayrightclick
	//go showgameview()
	dw := w32.GetWindowLong(mainwindow.GetHandle(), w32.GWL_STYLE)
	dw = dw & ^w32.WS_CAPTION    //取消标题栏
	dw = dw & ^w32.WS_THICKFRAME //取消拖动改变大小//不取消的话，自绘标题栏上面会有一条白边而且覆盖不了
	w32.SetWindowLong(mainwindow.GetHandle(), w32.GWL_STYLE, uint32(dw))

	//取消边框内的边缘，也就是取消3D效果
	dw = w32.GetWindowLong(mainwindow.GetHandle(), w32.GWL_EXSTYLE)
	dw = dw & ^w32.WS_EX_DLGMODALFRAME
	dw = dw & ^w32.WS_EX_CLIENTEDGE
	dw = dw & ^w32.WS_EX_WINDOWEDGE
	w32.SetWindowLong(mainwindow.GetHandle(), w32.GWL_EXSTYLE, uint32(dw))

	//w.OnCreate=oncreate
	icon := wingo.LoadIcon(101)
	mainwindow.SetIcon(icon)
	mainwindow.AddTrayIcon(icon, "繁星屠龙助手")
	traymenu = wingo.NewContextMenu()
	exitbtn := traymenu.AppendItemText("退出")
	exitbtn.OnClick = func(item *wingo.MenuItem) {
		mainview.Destroy()
		mainview.Terminate()
		mainwindow.Destroy()

		wingo.Exit()
	}

	oncreate(mainwindow, prefix)

	mainwindow.Show()

	wingo.Start()
}
