package lib

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"shin/src/config"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client
var mimeTypeToExt = map[string]string{
	"image/png":          "png",
	"image/jpeg":         "jpg",
	"application/pdf":    "pdf",
	"application/msword": "doc",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": "docx",
}

func hashFile(file io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate MD5: %v", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func Upload(file multipart.File, fileName string) (string, string) {

	fileContent, _ := io.ReadAll(file)

	hash, _ := hashFile(bytes.NewReader(fileContent))

	// Read the first 512 bytes of the file & Detect the file's content type
	buf := fileContent[:512]
	mimeType := http.DetectContentType(buf)

	filename := fmt.Sprintf("%s.%s", hash, mimeTypeToExt[mimeType])

	if _, err := S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(config.Config.S3.Bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(fileContent),
	}); err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf("%s/%s", config.Config.S3.CDNUrl, filename), fileName
}

// Initializing
func InitS3Lib() {
	S3Client = s3.New(s3.Options{
		Credentials:      credentials.NewStaticCredentialsProvider(config.Config.S3.AccessKeyId, config.Config.S3.SecretAccessKey, ""),
		Region:           config.Config.S3.DefaultRegion,
		RetryMaxAttempts: 5,
		RetryMode:        aws.RetryModeStandard,
		HTTPClient:       &http.Client{Timeout: 30 * time.Second},
		ClientLogMode:    aws.LogRequestWithBody | aws.LogResponseWithBody,
	})
}
