// FileUtil
package main

import (
	"bytes"
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
	Path         string     `json:"path"`
	RelativePath string     `json:"relaPath"`
	Name         string     `json:"Node"`
	Leafs        []FileTree `json:"Nodes,omitempty"`
	IsDir        bool       `json:"isDir"`
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
func GetCurFileList(path string, isSkip bool, idx int) ([]FileTree, bool, error) {
	i := strings.LastIndex(path, "/")
	length := len(path)
	if i != length-1 {
		i = strings.LastIndex(path, "\\")
	}
	if i != length-1 {
		path = path + "/"
		length++
	}
	relativePath := path[idx:length]

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
			fileInfo := path + v.Name() + "/README.md"
			os.OpenFile(fileInfo, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			hasNoneDir = false
			fileTree := FileTree{path + v.Name(), relativePath + v.Name() + "/README.md", v.Name(), nil, true}
			slice1 = append(slice1, fileTree)
		}
	}
	for _, v := range fileInfoList {
		if isSkip && strings.Compare(v.Name()[0:1], ".") == 0 {
			continue
		}
		if !v.IsDir() {
			fileTree := FileTree{path + v.Name(), relativePath + v.Name(), v.Name(), nil, false}
			slice1 = append(slice1, fileTree)
		}
	}
	return slice1, hasNoneDir, nil
}

func GetAllFileList(path string, isSkip bool, idx int) ([]FileTree, error) {
	fileTrees, isOver, e := GetCurFileList(path, isSkip, idx)
	if !isOver {
		for i, v := range fileTrees {
			if v.IsDir {
				fileTrees[i].Leafs, e = GetAllFileList(v.Path, true, idx)
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

func dealGitBookTree(fileTree []FileTree, headStr string) string {
	var buffer bytes.Buffer
	for _, v := range fileTree {
		if v.IsDir {
			str := dealGitBookTree(v.Leafs, "  "+headStr)
			buffer.WriteString(str)
		} else {
			buffer.WriteString(headStr)
			arr := strings.Split(v.Name, ".")
			// 图片暂时不处理
			if len(arr) > 1 && strings.Contains("png jpg gif bmp", strings.ToLower(arr[1])) {
				continue
			}
			buffer.WriteString("[" + v.Name + "](" + v.RelativePath + ")\n")
		}
	}
	return buffer.String()
}

func genGitBookTree(path string, headStr string) {
	fileTree, _ := GetAllFileList(path, true, len("F:/src/gitNote/"))
	jsonBytes, _ := json.Marshal(fileTree)
	jsonStr := string(jsonBytes)
	fmt.Println(jsonStr)

	rets := dealGitBookTree(fileTree, "- ")
	fileInfoName := path + "/README.md"
	os.OpenFile(fileInfoName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fileName := path + "/SUMMARY.md"
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	_, err = file.WriteString(rets)

	checkErr(err)
}

func main() {
	path := "F:/src/gitNote/"
	genGitBookTree(path, "- ")
	fmt.Println("---------------------")

	// var mapResult map[string]interface{}
	// err := json.Unmarshal(jsonBytes, &mapResult)
	// checkErr(err)
	// fmt.Println(mapResult)
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
