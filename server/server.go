package server

import (
	"github.com/xtaci/kcp-go"
	"log"
	"os"
	"io"
)

type Server struct {
	Died chan struct{}
}

func (s *Server) Listen(laddr string) {
	l, err := kcp.ListenWithOptions(laddr, nil, 0, 0)
	if err != nil {
		log.Panic(err)
	}
	for {
		udpSession, err := l.AcceptKCP()
		if err != nil {
			log.Panic(err)
		}
		udpSession.SetWindowSize(4096, 4096)
		udpSession.SetReadBuffer(4 * 1024 * 1024)
		udpSession.SetWriteBuffer(4 * 1024 * 1024)
		udpSession.SetStreamMode(true)
		udpSession.SetMtu(1400)
		udpSession.SetACKNoDelay(true)
		go func(udpSession *kcp.UDPSession, died chan struct{}) {
			f, err := os.Create("tmp.data")
			if err != nil {
				log.Panic(err)
			}
			buffer := make([]byte, 102400)
			io.CopyBuffer(f, udpSession, buffer)
			f.Close()
			udpSession.Close()
			died <- struct{}{}
		}(udpSession, s.Died)
	}
}
