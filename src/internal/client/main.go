package main

import (
	"fmt"
	"internal/client/config"
	"internal/client/message_receiver_service"
	"internal/client/message_sender_service"
	"io"
	"log"
	"os"
)

func main() {
	fmt.Println("Queues timeout client starts...")
	initLog()
	err := message_sender_service.InsertMessages()
	if err != nil {
		log.Fatal(err)
	}
	err = message_receiver_service.ReceiveMessages()
	if err != nil {
		log.Fatal(err)
	}
}

func initLog() {
	fmt.Println("Start initializing the log")
	logFile, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Failed to create log file")
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
