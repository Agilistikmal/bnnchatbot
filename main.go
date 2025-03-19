package main

import (
	"context"
	"embed"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

//go:embed assets/*
var assetsFS embed.FS

//go:embed views/**/*
var viewsFS embed.FS

func main() {
	log.Info("Starting...")

	os.Mkdir("public", 0755)

	log.Info("Loading database...")
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

		h := handlers.NewHandler(client, db, menuService)

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
					lib.GenerateQRToFile(evt.Code, "./public/qr.png")
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
					log.Info("Scan QR Code on Link Device Whatsapp")
				} else {
					log.Info("Login event ", evt.Event)
					os.Remove("./public/qr.png")
				}
			}
		} else {
			err = client.Connect()
			if err != nil {
				log.Fatal("Login from session error ::", err.Error())
			}
			log.Info("Login from session")
			os.Remove("./public/qr.png")
		}
	}()

	// ---
	// Website
	// ---
	go func() {
		log.Info("Preparing web server...")

		for client == nil {
			log.Info("Waiting for WhatsApp client to be ready...")
			time.Sleep(1 * time.Second)
		}

		subFS, err := fs.Sub(viewsFS, "views")
		if err != nil {
			log.Fatal("Gagal membuat sub-FS:", err)
		}
		views := html.NewFileSystem(http.FS(subFS), ".html")

		app := fiber.New(fiber.Config{
			Views: views,
		})

		app.Static("/public", "./public")

		app.Get("/assets/*", func(c *fiber.Ctx) error {
			filePath := c.Params("*")
			file, err := assetsFS.Open("assets/" + filePath)
			if err != nil {
				return c.SendStatus(fiber.StatusNotFound)
			}
			defer file.Close()

			return c.SendFile("assets/" + filePath)
		})

		dashboardController := controllers.NewDashboardController(client, menuService)
		menuController := controllers.NewMenuController(menuService)
		helpController := controllers.NewHelpController(db)

		app.Get("/", dashboardController.Dashboard)
		app.Get("/menu_part", dashboardController.MenuPart)
		app.All("/menu/add", menuController.Add)
		app.All("/menu/:id", menuController.Detail)
		app.All("/menu/:menuID/submenu", menuController.SubMenu)
		app.All("/menu/:menuID/submenu/position", menuController.SubMenuPosition)

		app.All("/help_part", helpController.HelpPart)
		app.Delete("/help/:jid", helpController.Delete)

		app.Get("/qrcode", dashboardController.QrCode)
		app.Post("/logout", dashboardController.Logout)

		app.Listen(":3000")
	}()

	c := make(chan os.Signal, 1)

	//
	// Webview
	//
	go func(exitChan chan os.Signal) {
		for client == nil {
			log.Info("Waiting for WhatsApp client to be ready...")
			time.Sleep(1 * time.Second)
		}

		lib.WaitForServer("http://localhost:3000")

		w := webview.New(false)
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
