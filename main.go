package main

/*
export GOROOT=/usr/local/go
export GOPATH=/Users/zzx/Desktop/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN:$GOROOT/bin

go build -buildmode=c-archive  -o DBmanage.a main.go

windows
go build -ldflags="-H windowsgui" -o DBmanage.exe

//(kill -9 $(pidof invoice))&&(nohup ./invoice &) &&(ps -aux | grep "invoice")

windows build  cannot find -lwebview 先去github.com/webview/webview"  script run bat
 $env:GOPROXY = "https://proxy.golang.com.cn,direct"

*/

import (
	"DataBaseManage/HTTPBusiness"
	"DataBaseManage/asset"
	"DataBaseManage/dal"
	"DataBaseManage/public"
	"fmt"

	"github.com/webview/webview"

	//"github.com/zserge/webview"

	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	//"DataBaseManage/dal"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
)

import "C"

var IP string
var PthSep string
var sandboxDir string

func main() {
	fmt.Printf("curdir= ")
	osname := public.GetOSName()
	public.Log("osname=" + osname)
	//if osname != "darwin" {
	curdir := public.GetCurDir()
	fmt.Printf("curdir= " + curdir)
	StartServer(C.CString(curdir), 0)
	//}
}

//export StartServer
func StartServer(dir *C.char, autoo int) {
	//curdir := string(dir)
	curdir := C.GoString(dir)
	public.Log("StartServer")
	if autoo == 1 {
		array := strings.Split(curdir, "/Data/")
		curdir = array[0] + "/Data/Documents"
		public.CurPath = curdir
		dal.ISmac = 1

	}

	fmt.Printf("curdir= \n")
	fmt.Printf(public.CurPath)
	fmt.Printf("\n")
	public.Log("copyLocalDB")
	copyLocalDB()

	if autoo != 1 {
		checkDB()
	}
	IP = public.GetLocalIP()

	if autoo == 1 {
		go WebServerBase()
	} else {
		http_port := public.GetHttpPort()
		url := "http://" + IP + ":" + http_port + "/web/"
		//go open(url)
		go WebServerBase()

		debug := true
		w := webview.New(debug)
		defer w.Destroy()
		w.SetTitle("DBManage wechat:elink9988   ruobinzhu@gmail.com  ")
		w.SetSize(1200, 800, webview.HintNone)
		w.Navigate(url)
		w.Run()

	}

}
func copyLocalDB() {

	curdir := public.GetCurDir()
	gopath := public.GetGoPath()
	public.Log("curdir=" + curdir)
	PthSep = string(os.PathSeparator)
	public.CreatePath(curdir + PthSep + "db")
	dbpath := curdir + PthSep + "db" + PthSep + "db.sqlite"
	if !public.ExistsPath(dbpath) {
		if PthSep == "\\" {
			public.Log("copy local db ")
			public.CopyFile(gopath+"\\src\\DataBaseManage\\db\\db.sqlite", dbpath)

		} else {

			public.Log("copy local db ")
			public.CopyFile(gopath+"/src/DataBaseManage/db/db.sqlite", dbpath)
		}
	}
}

func checkDB() {
	curdir := public.GetCurDir()
	PthSep = string(os.PathSeparator)
	dbpath := curdir + PthSep + "db" + PthSep + dal.SqliteFileName
	if !public.ExistsPath(dbpath) {
		downloadDB()
	}
}
func downloadDB() {
	if public.NetWorkAvailable() {
		curdir := public.GetCurDir()
		PthSep = string(os.PathSeparator)

		if !public.ExistsPath(curdir + PthSep + "db") {
			public.Log(curdir + "/db not exists create")
			public.CreatePath(curdir + PthSep + "db")
			//public.Log("download1")
			zippath := curdir + PthSep + "db.zip"
			err := public.HttpDownloadFile("http://www.iosbuy.com/db/db.zip", zippath)
			//public.Log("download2")
			//public.Unzip("D:\\temp\\dd\\db.zip")
			if err == nil {
				public.Log("download success")
				if public.ExistsPath(zippath) {
					ret := public.Unzip(zippath)

					public.Log(ret)

					public.Log(zippath + " exists")

				} else {
					public.Log(zippath + " not exists")
				}

			} else {
				public.Log("download fail")
			}

		} else {
			public.Log("db exists")
		}
	} else {
		public.Log("can not download")
	}

}

func returnString() *C.char {
	gostring := "hello world"
	return C.CString(gostring)
}

//export GetIpstr
func GetIpstr() *C.char {
	if IP != "" {
		public.Log("return ip= " + IP)

		return C.CString(IP)
	}
	return nil
}

func openUrl(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	public.Log(url)
	go open(url)
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {

	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

//192.168.1.218:"+http_port+"/api/video.ashx?path=video/macinvoicehelp.mov
func WebServerBase() {
	//http := ModifierMiddleware
	serv := mux.NewRouter()
	http_port := public.GetHttpPort()

	serv.HandleFunc("/api/action", HTTPBusiness.ActionTask)
	serv.HandleFunc("/api/kv", HTTPBusiness.AddKVTask)
	//验证码
	serv.HandleFunc("/api/getvericode", HTTPBusiness.GetVeriCode)
	serv.HandleFunc("/api/vericode", HTTPBusiness.VeriCode)
	serv.HandleFunc("/api/openurl", openUrl)
	serv.PathPrefix("/api/showcodeimg/").Handler(http.StripPrefix("/api/showcodeimg/", HTTPBusiness.ShowVeriCode))
	serv.HandleFunc("/api/download", downloadDBTask)
	fs := assetfs.AssetFS{
		Asset:     asset.Asset,
		AssetDir:  asset.AssetDir,
		AssetInfo: asset.AssetInfo,
	}

	curdir := public.GetCurDir()
	PthSep := string(os.PathSeparator)

	dbPath := curdir + PthSep + "db"
	//osname := public.GetOSName() //os name=darwin

	serv.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(&fs)))
	serv.PathPrefix("/db/").Handler(http.StripPrefix("/db/", http.FileServer(http.Dir(dbPath))))
	public.Log("http_port=" + http_port)
	errser := http.ListenAndServe(":"+http_port, Middleware(serv))
	if errser != nil {
		public.Log(errser)
	}

}
func downloadDBTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}
	public.Log("downloadDBTask running...")
	//userid:= r.FormValue("userid")
	curdir := public.GetCurDir()
	filename := r.FormValue("filename")
	filetype := r.FormValue("filetype")
	if filetype == "sqlite" {
		filename = dal.SqliteFileName
	}

	data, err := ioutil.ReadFile(string(curdir + PthSep + "db" + PthSep + filename))

	if err == nil {
		public.Log("download file " + filename)
		w.Header().Add("Content Type", "application/x-kexiproject-sqlite3")
		w.Write(data)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 fot found " + dal.SqliteFileName))
	}
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//public.Log("middleware " + r.Method)
		//public.Log(r.URL)
		//public.Log("middleware ----------")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin,Authorization,Origin, X-Requested-With, Content-Type, Accept,common")

		h.ServeHTTP(w, r)

		if r.Method == "OPTIONS" {
			return
		}
	})
}
