package lib

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func WaitForServer(url string) {
	for {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			logrus.Info("Server is ready:", url)
			resp.Body.Close()
			break
		}
		logrus.Info("Waiting for server to be ready...")
		time.Sleep(500 * time.Millisecond) // Tunggu sebentar sebelum cek lagi
	}
}
