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

type FileTree struct {
	Path  string
	Name  string
	Leafs []FileTree
	IsDir bool
}

/**
 * 获取当前执行文件（exe）的路径
 */
func GetCurrentPath() (string, error) {
	// strFiles, _ := filepath.Glob("*")
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
 * 获取当前目录下所有文件、文件夹，按文件夹、命名排序；
 * isSkip表示跳过'.git'点开头的文件夹或文件
 */
func GetCurFileList(path string, isSkip bool) ([]FileTree, bool, error) {
	i := strings.LastIndex(path, "/")
	length := len(path)
	if i != length-1 {
		i = strings.LastIndex(path, "\\")
	}
	if i != length-1 {
		path = path + "/"
	}

	// 读取目录列表
	fileInfoList, e := ioutil.ReadDir(path)
	// 如果为true代表没有文件夹，可以直接结束
	hasNoneDir := true
	if e != nil {
		fmt.Println("read dir error!", e)
		return nil, hasNoneDir, e
	}
	slice1 := make([]FileTree, 0, len(fileInfoList))
	for _, v := range fileInfoList {
		if isSkip && v.Name()[0] == '.' {
			continue
		}
		if v.IsDir() {
			hasNoneDir = false
			fileTree := FileTree{path + v.Name(), v.Name(), nil, true}
			slice1 = append(slice1, fileTree)
		}
	}
	for _, v := range fileInfoList {
		if isSkip && strings.Compare(v.Name()[0:1], ".") == 0 {
			continue
		}
		if !v.IsDir() {
			fileTree := FileTree{path + v.Name(), v.Name(), nil, false}
			slice1 = append(slice1, fileTree)
		}
	}
	return slice1, hasNoneDir, nil
}

func GetAllFileList(path string, isSkip bool) ([]FileTree, error) {
	fileTrees, isOver, e := GetCurFileList(path, isSkip)
	if !isOver {
		for i, v := range fileTrees {
			if v.IsDir {
				fileTrees[i].Leafs, e = GetAllFileList(v.Path, true)
				if nil != e {
					break
				}
			}
		}
	}
	if e != nil {
		fmt.Println("read dir error!", e)
	}
	return fileTrees, e
}

func genGitBookTree(fileTree []FileTree) string {
	return ""
}

func main() {
	path := "F:/src/gitNote"
	lst, _ := GetAllFileList(path, true)
	fmt.Println(lst)
	fmt.Println("---------------------")
	json1, _ := json.Marshal(lst)
	fmt.Println(string(json1))
	fmt.Println("---------------------")

	/*
		obj1 := FileTree{"", "ddd1", nil, true}
		obj2 := FileTree{"", "ddd2", nil, true}
		obj3 := FileTree{"", "ddd3", []FileTree{obj1}, true}
		fmt.Println(obj1)
		obj2.leafs = []FileTree{obj1}
		obj2.leafs = append(obj2.leafs, obj1)
		obj3.leafs = []FileTree{obj1}
		obj3.leafs = append(obj3.leafs, obj1)
		fmt.Println(obj2)
		fmt.Println(obj3)
	*/
	fileName := path + "/test.md"
	// 读取文件
	fileInfo, e := os.OpenFile(fileName, os.O_RDWR, os.ModeDir)
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
