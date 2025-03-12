package main

import (
	"context"
	"html/template"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agilistikmal/bnnchat/src/config"
	"github.com/agilistikmal/bnnchat/src/database"
	"github.com/agilistikmal/bnnchat/src/handlers"
	"github.com/agilistikmal/bnnchat/src/lib"
	"github.com/agilistikmal/bnnchat/src/services"
	"github.com/agilistikmal/bnnchat/src/web/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	log "github.com/sirupsen/logrus"
	webview "github.com/webview/webview_go"
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
					lib.GenerateQRToFile(evt.Code, "./assets/qr.png")
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
					log.Info("Scan QR Code on Link Device Whatsapp")
				} else {
					log.Info("Login event ", evt.Event)
					os.Remove("./assets/qr.png")
				}
			}
		} else {
			err = client.Connect()
			if err != nil {
				log.Fatal("Login from session error ::", err.Error())
			}
			log.Info("Login from session")
			os.Remove("./assets/qr.png")
		}
	}()

	// ---
	// Website
	// ---
	go func() {
		log.Info("Preparing web server...")
		views := html.New("./views", ".html")
		views.AddFunc("safeHTML", func(s string) template.HTML {
			return template.HTML(s)
		})

		app := fiber.New(fiber.Config{
			Views: views,
		})

		app.Static("/assets", "./assets")

		dashboardController := controllers.NewDashboardController(client, menuService)
		menuController := controllers.NewMenuController(menuService)

		app.Get("/", dashboardController.Dashboard)
		app.Get("/menu_part", dashboardController.MenuPart)
		app.All("/menu/add", menuController.Add)
		app.All("/menu/:id", menuController.Detail)
		app.All("/menu/:menuID/submenu", menuController.SubMenu)
		app.All("/menu/:menuID/submenu/position", menuController.SubMenuPosition)

		app.Post("/logout", dashboardController.Logout)

		app.Listen(":3000")
	}()

	c := make(chan os.Signal, 1)

	go func(exitChan chan os.Signal) {
		for client == nil {
			log.Info("Waiting for WhatsApp client to be ready...")
			time.Sleep(1 * time.Second)
		}
		w := webview.New(true)
		defer w.Destroy()
		w.SetTitle("BNN Chatbot Dashboard")
		w.SetSize(1024, 768, webview.HintNone)
		w.Navigate("http://localhost:3000")

		go func() {
			<-exitChan
			w.Terminate()
		}()

		w.Run()
		exitChan <- os.Interrupt
	}(c)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
