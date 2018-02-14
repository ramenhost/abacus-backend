package helpers

import (
	"fmt"
	"os"

	qrcode "github.com/skip2/go-qrcode"
)

//GenerateQR generates QR code png file
func GenerateQR(aid string) error {
	_ = os.Mkdir("qr", os.ModePerm)
	err := qrcode.WriteFile(aid, qrcode.Medium, 256, fmt.Sprintf("qr/%s.png", aid))
	return err
}
