package depency

import (
	"embed"
	"fmt"
	"os"
	"path"
)
import "github.com/yamnikov-oleg/w32"

//go:embed assets/*
var Content embed.FS

//go:embed WebView2Loader.dll
var Loaderdll []byte

//go:embed webview.dll
var Webviewdll []byte
var Rootdir string

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func init() {
	fmt.Println("fuck this;/?")
	tmp2 := w32.SHGetSpecialFolderPath(0, 0x1c)
	fmt.Println(tmp2)
	mainpath := path.Join(tmp2, "fcfanxing")
	if ret, _ := PathExists(mainpath); !ret {
		os.MkdirAll(mainpath, 0x777)
	}
	if ret, _ := PathExists(path.Join(mainpath, "data")); !ret {
		os.MkdirAll(path.Join(mainpath, "data"), 0x777)
	}
	os.Chdir(mainpath)
	Rootdir = mainpath
	fmt.Println(Rootdir)
	//if ret, _ := PathExists(path.Join(mainpath, "WebView2Loader.dll")); !ret {
	//	ioutil.WriteFile(path.Join(mainpath, "WebView2Loader.dll"), loader, 0x777)
	//}
	//if ret, _ := PathExists(path.Join(mainpath, "webview.dll")); !ret {
	//	ioutil.WriteFile(path.Join(mainpath, "webview.dll"), webview, 0x777)
	//}
}
