package public

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func RunCMD(cmd string) string {
	ret := ""
	//MyLog("run ->" + cmd)

	result, err := Exec_shell(cmd)
	if err != nil {
		ret += err.Error()
	} else {
		ret += result
	}

	return ret
}

var cmd *exec.Cmd

//阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出
func Exec_shell(s string) (string, error) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序

	if string(os.PathSeparator) == "\\" {
		cmd = exec.Command("cmd", "-c", s)
	} else {
		cmd = exec.Command("/bin/bash", "-c", s)
	}
	if 1 == 2 {
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + string(output))
			return "", err
		}
		fmt.Println(string(output))
		return string(output), nil
	} else {

		//}
		//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
		var out bytes.Buffer
		cmd.Stdout = &out

		//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
		err := cmd.Run()
		CheckErr(err)
		MyLog("ret=" + out.String())
		return out.String(), err

	}

}
func MyLog(a ...interface{}) (n int, err error) {

	//return public.Log(a...)
	return 1, nil
}

//错误处理函数
func CheckErr(err error) {
	if err != nil {
		MyLog("error")
		MyLog(err.Error())

	}
}
