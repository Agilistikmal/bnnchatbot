package lib

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/skip2/go-qrcode"
)

func EncodeBase62(num int) string {
	c := fmt.Sprint(num)

	for len(c) < 3 {
		c = "0" + c
	}
	result := base64.StdEncoding.EncodeToString([]byte(c))

	return result
}

func DecodeBase62(encoded string) int {
	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	encoded = strings.TrimLeft(string(decoded), "0")

	result, _ := strconv.Atoi(encoded)

	return result
}

func GenerateQRToFile(qrText, filePath string) error {
	qr, err := qrcode.Encode(qrText, qrcode.Medium, 256)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, qr, 0644)
	if err != nil {
		return err
	}

	return nil
}
