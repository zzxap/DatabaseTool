package public

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/otiai10/copy"

	//"strconv"
	//"net/http/cookiejar"
	"io/ioutil"
	//"log"
	//"path/filepath"
	//"path"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	//"github.com/kardianos/osext"
	"archive/zip"
	//"math/rand"

	"github.com/go-ini/ini"
	//"github.com/op/go-logging"
)

var IsShowLog = true
var CurPath string

func GetOSName() string {
	return runtime.GOOS

	//mac ->'darwin'    windows -->windows  linux -->linux
}
func GetLanguage() string {
	// Check the LANG environment variable, common on UNIX.
	// XXX: we can easily override as a nice feature/bug.
	envlang, ok := os.LookupEnv("LANG")
	if ok {
		return strings.Split(envlang, ".")[0]
	}

	// Exec powershell Get-Culture on Windows.
	cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
	output, err := cmd.Output()
	if err == nil {
		return strings.Trim(string(output), "\r\n")
	}

	return ""
}

func GetMd5(str string) string {

	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	//fmt.Println(md5str1)
	return md5str1
}
func GetGoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	if strings.Contains(gopath, ";") {

		arr := strings.Split(gopath, ";")
		gopath = arr[1]
	}

	// fmt.Println(gopath)
	return gopath
}

func GetCurRunPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}
func Unzip(src_zip string) string {
	// 解析解压包名
	dest := strings.Split(src_zip, ".")[0]
	// 打开压缩包
	unzip_file, err := zip.OpenReader(src_zip)
	if err != nil {
		return "压缩包损坏"
	}
	// 创建解压目录
	os.MkdirAll(dest, 0755)
	// 循环解压zip文件
	for _, f := range unzip_file.File {
		rc, err := f.Open()
		if err != nil {
			return "压缩包中文件损坏"
		}
		path := filepath.Join(dest, f.Name)
		// 判断解压出的是文件还是目录
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			// 创建解压文件
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return "创建本地文件失败"
			}
			// 写入本地
			_, err = io.Copy(f, rc)
			if err != nil {
				if err != io.EOF {
					return "写入本地失败"
				}
			} else {
				return "success"
			}
			f.Close()
		}
	}
	unzip_file.Close()
	return "OK"
}

func UnzipToest(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			Log(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				Log(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					Log(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
func NetWorkAvailable() bool {
	cmd := exec.Command("ping", "baidu.com", "-c", "1", "-W", "5")
	fmt.Println("NetWorkStatus Start:", time.Now().Unix())
	err := cmd.Run()
	fmt.Println("NetWorkStatus End  :", time.Now().Unix())
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		fmt.Println("Net Status , OK")
	}
	return true
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		Log(err.Error())
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		Log(err.Error())
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		Log(err.Error())
		return err
	}
	if ExistsPath(dst) {
		Log("copy success" + dst)
	} else {
		Log("copy fail" + dst)
	}

	return out.Close()
}

//拷贝文件  要拷贝的文件路径 拷贝到哪里 "github.com/otiai10/copy"
func CopyFiles(source, dest string) bool {
	if source == "" || dest == "" {
		Log("source or dest is null")
		return false
	}
	if ExistsPath(source) {
		err := copy.Copy(source, dest)
		if err == nil {
			return true
		} else {
			return false
		}
	} else {
		return false
	}

	/*
		//打开文件资源
		source_open, err := os.Open(source)
		//养成好习惯。操作文件时候记得添加 defer 关闭文件资源代码
		if err != nil {
			Log(err.Error())
			return false
		}
		defer source_open.Close()
		//只写模式打开文件 如果文件不存在进行创建 并赋予 644的权限。详情查看linux 权限解释
		dest_open, err := os.Open(dest) // os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 644)
		if err != nil {
			Log(err.Error())
			return false
		}
		//养成好习惯。操作文件时候记得添加 defer 关闭文件资源代码
		defer dest_open.Close()
		//进行数据拷贝
		_, copy_err := io.Copy(dest_open, source_open)
		if copy_err != nil {
			Log("copy fail")
			Log(copy_err.Error())
			return false
		} else {
			Log("copy success")
			return true
		}
	*/
}

func GetCurDir() string {
	if CurPath != "" {
		fmt.Printf("CurPath not null")

		return CurPath
	} else {
		dir, _ := GetCurrentPath()
		return dir
	}

}
func GetCurrentPath() (dir string, err error) {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		Log("exec.LookPath(%s), err: %s\n", os.Args[0], err)
		return "", err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		Log("filepath.Abs(%s), err: %s\n", path, err)
		return "", err
	}
	dir = filepath.Dir(absPath)
	return dir, nil
}

func CreateFile(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	fmt.Println("create success " + filepath)
	return nil

}
func ExistsPath(fullpath string) bool {
	//dir, _ := GetCurrentPath() //os.Getwd() //当前的目录
	//fullpath := dir + "/" + path
	_, err := os.Stat(fullpath)

	//Log("fullpath==" + fullpath)
	return err == nil || os.IsExist(err)
}

func CreatePath(fullpath string) {
	//dir, _ := GetCurrentPath() //os.Getwd() //当前的目录
	//fullpath := dir + "/" + newPath
	//fullpath = strings.Replace(fullpath, "/", "\\", -1)
	//fullpath = strings.Replace(fullpath, " ", "", -1)

	//newPath = strings.Replace(newPath, " ", "", -1)
	ff, errr := os.Stat(fullpath)
	if errr != nil && os.IsNotExist(errr) {
		fmt.Println(ff, fullpath+" 文件不存在 创建") //为什么打印nil 是这样的如果file不存在 返回f文件的指针是nil的 所以我们不能使用defer f.Close()会报错的

		var path string
		if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
			path = "\\"
		} else {
			path = "/"
		}
		fmt.Println(path)

		if err := os.MkdirAll(fullpath, 0777); err != nil {
			if os.IsPermission(err) {
				fmt.Println("你不够权限创建文件")
			}

		} else {
			fmt.Println("创建目录" + fullpath + "成功")
		}

		//err := os.Mkdir(fullpath, os.ModePerm) //在当前目录下生成md目录
		//if err != nil {
		//	Log(err)
		//}

	} else {
		//Log(ff, fullpath+"文件存在 ")
	}

}

func SetCookie(r *http.Request, name string, value string) {
	COOKIE_MAX_MAX_AGE := time.Hour * 24 / time.Second // 单位：秒。
	maxAge := int(COOKIE_MAX_MAX_AGE)

	uid_cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: false,
		MaxAge:   maxAge}

	r.AddCookie(uid_cookie)

}
func GetTotal(price string, num string) string {
	//Log("buyprice=" + price + "num=" + num)
	fPrice, err1 := strconv.ParseFloat(price, 64)
	fnum, err2 := strconv.ParseFloat(num, 64)
	if err1 == nil && err2 == nil {
		return fmt.Sprintf("%1.2f", fPrice*fnum)
	} else {
		if err1 != nil {
			Log(err1)
		}
		if err2 != nil {
			Log(err2)
		}
	}
	return ""

}

func RemoveFile(path string) bool {
	//Log("upload picture Task is running...")
	curdir := GetCurDir()
	fullPath := curdir + "/" + path + "/"
	err := os.Remove(fullPath) //删除文件test.txt
	if err != nil {

		Log("remove fail " + fullPath)
		return false
	} else {
		//如果删除成功则输出 file remove OK!
		return true
	}

}
func SavePictureTask(res http.ResponseWriter, req *http.Request, path string, userid string, typeid string) string {
	//Log("upload picture Task is running...")
	curdir := GetCurDir()
	PthSep := string(os.PathSeparator)
	var fileNames string = "#"
	if req.Method == "GET" {

	} else {

		ff, errr := os.Open(curdir + PthSep + path + PthSep)
		if errr != nil && os.IsNotExist(errr) {
			Log(ff, ""+path+"文件不存在,创建") //为什么打印nil 是这样的如果file不存在 返回f文件的指针是nil的 所以我们不能使用defer f.Close()会报错的
			CreatePath(curdir + PthSep + path + PthSep)

		}

		var (
			status int
			err    error
		)
		defer func() {
			if nil != err {
				http.Error(res, err.Error(), status)
			}
		}()
		// parse request
		const _24K = (1 << 20) * 24
		if err = req.ParseMultipartForm(_24K); nil != err {
			status = http.StatusInternalServerError
			return ""
		}
		for _, fheaders := range req.MultipartForm.File {
			for _, hdr := range fheaders {
				// open uploaded
				var infile multipart.File
				if infile, err = hdr.Open(); nil != err {
					status = http.StatusInternalServerError
					return ""
				}
				filename := hdr.Filename
				arr := strings.Split(filename, ".")
				if len(arr) > 1 {
					filename = GetRandom() + "." + arr[len(arr)-1]
				}

				if strings.Contains(strings.ToLower(filename), ".mp3") || strings.Contains(strings.ToLower(filename), ".mov") {
					//如果是音频文件，直接存到picture文件夹，不存temp文件夹
					path = "picture" + PthSep + userid + PthSep + typeid
					CreatePath(curdir + PthSep + path + PthSep)
				}

				// open destination
				var outfile *os.File
				savePath := curdir + PthSep + path + PthSep + filename
				if outfile, err = os.Create(savePath); nil != err {
					status = http.StatusInternalServerError
					return ""
				}
				// 32K buffer copy
				//var written int64
				if _, err = io.Copy(outfile, infile); nil != err {
					status = http.StatusInternalServerError
					return ""
				}

				infile.Close()
				outfile.Close()
				//CreatePath(curdir + "/" + path + "/thumbnail")
				//ImageFile_resize(infile, curdir+"/"+path+"/thumbnail/"+filename, 200, 200)
				fileNames += "," + filename
				//outfile.Close()
				//res.Write([]byte("uploaded file:" + filename + ";length:" + strconv.Itoa(int(written))))
			}
		}
	}
	fileNames = strings.Replace(fileNames, "#,", "", -1)
	fileNames = strings.Replace(fileNames, "#", "", -1)
	return fileNames
}

func copyFile(source, dest string) bool {
	if source == "" || dest == "" {
		Log("source or dest is null")
		return false
	}

	source_open, err := os.Open(source)

	if err != nil {
		Log(err.Error())
		return false
	}
	defer source_open.Close()

	dest_open, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 666)
	if err != nil {
		Log(err.Error())
		return false
	}

	defer dest_open.Close()

	_, copy_err := io.Copy(dest_open, source_open)
	if copy_err != nil {
		Log(copy_err.Error())
		return false
	} else {
		return true
	}
}

func GetMapByJsonStr(jsonstr string) map[string]interface{} {
	//var jsonstr:='{\"data\": { \"mes\": [ {\"fromuserid\": \"25\", \"touserid\": \"56\",\"message\": \"hhhhhaaaaaa\",\"time\": \"2017-12-12 12:11:11\"}]}}';
	if len(jsonstr) > 4 {
		var d map[string]interface{}
		err := json.Unmarshal([]byte(jsonstr), &d)
		if err != nil {
			Log(err)
			Log("bad json")
			return nil
		}
		return d
	}
	return nil
}

func GetMessageMapByJson(jsonstr string) map[string]interface{} {
	//var jsonstr:='{\"data\": { \"mes\": [ {\"fromuserid\": \"25\", \"touserid\": \"56\",\"message\": \"hhhhhaaaaaa\",\"time\": \"2017-12-12 12:11:11\"}]}}';

	if len(jsonstr) > 4 && strings.Index(jsonstr, "{") > -1 && strings.Index(jsonstr, "}") > -1 {
		mapp := GetMapByJsonStr(jsonstr)
		//Log(mapp)
		mappp := mapp["data"]
		//Log(mappp)
		kll := mappp.(map[string]interface{})["mes"]
		//Log(kll)
		mymap := kll.(map[string]interface{})
		//Log(mymap["fromuserid"])

		return mymap
	}
	return nil
}

func GetJsonStrByMap(MapList map[int]map[string]string) string {
	var str string = "##"
	for _, v := range MapList {
		jsonStr, err := json.Marshal(v)
		if err != nil {
			Log(err)
		}
		//Log("map to json", string(str))
		str += "," + string(jsonStr)
	}
	str = strings.Replace(str, "##,", "", -1)
	str = strings.Replace(str, "+0000 UTC", "", -1)
	str = strings.Replace(str, " +0000", "", -1)
	str = strings.Replace(str, "##", "", -1)
	return str
}
func ConverToStr(v interface{}) string {
	str := ""
	if v == nil {
		return ""
	} else {

		switch t := v.(type) {

		case int:
			str = string(v.(int))
		case int64:
			str = string(v.(int64))
		case float64:
			str = fmt.Sprintf("%f", v)
		case float32:
			str = fmt.Sprintf("%f", v)
		case string:
			str = v.(string)
		case []byte:
			str = string(v.([]byte))
		case time.Time:
			str = fmt.Sprintf("%s", v)
		default:
			Log(t)
			str = v.(string)
		case bool:
			str = v.(string)
		}

	}

	return strings.Replace(str, ".000000", "", -1)
}
func GetCurDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
func GetTimeStr(str string) string {
	//postgresql 2019-12-12 不能插入Time类型的字段，需要加00:00:00
	if str == "" {
		return "00:00:00"
	}
	str = strings.Replace(str, "'", "", -1)
	str = strings.Replace(str, "exec ", "", -1)

	if !strings.Contains(str, ":") {

		//str += " 00:00:00"
	}
	return str
}
func ReplaceStr(str string) string {
	str = strings.Replace(str, "'", "", -1)
	str = strings.Replace(str, "exec ", "", -1)
	return str //.Replace(str, ",", "", -1).Replace(str, "-", "\-", -1) //-1表示替换所有
}
func GetCurDay() string {

	return time.Now().Format("2006-01-02")
}
func GetCurMinutes() string {

	return time.Now().Format("2006-01-02 15:04")
}

func GetOrderNum() string {
	day := time.Now().AddDate(0, 0, 0).Format("2006-01-02 15:04:05")
	day = strings.Replace(day, "-", "", -1)
	day = strings.Replace(day, " ", "", -1)
	day = strings.Replace(day, ":", "", -1)
	return day
}

var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xForwardedFor2 = http.CanonicalHeaderKey("x-forwarded-for")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")
var xRealIP2 = http.CanonicalHeaderKey("x-real-ip")
var xRealIP3 = http.CanonicalHeaderKey("x-real-client-ip")

var ProxyClientIP = http.CanonicalHeaderKey("Proxy-Client-IP")
var WLProxyClientIP = http.CanonicalHeaderKey("WL-Proxy-Client-IP")
var HTTPXFORWARDEDFOR = http.CanonicalHeaderKey("HTTP_X_FORWARDED_FOR")

func RealIP(r *http.Request) string {

	//PrintHead(r)

	var ip string

	//clientIP := realip.FromRequest(r)
	//log.Println("GET / from", clientIP)

	if xff := r.Header.Get(xForwardedFor); xff != "" {
		//Log(xff)
		i := strings.Index(xff, ", ")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	} else if xff := r.Header.Get(xForwardedFor2); xff != "" {
		//Log(xff)
		i := strings.Index(xff, ", ")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else if xrip := r.Header.Get(xRealIP2); xrip != "" {
		ip = xrip
	} else if xrip := r.Header.Get(xRealIP3); xrip != "" {
		ip = xrip
	} else if xrip := r.Header.Get(ProxyClientIP); xrip != "" {
		ip = xrip
	} else if xrip := r.Header.Get(WLProxyClientIP); xrip != "" {
		ip = xrip
	} else {
		ip = r.RemoteAddr
	}

	return ip

	//return realip.FromRequest(r)
}
func GetNameSinceNow(after int) string {
	day := time.Now().AddDate(0, 0, after).Format("2006-01-02")
	day = strings.Replace(day, "-", "", -1)
	return day
}
func GetDaySinceNow(after int) string {
	return time.Now().AddDate(0, 0, after).Format("2006-01-02")
}
func GetMonthSinceNow(after int) string {
	return time.Now().AddDate(0, after, 0).Format("2006-01")
}
func GetYearSinceNow(after int) string {
	return time.Now().AddDate(after, 0, 0).Format("2006-01")
}

var logfile *os.File
var oldFileName string

func Log(a ...interface{}) (n int, err error) {
	//log.SetFlags(log.LstdFlags | log.Lshortfile)

	if EnableFmtLog() {
		log.Println(a...)

	}
	if EnableLog() {
		name := GetCurDay()
		name = strings.Replace(name, "-", "", -1)
		if logfile == nil || name != oldFileName {
			curpath := GetCurDir()
			PthSep := string(os.PathSeparator)
			path := curpath + PthSep + "log" + PthSep + name + ".txt"
			oldFileName = name
			if !ExistsPath(path) {

				os.Create(path)
			}

			logfile, _ = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

			//defer logfile.Close()
			//删除7天前的日志文件
			name_delete := GetNameSinceNow(-7)
			path_delete := curpath + PthSep + "log" + PthSep + name_delete + ".txt"
			if ExistsPath(path_delete) {
				os.Remove(path_delete)
			}

		}
		log.SetOutput(logfile)
		log.Println(a...)

	}
	return 1, nil
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP22() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipstr := ipnet.IP.String()
				index := strings.Index(ipstr, "127.0")
				if index > -1 {
					continue
				}

				index = strings.Index(ipstr, "192.168.")
				if index > -1 {
					return ipstr
					break
				}

				index = strings.Index(ipstr, "169.254.")
				if index > -1 {
					continue
				}

				return ipstr
			}
		}
	}
	return ""
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err == nil {
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)

		return localAddr.IP.String()
	} else {
		return GetLocalIPP()
	}

}
func GetLocalIPP() string {
	//GetIpList()
	var ipstr string = ""
	//windows 获取IP
	host, _ := os.Hostname()
	addrss, err := net.LookupIP(host)
	if err != nil {
		Log("error", err.Error())
		//return ""
	}
	var ipArray []string

	for _, addr := range addrss {
		if ipv4 := addr.To4(); ipv4 != nil {
			Log("ippppp=: ", ipv4)
			ipstr = ipv4.String()

			if !strings.HasPrefix(ipstr, "127.0") && !strings.HasPrefix(ipstr, "169.254") && !strings.HasPrefix(ipstr, "172.16") {
				ipArray = append(ipArray, ipstr)
			}
		}
	}

	//提取公网IP
	//var pubIpArray []string
	for i := 0; i < len(ipArray); i++ {
		//Log("pubip===" + ipArray[i])
		if !strings.HasPrefix(ipArray[i], "10.") && !strings.HasPrefix(ipArray[i], "192.168") && !strings.HasPrefix(ipArray[i], "172.") {
			return ipArray[i]
			//pubIpArray = append(pubIpArray, ipstr)
		}
	}

	//如果没有公网IP 就返回一个本地IP

	if len(ipArray) > 0 {

		return ipArray[0]
	}

	//linux 获取IP
	if ipstr == "" {

		ifaces, errr := net.Interfaces()
		// handle err
		if errr != nil {
			Log("error", errr.Error())
			return ""
		}

		for _, i := range ifaces {
			addrs, _ := i.Addrs()
			// handle err
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				// process IP address
				//Log("ip=", ip)
				ipstr = fmt.Sprintf("%s", ip)
				Log("ipstr=", ipstr)

				index := strings.Index(ipstr, "127.0")
				if index > -1 {
					continue
				}
				index = strings.Index(ipstr, "192.168.")
				if index > -1 {
					return ipstr
					break
				}

				index = strings.Index(ipstr, "169.254.")
				if index > -1 {
					continue
				}
				if len(ipstr) > 6 {
					array := strings.Split(ipstr, ".")
					if len(array) == 4 {
						return ipstr
					}

				}

			}
		}

	}
	Log("貌似获取不到IP，请确保你的电脑连接了Wi_Fi,或是公网服务器 ，\n  can not get your ip，make sure your PC connect WIFI,or is a public internet server \n")
	//Log("暂时不支持安装在广域网只支持 192.168开头的内网IP，\n 如需要安装互联网版本请联系微信 ELink9988 \n")
	//Log("this version only suport 192.168.x.x \n if you need to suport internet please  contact ELink9988 \n")
	return ""
}
func HttpPostJson(url string, json string) string {
	//Log("url=" + url + " paras=" + paras)
	client := &http.Client{}
	var jsonStr = []byte(json)
	req, err := http.NewRequest("POST",
		url,
		bytes.NewBuffer(jsonStr))

	if err != nil {
		// handle error
		return ""
	}

	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return ""
	}

	//Log(string(body))
	return string(body)
}
func HttpPost(url string, paras string) string {
	//Log("url=" + url + " paras=" + paras)
	client := &http.Client{}

	req, err := http.NewRequest("POST",
		url,
		strings.NewReader(paras))

	if err != nil {
		// handle error
		return ""
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return ""
	}

	//Log(string(body))
	return string(body)
}

func HttpGet(url string) string {
	//Log("get =" + url)
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		Log(err.Error())
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		Log(err.Error())
		return ""
	}

	//Log("response =" + string(body))
	if strings.Contains(string(body), "context.Request") {
		return "error"
	} else {
		return string(body)
	}

}
func GetRandom() string {

	return GetOrderNum()

}
func BytesToString(b []byte) string {
	str := hex.EncodeToString(b[:])
	str = strings.Replace(str, "01", "1", -1)
	str = strings.Replace(str, "02", "2", -1)
	str = strings.Replace(str, "03", "3", -1)
	str = strings.Replace(str, "04", "4", -1)
	str = strings.Replace(str, "05", "5", -1)
	str = strings.Replace(str, "06", "6", -1)
	str = strings.Replace(str, "07", "7", -1)
	str = strings.Replace(str, "08", "8", -1)
	str = strings.Replace(str, "09", "9", -1)
	str = strings.Replace(str, "00", "0", -1)
	return str
}

func SaveUploadFileask(w http.ResponseWriter, r *http.Request, savePath string) string {
	req := r

	if req.MultipartForm == nil || len(req.MultipartForm.File) == 0 {
		Log("no file upload")
		return ""
	} else {
		Log("has file upload")
	}
	var fileNames string = "#"
	var (
		//status int
		err error
	)
	//status := 0
	defer func() {
		if nil != err {
			Log(err.Error())

		}
	}()
	//PthSep := string(os.PathSeparator)
	// parse request
	const _24K = (1 << 20) * 24
	if err = req.ParseMultipartForm(_24K); nil != err {
		//status = http.StatusInternalServerError
		return ""
	}

	for _, fheaders := range req.MultipartForm.File {
		for _, hdr := range fheaders {
			// open uploaded
			var infile multipart.File
			if infile, err = hdr.Open(); nil != err {
				//status = http.StatusInternalServerError
				return ""
			}
			filename := hdr.Filename
			var outfile *os.File

			if outfile, err = os.Create(savePath + filename); nil != err {

				return ""
			}

			if _, err = io.Copy(outfile, infile); nil != err {

				return ""
			}
			infile.Close()
			outfile.Close()
			fileNames += "$" + filename

		}
	}
	fileNames = strings.Replace(fileNames, "#$", "", -1)
	fileNames = strings.Replace(fileNames, "#", "", -1)

	return fileNames
}
func HttpDownloadFile(url string, toPath string) error {
	Log("get =" + url)
	res, err := http.Get(url)
	if err != nil {
		Log(err)
		return err
	}
	f, err := os.Create(toPath)
	defer f.Close()
	if err != nil {
		Log(err)
		return err
	}
	Log("copy ")
	_, errr := io.Copy(f, res.Body)
	return errr
	//Log("size =" + size)
}

//http中转传输
func TransHttpPostTask(w http.ResponseWriter, r *http.Request) {
	urlpath := r.URL.Path

	//Log("urlpath=" + urlpath)
	//userid := r.FormValue("userid")
	//pid := r.FormValue("pid")
	//str := dal.GetProductTypeTwoList(userid, pid, "")

	r.ParseForm()
	//log.Println(r.Form)
	var paras string = "##"
	for key, value := range r.Form {
		//log.Println(key)
		//log.Println(value[0])
		paras += "&" + string(key) + "=" + string(value[0])
	}
	paras = strings.Replace(paras, "##&", "", -1)
	//Log(paras)

	str := HttpPost("http://www.iosbuy.com"+urlpath, paras+"&from=invoiceServer")
	//Log("login result=" + str)
	w.Write([]byte(str))

}

func TransHttpGetTask(w http.ResponseWriter, r *http.Request) string {
	urlpath := r.URL.Path

	//Log("urlpath=" + urlpath)
	//userid := r.FormValue("userid")
	//pid := r.FormValue("pid")
	//str := dal.GetProductTypeTwoList(userid, pid, "")

	r.ParseForm()
	//log.Println(r.Form)
	var paras string = "##?"
	for key, value := range r.Form {
		//log.Println(key)
		//log.Println(value[0])
		paras += "&" + string(key) + "=" + string(value[0])
	}
	paras = strings.Replace(paras, "##?&", "?", -1)
	//Log(paras)

	str := HttpGet("http://www.xx.com" + urlpath + paras + "&from=invoiceServer")
	//Log("login result=" + str)
	//w.Write([]byte(str))
	return str

}

func GetLogoPath() string {
	cfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, "config.ini")
	if err == nil {
		//fromfolder := cfg.Section("path").Key("fromfolder").String()
		//tofolder := cfg.Section("path").Key("tofolder").String()
		return cfg.Section("path").Key("logopath").String()

	}
	return ""
}
