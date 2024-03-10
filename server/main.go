package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Vit0Corleone/vpn/command"

	"github.com/songgao/water"
)

var connection net.Conn

func main() {
	iface, err := createTun()
	if err != nil {
		log.Fatal(err)
	}

	l, err := createListener()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(l.Addr())

	go listen(l, iface)
	listenInterface(l, iface)
}

func createTun() (*water.Interface, error) {
	ip := "192.168.0.141"

	config := water.Config{
		DeviceType: water.TUN,
	}

	iface, err := water.New(config)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Interface name: %s\n", iface.Name())

	out, err := command.RunCommand(
		fmt.Sprintf("ip addr add %s/24 dev %s", ip, iface.Name()),
	)
	if err != nil {
		fmt.Println(out)
		return nil, err
	}

	out, err = command.RunCommand(
		fmt.Sprintf("ip link set dev %s up", iface.Name()),
	)
	if err != nil {
		fmt.Println(out)
		return nil, err
	}

	return iface, nil
}

func createListener() (net.Listener, error) {
	return net.Listen("tcp", ":44433")
}

func listen(l net.Listener, iface *water.Interface) {
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	connection = conn

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

func listenInterface(l net.Listener, iface *water.Interface) {
	fmt.Println("interface listening")
	packet := make([]byte, 65535)
	for {
		n, err := iface.Read(packet)
		if err != nil {
			log.Fatal("ifce read error:", err)
		}

		if connection != nil {
			_, err = connection.Write(packet[:n])
			if err != nil {
				log.Fatal("conn write error:", err)
			}
			fmt.Println("conn write done")
		}
	}
}
