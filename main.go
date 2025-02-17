package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/agilistikmal/bnnchat/src/config"
	"github.com/agilistikmal/bnnchat/src/database"
	"github.com/agilistikmal/bnnchat/src/handlers"
	"github.com/agilistikmal/bnnchat/src/services"
	_ "github.com/lib/pq"
	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/viper"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	config.NewConfig()
	db := database.NewDatabase()

	questionService := services.NewQuestionService(db)

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("postgres", viper.GetString("postgres.dsn"), dbLog)
	if err != nil {
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	client := whatsmeow.NewClient(deviceStore, nil)

	h := handlers.NewHandler(client, questionService)

	client.AddEventHandler(h.MessageEvent)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
