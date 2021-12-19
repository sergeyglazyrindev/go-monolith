package core

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type IStorageInterface interface {
	Save(f *FileForStorage) (string, error)
	Read(filename string) ([]byte, error)
	Stats(filename string) (os.FileInfo, error)
	Exists(filename string) (bool, error)
	Delete(filename string) (bool, error)
	GetUploadURL() string
}

type FileForStorage struct {
	Content           []byte
	PatternForTheFile string
	Filename          string
}

type FsStorage struct {
	UploadPath string
	URLPath    string
}

func (s *FsStorage) GetUploadURL() string {
	return s.URLPath
}

func (s *FsStorage) Save(f *FileForStorage) (string, error) {
	tmpfile, err := ioutil.TempFile(s.UploadPath, f.PatternForTheFile)
	if err != nil {
		return "", err
	}
	_, err = tmpfile.Write(f.Content)
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()
	return strings.TrimPrefix(strings.Replace(tmpfile.Name(), s.UploadPath, "", 1), "/"), nil
}

func (s *FsStorage) Read(filename string) ([]byte, error) {
	content, err := os.ReadFile(fmt.Sprintf("%s/%s", s.UploadPath, filename))
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (s *FsStorage) Stats(filename string) (os.FileInfo, error) {
	f, err := os.OpenFile(fmt.Sprintf("%s/%s", s.UploadPath, filename), os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (s *FsStorage) Exists(filename string) (bool, error) {
	filepath := fmt.Sprintf("%s/%s", s.UploadPath, filename)
	var err error
	if _, err = os.Stat(filepath); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

func (s *FsStorage) Delete(filename string) (bool, error) {
	filepath := fmt.Sprintf("%s/%s", s.UploadPath, filename)
	err := os.Remove(filepath)
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewFsStorage() IStorageInterface {
	return &FsStorage{UploadPath: CurrentConfig.GetPathToUploadDirectory(), URLPath: CurrentConfig.GetURLToUploadDirectory()}
}

type AWSS3Storage struct {
	URLPath    string
	Config     *AWSConfig
	Timeout    time.Duration
	Bucket     string
	NameLength int
	Domain     string
}

func (s *AWSS3Storage) GetUploadURL() string {
	return fmt.Sprintf("%s/%s", s.Domain, s.URLPath)
}

func (s *AWSS3Storage) Save(f *FileForStorage) (string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s.Config.S3.Region),
		Credentials: credentials.NewStaticCredentials(
			s.Config.S3.AccessKey,
			s.Config.S3.SecretKey,
			"",
		),
	}))
	svc := s3.New(sess)
	ctx := context.Background()
	var cancelFn func()
	if s.Timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, s.Timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}
	s3Key := GenerateRandomString(s.NameLength, &OnlyLetersNumbersStringAlphabet)
	ext := ""
	filenameParts := strings.Split(f.Filename, ".")
	if len(filenameParts) > 1 {
		ext = "." + filenameParts[1]
	}
	s3Key += ext
	_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s.URLPath + "/" + s3Key),
		Body:   bytes.NewReader(f.Content),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			// If the SDK can determine the request or retry delay was canceled
			// by a context the CanceledErrorCode error code will be returned.
			return "", fmt.Errorf("upload canceled due to timeout, %v", err)
		}
		return "", fmt.Errorf("failed to upload object, %v", err)
	}
	return s3Key, nil
}

func (s *AWSS3Storage) Read(s3Key string) ([]byte, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s.Config.S3.Region),
		Credentials: credentials.NewStaticCredentials(
			s.Config.S3.AccessKey,
			s.Config.S3.SecretKey,
			"",
		),
	}))
	svc := s3.New(sess)
	ctx := context.Background()
	var cancelFn func()
	if s.Timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, s.Timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}
	res, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return nil, err
	}
	b := make([]byte, *res.ContentLength)
	_, err = res.Body.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s *AWSS3Storage) Stats(s3Key string) (os.FileInfo, error) {
	return nil, nil
}

func (s *AWSS3Storage) Exists(s3Key string) (bool, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s.Config.S3.Region),
		Credentials: credentials.NewStaticCredentials(
			s.Config.S3.AccessKey,
			s.Config.S3.SecretKey,
			"",
		),
	}))
	svc := s3.New(sess)
	ctx := context.Background()
	var cancelFn func()
	if s.Timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, s.Timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}
	res, err := svc.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return false, err
	}
	return *res.ContentLength > 0, nil
}

func (s *AWSS3Storage) Delete(s3Key string) (bool, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s.Config.S3.Region),
		Credentials: credentials.NewStaticCredentials(
			s.Config.S3.AccessKey,
			s.Config.S3.SecretKey,
			"",
		),
	}))
	svc := s3.New(sess)
	ctx := context.Background()
	var cancelFn func()
	if s.Timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, s.Timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}
	res, err := svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return false, err
	}
	return *res.DeleteMarker, nil
}

func NewAWSS3Storage(uploadPath string, s3Config *AWSConfig) IStorageInterface {
	return &AWSS3Storage{
		URLPath:    uploadPath,
		Config:     s3Config,
		NameLength: 20,
	}
}

type AWSS3Config struct {
	Region    string `yaml:"region"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
}

type AWSConfig struct {
	S3 *AWSS3Config `yaml:"s3"`
}

func NewAWSConfig() *AWSConfig {
	config := &AWSConfig{}
	err := yaml.Unmarshal(CurrentConfig.ConfigContent, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return config
}
