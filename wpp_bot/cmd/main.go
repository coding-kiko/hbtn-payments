package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"wpp_bot/bot"

	_ "github.com/lib/pq"
)

func main() {
	GroupId := os.Getenv("GROUP_JID")
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASSWORD")
	backendUri := os.Getenv("BACKEND_URI")
	registerPaymentPath := os.Getenv("REGISTER_PAYMENT_PATH")
	getSummaryPath := os.Getenv("GET_SUMMARY_PATH")

	conn := bot.NewWppClientConn(dbUser, dbPwd)

	registerPaymentUri := backendUri + registerPaymentPath
	getSummaryUri := backendUri + getSummaryPath
	client := bot.NewClient(conn, GroupId, registerPaymentUri, getSummaryUri)

	conn.AddEventHandler(client.EventHandler)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill, syscall.SIGKILL)
	<-c

	conn.Disconnect()
	log.Println("Gracefully shutdown")
}
