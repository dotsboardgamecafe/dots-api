package handler

import (
	"dots-api/lib/s3"
	"dots-api/lib/upload"
	"net/http"
)

func (h *Contract) UploadFileAct(w http.ResponseWriter, r *http.Request) {
	var (
		resFileName string
		info        = new(upload.Info)
		name        = "upload"
		allowedExt  = []string{"jpg", "png", "jpeg", "pdf", "webp"}

		// Import contract s3
		uploadContract = s3.New(h.App)
	)
	info.MaxSize = 10
	file, fileInfo, err := info.MultipartHandler(w, r, name, allowedExt)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}
	defer file.Close()

	fileHeader := make([]byte, fileInfo.FileSize)
	if _, err := file.Read(fileHeader); err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}
	resFileName, err = uploadContract.UploadFileS3(fileInfo.Filename, fileInfo.FileMime, fileInfo.FileSize, fileHeader)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, resFileName, nil)
}
