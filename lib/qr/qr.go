package qr

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func GenerateQRCode(data, fileName, filePath string) (string, error) {
	if fileName == "" {
		return "", errors.New("qrcode file name can not be empty")
	}

	qrc, err := qrcode.New(data)
	if err != nil {
		return "", err
	}

	fileName = fmt.Sprintf("%s.jpg", fileName)
	w, err := standard.New(
		fmt.Sprintf("%s/%s", filePath, fileName), standard.WithQRWidth(32),
		standard.WithBgColor(color.White),
		standard.WithFgColor(color.Black),
	)
	if err != nil {
		return "", err
	}

	defer w.Close()
	if err := qrc.Save(w); err != nil {
		return "", err
	}

	return fileName, nil
}
