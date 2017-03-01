package main

import (
	"log"
	"net/http"
	"flag"
	"fmt"
	"strings"
)

var fileDir = flag.String("dir", ".", "server root dir")

type MuxHandle struct {
	rootDir string
	fileHandle http.Handler
}

func newMuxHandle(rootDir string) *MuxHandle {
	muxHandle := &MuxHandle{rootDir:rootDir}
	rootSystem := http.Dir(rootDir)
	muxHandle.fileHandle = http.FileServer(rootSystem)
	return muxHandle
}

func (mux *MuxHandle) ServeHTTP(rsp http.ResponseWriter,req *http.Request) {
	remoteIp := strings.Split(req.RemoteAddr, ":")[0]
	reqInfo := fmt.Sprintf("%s:%s%s\n", remoteIp,mux.rootDir, req.URL.Path)
	log.Println(reqInfo)
	mux.fileHandle.ServeHTTP(rsp, req)
}



func main() {

	flag.Parse()

	if flag.NFlag() < 1 {
		flag.Usage()
		return
	}

	fmt.Printf("file server root dir: %s\n", *fileDir)
	muxHandle := newMuxHandle(*fileDir)

	http.Handle("/", muxHandle)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}