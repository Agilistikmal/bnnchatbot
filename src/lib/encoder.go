package lib

import (
	"bytes"
	"encoding/base64"
	"fmt"
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

func GenerateQRBase64(qrText string) (string, error) {
	var buf bytes.Buffer
	qr, _ := qrcode.Encode(qrText, qrcode.Medium, 256)
	_, err := buf.Write(qr)
	if err != nil {
		return "", err
	}

	base64Img := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/png;base64," + base64Img, nil
}
