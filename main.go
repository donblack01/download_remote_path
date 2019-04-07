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

	"github.com/mattn/go-gtk/gtk"
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
	gtk.Init(&os.Args)
	mainWindow := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	mainWindow.SetTitle("下载远程文件")
	mainWindow.SetIconFromFile("./images/1.png")
	mainWindow.SetSizeRequest(400, 300)
	layout := gtk.NewFixed()
	mainWindow.Add(layout)
	button := gtk.NewButtonWithLabel("确定")
	label1 := gtk.NewLabel("源地址：")
	label2 := gtk.NewLabel("保存目录：")
	label3 := gtk.NewLabel("")
	entry1 := gtk.NewEntry()
	entry2 := gtk.NewEntry()
	button.SetSizeRequest(100, 30)
	label1.SetSizeRequest(100, 30)
	entry1.SetSizeRequest(200, 30)
	label2.SetSizeRequest(100, 30)
	entry2.SetSizeRequest(200, 30)
	label3.SetSizeRequest(225, 30)
	layout.Put(button, 150, 233)
	layout.Put(label1, 23, 27)
	layout.Put(label2, 23, 99)
	layout.Put(label3, 77, 173)
	layout.Put(entry1, 140, 27)
	layout.Put(entry2, 140, 99)

	button.Clicked(func() {
		content1 := entry1.GetText()
		content2 := entry2.GetText()
		if content1 == "" || content2 == "" {
			label3.SetText("参数不能为空")
			return
		}
		label3.SetText("下载中，请稍后")
		download(content1, content2)
		label3.SetText("下载完成")
	})

	mainWindow.Connect("destroy", func() {
		gtk.MainQuit()
	})
	mainWindow.ShowAll()
	gtk.Main()
}
