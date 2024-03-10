package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Vit0Corleone/vpn/command"

	"github.com/songgao/water"
)

var conn net.Conn
var serverIP string

// to start set ip addr of vpn server
func main() {
	iface, err := createTun("192.168.0.10")
	if err != nil {
		fmt.Println("interface can not created:", err)
		return
	}
	conn, err = createConn(serverIP)
	if err != nil {
		fmt.Println("tcp conn create error:", err)
	}

	go listen(iface)
	go listenInterface(iface)

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, os.Interrupt, syscall.SIGTERM)
	<-termSignal
	fmt.Println("closing")
}

func createConn(ip string) (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("%s:44433", ip)) // vpn server addr
}

func listen(iface *water.Interface) {
	for {
		fmt.Println("tcp connection listening")
		message := make([]byte, 65535)
		for {
			n, err := conn.Read(message)
			if err != nil {
				log.Fatal("conn read error:", err)
			}

			fmt.Printf("Read: %s\n", message[:n])

			if iface != nil {
				_, err = iface.Write(message[:n])
				if err != nil {
					log.Fatal("ifce write err:", err)
				} else {
					fmt.Println("iface write done")
				}
			}
		}
	}
}

func listenInterface(iface *water.Interface) {
	fmt.Println("interface listening")
	packet := make([]byte, 65535)
	for {
		n, err := iface.Read(packet)
		if err != nil {
			log.Fatal("ifce read error:", err)
		}

		if err == nil {
			_, err = conn.Write(packet[:n])
			if err != nil {
				log.Fatal("conn write error:", err)
			}
			fmt.Println("conn write done")
		}
	}
}

func createTun(ip string) (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}

	iface, err := water.New(config)
	if err != nil {
		return nil, err
	}
	log.Printf("Interface Name: %s\n", iface.Name())
	out, err := command.RunCommand(fmt.Sprintf("ip addr add %s/24 dev %s", ip, iface.Name()))
	if err != nil {
		fmt.Println(out)
		return nil, err
	}

	out, err = command.RunCommand(fmt.Sprintf("ip link set dev %s up", iface.Name()))
	if err != nil {
		fmt.Println(out)
		return nil, err
	}
	return iface, nil
}
