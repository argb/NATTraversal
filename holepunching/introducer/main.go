package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	udpAddr := &net.UDPAddr{
		IP: net.IPv4zero, // 0地址，相当于监听所有本地ip地址。 https://blog.csdn.net/liyi1009365545/article/details/84780476
		Port: 9527,
	}
	listener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("local address: %s \n", listener.LocalAddr().String())
	peers := make([]net.UDPAddr, 0, 2)
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("read error: %s", err)
		}
		log.Printf("<%s>, data: %s\n", remoteAddr.String(), data[:n])
		peers = append(peers, *remoteAddr)
		if len(peers) == 2 {
			log.Printf("punching hole with udp, establish connection: %s <---> %s\n", peers[0].String(), peers[1].String())
			listener.WriteToUDP([]byte(peers[0].String()), &peers[0])
			listener.WriteToUDP([]byte(peers[0].String()), &peers[1])
			time.Sleep(time.Second * 8)
			log.Println("完成介绍，介绍人服务器退出(介绍人服务退出后不影响建立连接的两个peer)")
			return
		}
	}
}
