package com

import (
	"bufio"
	"log"
	"syscall"
	"time"

	"github.com/schleibinger/sio"
)

var port *sio.Port
var callback func([]byte)

func Init(f func([]byte)) {
	// устанавливаем соединение
	porter, err := sio.Open("/dev/ttyS10", syscall.B9600)
	if err != nil {
		log.Fatal(err)
	}
	port = porter
	callback = f

	// отправляем данные
	_, err = port.Write([]byte("This test string is sended to COM and received back!\n"))
	if err != nil {
		log.Fatal(err)
	}

	go comRecv()
}

func Send(data []byte) {
	var err error
	// отправляем данные
	_, err = port.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func comSend() {
	var err error
	for i := 0; ; i++ {
		time.Sleep(time.Second)

		// отправляем данные
		_, err = port.Write([]byte("test\n"))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func comRecv() {
	reader := bufio.NewReader(port)
	for i := 0; ; i++ {
		//time.Sleep(time.Second)
		// получаем данные
		reply, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("recieved: %q", reply)
		callback(reply)
	}
}
