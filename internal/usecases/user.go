package usecases

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/akshaybt001/DatingApp_UserService/internal/adapters"
	"github.com/akshaybt001/DatingApp_proto_files/pb"
	"github.com/go-redis/redis"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type UserUseCase struct {
	userAdapter adapters.AdapterInterface
}

func NewUserUseCase(useradapter adapters.AdapterInterface) *UserUseCase {
	return &UserUseCase{
		userAdapter: useradapter,
	}
}
var redisClient *redis.Client

func init() {
	fmt.Println("hii from init redis")
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func (user *UserUseCase) UploadImage(req *pb.UserImageRequest, profileId string) (string, error) {
	minioClient, err := minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESSKEY"),os.Getenv("MINIO_SECRETKEY"),""),
		Secure: false,
	})
	if err != nil {
		log.Print("error while initialising minio", err)
		return "", err
	}
	objectName := "images/" + req.ObjectName
	contentType := `image/jpeg`
	n, err := minioClient.PutObject(context.Background(), os.Getenv("BUCKET_NAME"), objectName, bytes.NewReader(req.ImageData), int64(len(req.ImageData)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Println("error while uploading to minio", err)
		return "", err
	}
	log.Printf("Successfully uploaded %s of size %v\n", objectName, n)
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), os.Getenv("BUCKET_NAME"), objectName, time.Second*24*60*60, nil)
	if err != nil {
		log.Println("error while generating presigned URL", err)
		return "", err
	}
	url, err := user.userAdapter.UploadProfileImage(presignedURL.String(), profileId)
	return url, err
}


