package main

import (
	"fmt"
	"gofile/com"
	"gofile/config"
	"gofile/handler"
	"gofile/msg"
	"gofile/util"
	"io/ioutil"
	"os"

	// "path/filepath"

	log "github.com/donnie4w/go-logger/logger"
)

var hlr handler.Handler

func Opencom(comnum string, baudrate int) bool {
	irw, err := com.New(comnum, uint(baudrate))
	if err != nil {
		return false
	}

	hlr = handler.Handler{
		Rwc:    irw,
		Listch: make(chan []string, 10),

		Chandler: &DefaultComtask,
	}
	log.Debug("Opencom:", hlr)
	go hlr.HandleLoop()
	go Run()
	return true
}

func Sendmsg(message msg.Msg) {
	bs, err := msg.Pack(message)
	if err != nil {
		log.Error("sendmsg err:", err)
		return
	}
	hlr.Send(bs)
}

func Browseclientpath(bpath string) byte {

	curpath := config.Cfg.Section("file1").Key("clientpath").MustString(config.GetRootdir())
	log.Debug(curpath)

	if bpath != "" {
		if string(curpath[len(curpath)-1]) != "/" {
			curpath += `/` + bpath
		} else {
			curpath += bpath
		}

	}
	log.Debug(curpath)

	s, err := os.Stat(curpath)
	if err != nil {
		log.Error(err)
		log.Debug(config.GetRootdir())
		config.Cfg.Section("file1").Key("clientpath").SetValue(config.GetRootdir())

		config.Save()
		return 3
	}
	if s.IsDir() {
		files, _ := ioutil.ReadDir(curpath)
		jsStr1 := fmt.Sprintf(`$('#clientpath').val("%s");$("#clientfiles").find("li").remove()`, curpath)
		Defaultweb.UI.Eval(jsStr1)
		for _, f := range files {
			log.Debug(f.Name())

			jsStr := fmt.Sprintf(`$('#clientfiles').append("<li>%s</li>")`, f.Name())
			Defaultweb.UI.Eval(jsStr)
		}
		config.Cfg.Section("file1").Key("clientpath").SetValue(curpath)

		config.Save()
		return 0
	} else {
		// files, _ := ioutil.ReadDir(curpath)
		// jsStr1 := `$("#filesgroup").find("li").remove()`
		// Defaultweb.UI.Eval(jsStr1)
		// for _, f := range files {
		// 	log.Debug(f.Name())

		// 	jsStr := fmt.Sprintf(`$('#filesgroup').append("<li>%s</li>")`, f.Name())
		// 	Defaultweb.UI.Eval(jsStr)
		// }
		return 1
	}

}

func Browseclientuppage() {
	curpath := config.Cfg.Section("file1").Key("clientpath").MustString(config.GetRootdir())
	log.Debug(curpath)

	curpath = util.GetParentDirectory(curpath)

	_, err := os.Stat(curpath)
	if err != nil {
		log.Error(err)
		return
	}

	files, _ := ioutil.ReadDir(curpath)
	jsStr1 := fmt.Sprintf(`$('#clientpath').val("%s");$("#clientfiles").find("li").remove()`, curpath)

	Defaultweb.UI.Eval(jsStr1)
	for _, f := range files {
		log.Debug(f.Name())

		jsStr := fmt.Sprintf(`$('#clientfiles').append("<li>%s</li>")`, f.Name())
		Defaultweb.UI.Eval(jsStr)
	}
	config.Cfg.Section("file1").Key("clientpath").SetValue(curpath)

	config.Save()

}

func Run() {
	select {
	case listdir := <-hlr.Listch:
		log.Debug(listdir)
		jsStr1 := `$("#downfilesgroup").find("li").remove()`
		Defaultweb.UI.Eval(jsStr1)
		for _, f := range listdir {
			log.Debug(f)

			jsStr := fmt.Sprintf(`$('#downfilesgroup').append("<li>%s</li>")`, f)
			Defaultweb.UI.Eval(jsStr)
		}

	}

}
