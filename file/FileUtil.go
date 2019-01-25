// FileUtil
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/**
 * 获取当前执行文件（exe）的路径
 */
func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

/**
 * 获取目录下所有文件、文件夹
 */
func getFileList(path string) {
	// files, _ := filepath.Glob("*")

}

func main() {
	path := "F:/src/gitNote"
	// 读取目录列表
	fileInfoList, e := ioutil.ReadDir(path)
	if e != nil {
		fmt.Println("read dir error")
		return
	}
	// for i, v := range fileInfoList {
	// 	fmt.Println(i, "=", v.Name(), v.IsDir())
	// }
	fmt.Println("---------")
	json1, _ := json.Marshal(fileInfoList)
	fmt.Println(string(json1))
	map1 := map[string]interface{}{"name": "John", "age": 10}
	json1, _ = json.Marshal(map1)
	fmt.Println(string(json1))

	fileName := path + "/test.md"
	// 读取文件
	fileInfo, e := os.Open(fileName, os.O_RDWR, os.ModeDir)
	if e != nil {
		fmt.Println(fileName + "doesn't exist!")
		// 不存在，则创建文件
		fileInfo, e = os.Create(fileName)
		if e != nil {
			fmt.Println("create dir error")
			return
		}
	}

	// 写入内容（末尾追加），返回字节大小n
	n, err := fileInfo.Write([]byte("白日依山尽"))
	fmt.Println("write dir", n)
	checkErr(err)
	n, err = fileInfo.Write([]byte("黄河入海流"))
	fmt.Println("write dir", n)
	checkErr(err)
	n, err = fileInfo.WriteString("孤帆远影碧空尽")
	fmt.Println("write dir", n)
	checkErr(err)
	// 指定位置追加写入内容，会覆盖后面的内容
	bytesContent := []byte("李白\r\n")
	// size := len(bytesContent)
	n, err = fileInfo.WriteAt(bytesContent, 4)
	fmt.Println("write dir", n)
	checkErr(err)
}

/**
 * 默认return
 */
func checkErr(e error, isRet ...bool) {
	isReturn := len(isRet) == 0 || isRet[0]
	if e != nil {
		if !isReturn {
			panic(e)
			return
		}
	}
}
