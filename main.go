package main

import (
	"flag"
	"github.com/withliyh/filetransfer/server"
	"github.com/withliyh/filetransfer/client"
)

var raddr = flag.String("raddr", "127.0.0.1:10000", "指定文件接收服务器地址")
var file = flag.String("file", "", "指定传输的文件路径")
var d = flag.Bool("d", false, "作为服务端接收文件")



func main() {

	flag.Parse()

	if flag.NFlag() < 1 {
		flag.Usage()
		return
	}

	if *d {
		s := server.Server{}
		s.Listen("0.0.0.0:10000")
		select {
		case <- s.Died:
		}
	} else {
		c := client.Client{}
		c.Connect(*raddr)
		c.SendFile(*file)
	}
}