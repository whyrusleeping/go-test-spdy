package main

import (
	"flag"
	"fmt"
	gmux "github.com/jbenet/go-stream-muxer"
	spdy "github.com/whyrusleeping/go-smux-spdystream"
	"io/ioutil"
	"log"
	"net"
)

func main() {
	listenport := flag.Int("l", 0, "specify listen port")
	flag.Parse()
	if *listenport != 0 {
		log.Println("listening")
		list, err := net.Listen("tcp", fmt.Sprintf(":%d", *listenport))
		if err != nil {
			log.Fatalln(err)
		}
		defer list.Close()

		con, err := list.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		defer con.Close()

		sc, err := spdy.Transport.NewConn(con, true)
		if err != nil {
			log.Fatalln(err)
		}

		sc.Serve(func(s gmux.Stream) {
			log.Println("Got a stream!")
			out, err := ioutil.ReadAll(s)
			if err != nil {
				log.Println("error reading: ", err)
				s.Close()
				return
			}
			log.Printf("read: %q\n", out)
			s.Close()
			log.Println("closed stream")
		})
	} else {
		log.Println("dialing!")
		con, err := net.Dial("tcp", flag.Args()[0])
		if err != nil {
			log.Fatalln(err)
		}
		defer con.Close()

		cc, err := spdy.Transport.NewConn(con, false)
		if err != nil {
			log.Fatalln(err)
		}
		defer cc.Close()
		go cc.Serve(func(s gmux.Stream) {
			log.Println("client got a stream?")
		})

		log.Println("creating stream")
		s, err := cc.OpenStream()
		if err != nil {
			log.Fatalln(err)
		}
		defer s.Close()

		log.Println("reading stream")
		s.Write([]byte("hello everyone this is a test"))
	}
}

var _ = ioutil.NopCloser
