package s3

import (
	"bytes"
	"dots-api/bootstrap"
	"dots-api/lib/upload"
	"encoding/base64"
	"strconv"
	"strings"
	"time"
)

type contract struct {
	app *bootstrap.App
}

func New(app *bootstrap.App) *contract {
	return &contract{app}
}

func (s *contract) UploadFileS3(paramName, fileMime string, fileSize int64, fileHeader []byte) (string, error) {

	var err error
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	encodingName := base64.StdEncoding.EncodeToString([]byte(timestamp))

	splitName := strings.Split(paramName, ".")
	extName := splitName[len(splitName)-1]

	filename := s.app.Config.GetString("aws.s3.filepath") + "/" + encodingName + "." + extName
	s3Info := new(upload.S3Info)
	s3Info.Key = s.app.Config.GetString("aws.s3.key")
	s3Info.Secret = s.app.Config.GetString("aws.s3.secret")
	s3Info.Region = s.app.Config.GetString("aws.s3.region")
	s3Info.Bucket = s.app.Config.GetString("aws.s3.bucket")
	s3Info.Filename = filename
	s3Info.Filemime = fileMime
	s3Info.Filesize = fileSize

	buffer := bytes.NewReader(fileHeader)

	err = upload.PushS3Buffer(buffer, *s3Info)
	if err != nil {
		return "", err
	}

	filePath := s.app.Config.GetString("aws.s3.public_url") + filename
	return filePath, nil
}
