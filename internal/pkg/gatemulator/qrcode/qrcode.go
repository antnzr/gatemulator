package qrcode

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"

	qrcode "github.com/skip2/go-qrcode"
)

const qrCodeSize = 256

func Generate(content string) (string, error) {
	qrCode, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return "", fmt.Errorf("could not generate a QR code: %v", err)
	}
	fmt.Printf("%v", qrCode.Content)

	var pngBuffer bytes.Buffer
	err = png.Encode(&pngBuffer, qrCode.Image(qrCodeSize))
	if err != nil {
		return "", err
	}

	qrBase64 := base64.StdEncoding.EncodeToString(pngBuffer.Bytes())
	qrDataURL := "data:image/png;base64," + qrBase64

	return qrDataURL, nil
}
