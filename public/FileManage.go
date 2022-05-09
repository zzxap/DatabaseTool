package public

import (
	//"flag"
	//"bufio"
	"fmt"
	"io/ioutil"
	"os"

	//"path/filepath"
	"strings"

	"github.com/skip2/go-qrcode"
)

func GetFileSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}
	// get the size
	return fi.Size()
}
func RenameFolder(oldpath string, newpath string) bool {
	// 重命名文件夹
	if ExistsPath(oldpath) {

		err2 := os.Rename(oldpath, newpath)
		if err2 != nil {
			return false
		} else {
			return true
		}
	} else {

		return false
	}
}

//读取文件需要经常进行错误检查，这个帮助方法可以精简下面的错误检查过程。
func check(e error) {
	if e != nil {
		Log(e.Error())
	}
}

func MakeQRCode(content string, path, name string) {

	var png []byte
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		Log(err.Error())
		return
	}
	curdir := GetCurDir()
	PthSep := string(os.PathSeparator)
	fullpath := curdir + PthSep + "picture" + PthSep
	if len(path) > 0 {
		fullpath = fullpath + path + PthSep
		if !ExistsPath(fullpath) {
			CreatePath(fullpath)
		}
	}
	if !ExistsPath(fullpath + name) {
		WriteToFile(png, fullpath+name)
	}

}

//获取指定目录下的所有文件夹，不进入下一级目录搜索
func GetFolders(dirPth string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		fmt.Println("read error")
		return nil, err
	}
	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() { // 目录
			files = append(files, dirPth+PthSep+fi.Name())
		}

	}
	return files, nil
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		fmt.Println("read error")
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if len(suffix) > 0 {
			if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
				//fmt.Println("path==" + dirPth + PthSep + fi.Name())
				files = append(files, dirPth+PthSep+fi.Name())
			}
		} else {
			files = append(files, dirPth+PthSep+fi.Name())
		}

	}
	return files, nil
}

func ReadFile(path string) []byte {

	fi, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer fi.Close()

	fd, _ := ioutil.ReadAll(fi) //ioutil.ReadFile(path) //
	//fmt.Println(string(fd[:]))
	return fd
}
func WriteToFile(writeByte []byte, path string) {
	//var d1 = []byte(writeStr)
	//fmt.Println("path=")
	//fmt.Println(path)
	if ExistsPath(path) {
		err := os.Remove(path) //删除文件test.txt
		check(err)
	}
	err2 := ioutil.WriteFile(path, writeByte, 0666) //写入文件(字节数组)
	check(err2)

}

//加密文件
func EncryptFile(fromPath string, toPath string) {
	fileByte := ReadFile(fromPath)

	arrEncrypt, err := AesEncrypt(fileByte)

	if fromPath == toPath {
		DeleteFile(fromPath)
	}
	if err != nil {
		//fmt.Println("encrypt=" + string(err))
		return
	}
	//fmt.Println("encrypt=" + string(arrEncrypt))
	WriteToFile(arrEncrypt, toPath)

}
func DeleteFile(filePath string) {
	err := os.Remove(filePath) //删除文件
	if err != nil {
		//如果删除失败则输出 file remove Error!
		//fmt.Println("file remove Error!")
		//输出错误详细信息
		fmt.Printf("删除失败 %s", err)
	} else {
		//如果删除成功则输出 file remove OK!
		fmt.Print("file remove OK!")

	}
}

//解密文件
func DecryptFile(fromPath string, toPath string) {
	fileByte := ReadFile(fromPath)

	arrEncrypt, err := AesDecrypt(fileByte)
	if fromPath == toPath {
		DeleteFile(fromPath)
	}

	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("encrypt=" + string(arrEncrypt))
	WriteToFile(arrEncrypt, toPath)

}

//加密文件夹下面的所有文件夹的文件
func EncryptFolders(fromFolder string) {
	filePathArray, err := GetFolders(fromFolder)
	if err == nil {
		for i := 0; i < len(filePathArray); i++ {
			fullPath := filePathArray[i]
			EncryptFolder(fullPath, "")
		}
	}
}

//解密文件夹下面的所有文件夹的文件
func DecryptFolders(fromFolder string) {
	fmt.Println(fromFolder)
	filePathArray, err := GetFolders(fromFolder)
	fmt.Println(len(filePathArray))
	if err == nil {
		for i := 0; i < len(filePathArray); i++ {
			fullPath := filePathArray[i]
			DecryptFolder(fullPath, "")
		}
	}
}

//加密文件夹
func EncryptFolder(fromFolder string, typename string) {
	filePathArray, err := ListDir(fromFolder, typename)
	fmt.Println(len(filePathArray))
	PthSep := string(os.PathSeparator)
	if err == nil {
		for i := 0; i < len(filePathArray); i++ {
			fullPath := filePathArray[i]

			if strings.Contains(fullPath, PthSep) {
				array := strings.Split(fullPath, PthSep)
				filename := array[len(array)-1]
				folder := strings.Replace(fullPath, filename, "", -1)
				//toPath := folder + PthSep + filename
				toPath := folder + "Encrypt" + PthSep + filename
				if i == 0 {
					CreatePath(folder + "Encrypt")
				}
				//fmt.Println("topath=" + toPath)
				EncryptFile(fullPath, toPath)

			}

		}
	} else {
		fmt.Println("eeerrrr")
	}

}

//解密文件夹
func DecryptFolder(fromFolder string, typename string) {
	filePathArray, err := ListDir(fromFolder, typename)
	fmt.Println(len(filePathArray))
	PthSep := string(os.PathSeparator)
	if err == nil {
		for i := 0; i < len(filePathArray); i++ {
			fullPath := filePathArray[i]

			if strings.Contains(fullPath, PthSep) {
				array := strings.Split(fullPath, PthSep)
				filename := array[len(array)-1]
				folder := strings.Replace(fullPath, filename, "", -1)
				//toPath := folder + PthSep + filename
				toPath := folder + "Decrypt" + PthSep + filename
				if i == 0 {
					CreatePath(folder + "Decrypt")
				}
				fmt.Println("topath=" + toPath)
				DecryptFile(fullPath, toPath)

			}

		}
	} else {
		fmt.Println("eeerrrr")
	}

}
