package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/agilistikmal/bnnchat/src/config"
	"github.com/agilistikmal/bnnchat/src/database"
	"github.com/agilistikmal/bnnchat/src/handlers"
	"github.com/agilistikmal/bnnchat/src/services"
	"github.com/agilistikmal/bnnchat/src/web/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	log "github.com/sirupsen/logrus"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	log.Info("Starting...")

	log.Info("Loading config and database...")
	config.NewConfig()
	db := database.NewDatabase()

	log.Info("Loading services...")
	menuService := services.NewMenuService(db)
	log.Info(menuService)

	// ---
	// Website
	// ---
	go func() {
		log.Info("Preparing web server...")
		views := html.New("./src/web/views", ".html")

		app := fiber.New(fiber.Config{
			Views: views,
		})

		app.Static("/assets", "./assets")

		dashboardController := controllers.NewDashboardController(nil, menuService)
		menuController := controllers.NewMenuController(menuService)

		app.Get("/", dashboardController.Dashboard)
		app.Get("/menu_part", dashboardController.MenuPart)
		app.All("/menu/add", menuController.Add)
		app.All("/menu/:id", menuController.Detail)
		app.All("/menu/:menuID/submenu", menuController.SubMenu)

		app.Listen(":3000")
	}()

	// ---
	// WhatsApp
	// ---
	var client *whatsmeow.Client
	go func() {
		log.Info("Preparing whatsapp client...")
		dbLog := waLog.Stdout("Database", "DEBUG", true)
		container, err := sqlstore.New("sqlite3", "file:data.db?_foreign_keys=on", dbLog)
		if err != nil {
			panic(err)
		}
		deviceStore, err := container.GetFirstDevice()
		if err != nil {
			panic(err)
		}
		client = whatsmeow.NewClient(deviceStore, nil)

		h := handlers.NewHandler(client, menuService)

		client.AddEventHandler(h.MessageEvent)

		if client.Store.ID == nil {
			log.Info("New Session Created")
			qrChan, _ := client.GetQRChannel(context.Background())
			err = client.Connect()
			if err != nil {
				panic(err)
			}
			for evt := range qrChan {
				if evt.Event == "code" {
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
					log.Info("Scan QR Code on Link Device Whatsapp")
				} else {
					log.Info("Login event ::", evt.Event)
				}
			}
		} else {
			err = client.Connect()
			if err != nil {
				log.Fatal("Login from session error ::", err.Error())
			}
			log.Info("Login from session")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
