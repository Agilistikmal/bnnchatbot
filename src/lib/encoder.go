package lib

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
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
