package upload

import (
	"bytes"
	"encoding/base64"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Info file information
type S3Info struct {
	Key      string
	Secret   string
	Region   string
	Bucket   string
	Filename string
	Filemime string
	Filesize int64
}

// PushS3Buffer ...
func PushS3Buffer(buffer *bytes.Reader, in S3Info) error {
	session, err := session.NewSession(&aws.Config{
		Region:      &in.Region,
		Credentials: credentials.NewStaticCredentials(in.Key, in.Secret, ""),
	})
	if err != nil {
		return err
	}

	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(in.Bucket),
		Key:                  aws.String(in.Filename),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 buffer,
		ContentLength:        aws.Int64(in.Filesize),
		ContentType:          aws.String(in.Filemime),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return err
	}

	return nil
}

// PushS3ByPath ...
func PushS3ByPath(fPath string, in S3Info) error {
	session, err := session.NewSession(&aws.Config{
		Region:      &in.Region,
		Credentials: credentials.NewStaticCredentials(in.Key, in.Secret, ""),
	})
	if err != nil {
		return err
	}

	file, err := os.Open(fPath)
	if err != nil {
		log.Printf("os.Open - filename: %s, err: %v", fPath, err)
		return err
	}
	defer file.Close()

	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(in.Bucket),
		Key:                  aws.String(in.Filename),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 file,
		ContentLength:        aws.Int64(in.Filesize),
		ContentType:          aws.String(in.Filemime),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return err
	}

	return nil
}

func PushS3Base64(in S3Info, base64File string) error {
	b64data := base64File[strings.IndexByte(base64File, ',')+1:]
	decode, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return err
	}

	session, err := session.NewSession(&aws.Config{
		Region:      &in.Region,
		Credentials: credentials.NewStaticCredentials(in.Key, in.Secret, ""),
	})
	if err != nil {
		return err
	}

	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(in.Bucket),
		Key:                  aws.String(in.Filename),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 bytes.NewReader(decode),
		ContentType:          aws.String(in.Filemime),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		// ContentLength:        aws.Int64(in.Filesize),
	})
	if err != nil {
		return err
	}

	return err
}
