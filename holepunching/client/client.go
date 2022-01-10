package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const HANDSHAKEMSG = "我是打洞消息"

var tag string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请输入一个客户端标志")
		os.Exit(0)
	}
	// 当前进程标记字符串，便与显示
	tag = os.Args[1]
	srcAddr := &net.UDPAddr{
		IP: net.IPv4zero,
		Port: 9902, // 注意端口必须固定
	}
	dstAddr := &net.UDPAddr{
		IP: net.ParseIP("192.168.0.100"),
		Port: 9527,
	}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println("connection error:", err)
	}
	if _, err = conn.Write([]byte("hello, I'am new peer:" + tag)); err != nil {
		log.Panic(err)
	}
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Printf("error during read: %s", err)
	}
	conn.Close()
	anotherPeer := parseAddr(string(data[:n]))
	fmt.Printf("local:%s, server: %s, another: %s\n", srcAddr, remoteAddr, anotherPeer)
	// 开始打洞
	bidirectionHole(srcAddr, &anotherPeer)
}

func parseAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP: net.ParseIP(t[0]),
		Port: port,
	}
}

func bidirectionHole(srcAddr, anotherAddr *net.UDPAddr)  {
	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	// 向另一个peer发送一条udp消息(对方peer的nat设备会丢弃该消息,非法来源),
	//用意是在自身的nat设备打开一条可进入的通道,这样对方peer就可以发过来udp消息
	if _, err = conn.Write([]byte(HANDSHAKEMSG)); err != nil {
		log.Println("send handshake msg(这个错误是预期中的错误):", err)
	}
	go func() {
		for {
			time.Sleep(10 * time.Second)
			msg := "from [" + tag + "]"
			if _, err = conn.Write([]byte(msg)); err != nil {
				log.Println("send msg failed", err)
			}
		}
	}()
	for {
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		}else {
			log.Printf("the receved data: %s\n", data[:n])
		}
	}
}
