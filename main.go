package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func download(path, downPath string) {
	_, err := os.Stat(downPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(downPath, os.ModePerm)
		if err != nil {
			fmt.Printf("%s", err)
			return
		}
	}
	pathSlice := strings.Split(path, "/")
	fileName := pathSlice[len(pathSlice)-1]
	if fileName == "" {
		fileName = pathSlice[len(pathSlice)-2]
	}
	nameSlice := strings.Split(fileName, ".")
	resp, err := http.Get(path)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("地址内容获取失败")
		return
	}
	content, _ := ioutil.ReadAll(resp.Body)
	nameSliceLenght := len(nameSlice)
	if nameSliceLenght > 1 {
		var newPath string
		var pathRoute string

		if nameSliceLenght > 2 {
			for i := 0; i < nameSliceLenght-1; i++ {
				pathRoute += nameSlice[i] + "."
			}
		} else {
			pathRoute = nameSlice[0] + "."
		}
		newPath = downPath + "/" + pathRoute + nameSlice[nameSliceLenght-1]

		out, _ := os.Create(newPath)
		io.Copy(out, bytes.NewReader(content))
	} else {

		reg, _ := regexp.Compile(`<a href="([^ \f\n\r\t\v]*)">([\s]?[^ \f\n\r\t\v]*)</a>`)
		regName := reg.FindAllStringSubmatch(string(content), -1)
		for _, v := range regName {
			if v[2] == "Name" || v[2] == "Last modified" || v[2] == "Size" || v[2] == "Description" || v[2] == "Parent Directory" {
				continue
			}
			download(path+v[1], downPath+"/"+nameSlice[0])
		}
	}

}
func main() {
	var path, downPath string
	fmt.Println("请输入源链接：")
	fmt.Scanln(&path)
	fmt.Println("请输入存放目录：")
	fmt.Scanln(&downPath)
	download(path, downPath)

}
