package client

import (
	kcp "github.com/xtaci/kcp-go"
	"log"
	"os"
	"fmt"
	"bufio"
	"io"
	"time"
)

type Client struct {
	udpSession *kcp.UDPSession
}

func (conn *Client) Connect(raddr string)  {
	udpSession, err := kcp.DialWithOptions(raddr, nil, 0, 0)
	if err != nil {
		log.Panic(err)
	}
	udpSession.SetWindowSize(4096, 4096)
	udpSession.SetReadBuffer(4 * 1024 * 1024)
	udpSession.SetWriteBuffer(4 * 1024 * 1024)
	udpSession.SetStreamMode(true)
	udpSession.SetMtu(1400)
	udpSession.SetACKNoDelay(true)
	conn.udpSession = udpSession
}

func (conn *Client) SendFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	bf := bufio.NewReader(f)
	stat, _ := f.Stat()
	fmt.Printf("文件大小：%dm\n", stat.Size()/1024/1024)
	buffer  := make([]byte, 102400)

	start := time.Now()
	io.CopyBuffer(conn.udpSession, bf, buffer)
	end := time.Now()
	delta := end.Sub(start)
	fmt.Printf("传输用时%fs, %fm/s\n", delta.Seconds(), float64(stat.Size())/delta.Seconds()/1024/1024)
	f.Close()
	conn.udpSession.Close()
}

