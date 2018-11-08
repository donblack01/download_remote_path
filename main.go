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
		}
		return
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
	_, err = os.Stat(downPath)
	if os.IsNotExist(err) {
		os.Mkdir(downPath, os.ModePerm)
	}
	if len(nameSlice) != 1 {
		newPath := downPath + "/" + nameSlice[0] + "." + nameSlice[1]
		out, _ := os.Create(newPath)
		io.Copy(out, bytes.NewReader(content))
	} else {

		reg, _ := regexp.Compile(`<li><a href="([^ \f\n\r\t\v]*)"> (?:[^ \f\n\r\t\v]*)</a></li>`)
		regName := reg.FindAllStringSubmatch(string(content), -1)
		for _, v := range regName {
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
