package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ResponseService struct {
	Id    string
	Error error
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env not loaded")
	}
	service := NewVariableService()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-type", "application/json")
		//call second service
		//str, err := getIdFromService()
		fromService := <-getIfFromServiceASync(*service)

		if fromService.Error != nil {
			w.WriteHeader(500)
			response := map[string]interface{}{
				"status":  500,
				"message": err,
			}
			json, _ := json.Marshal(response)
			w.Write(json)
			return
		}
		//call minio
		info, err := uploadToMinio(fromService.Id, *service)

		if err != nil {
			w.WriteHeader(500)
			response := map[string]interface{}{
				"status":  500,
				"message": "error while upload file to MinIo",
				"stack":   err,
			}
			json, _ := json.Marshal(response)
			w.Write(json)
			return
		}
		response := map[string]interface{}{
			"etag":      info.ETag,
			"versionId": info.VersionID,
		}
		json, _ := json.Marshal(response)
		w.WriteHeader(200)
		w.Write(json)
	})
	fmt.Println("servicio iniciado en puerto 3000")

	if err = http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}

}

func getIfFromServiceASync(service GetOsVariables) <-chan ResponseService {
	url := service.SecondServiceUrl
	res := make(chan ResponseService)
	go func() {
		var result ResponseService
		response, err := http.Get(url)
		if err != nil {
			result.Id = ""
			result.Error = err
			res <- result
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		result.Id = string(body)
		result.Error = err
		res <- result

	}()
	return res
}

func uploadToMinio(str string, service GetOsVariables) (info minio.UploadInfo, err error) {

	//url := fmt.Sprintf("%s:%s", os.Getenv("MINIO_ENDPOINT"), os.Getenv("MINIO_PORT"))

	/* 	minioClient, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), ""),
		Secure: false,
	}) */

	minioClient, err := minio.New(service.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(service.MinioAccessKey, service.MinioAccessSecret, ""),
		Secure: service.MinioSSL,
	})

	if err != nil {
		panic(err)
	}

	name := fmt.Sprintf("golangFile%s.txt", str)
	rd := strings.NewReader(str)

	//return minioClient.PutObject(context.Background(), os.Getenv("MINIO_BUCKET_NAME"), name, rd, rd.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	return minioClient.PutObject(context.Background(), service.MinioBucket, name, rd, rd.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})

}
