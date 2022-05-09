package public

import (
	"fmt"
	"io/ioutil"
	"os"
)

//批量修改文件夹名称
func ChangeFilePath() {

	dir, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println(err)
		return
	}
	num := 0
	for _, fi := range dir {
		if fi.IsDir() { // 目录
			path := "./" + fi.Name() + "/Thumbnail"
			_, err := os.Stat(path)
			newpath := "./" + fi.Name() + "/thumbnail"
			_, errsmall := os.Stat(newpath)

			//Thumbnail thumbnail 目录同时存在
			if (err == nil || os.IsExist(errsmall)) && (err == nil || os.IsExist(err)) {
				backuppath := "./" + fi.Name() + "/thumbnail_22"
				err2 := os.Rename(newpath, backuppath)
				if err2 != nil {
					fmt.Println(err2)
					fmt.Println("111change false" + fi.Name() + "\n")
				} else {
					num++
					fmt.Println("111change success " + fi.Name() + "\n")
				}
			} else {
				fmt.Println(err)
			}

			//目录存在
			if err == nil || os.IsExist(err) {

				err2 := os.Rename(path, newpath)
				if err2 != nil {
					fmt.Println(err2)
					fmt.Println("change false" + fi.Name() + "\n")
				} else {
					num++
					fmt.Println("change success " + fi.Name() + "\n")
				}
			} else {
				fmt.Println(err)
			}
		}

	}
	fmt.Println("change finish  ")
	fmt.Println(num)
}
