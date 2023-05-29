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

	if err = http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
	fmt.Println("servicio iniciado en puerto 3000")

	/* app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {

		//call second service
		str, err := getIdFromService()
		if err != nil {
			c.Response().SetStatusCode(fiber.StatusInternalServerError)
			return c.JSON("Error al consultar segundo servicio")
		}
		//call minio
		info, err := uploadToMinio(str)
		if err != nil {
			c.Response().SetStatusCode(fiber.StatusInternalServerError)
			return c.JSON(err)
		}

		c.Response().SetStatusCode(fiber.StatusOK)
		return c.JSON(map[string]interface{}{
			"message": info,
			"status":  "Ok",
		})
	})

	app.Listen(":3000") */
}

func getIfFromServiceASync(service GetOsVariables) <-chan ResponseService {
	url := service.SecondServiceUrl
	res := make(chan ResponseService)
	go func() {
		var result ResponseService
		//	response, err := http.Get(os.Getenv("SECOND_SERVICE"))
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

/* func getIdFromService() (string, error) {

	response, err := http.Get(os.Getenv("SECOND_SERVICE"))

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil

} */

func uploadToMinio(str string, service GetOsVariables) (info minio.UploadInfo, err error) {

	//url := fmt.Sprintf("%s:%s", os.Getenv("MINIO_ENDPOINT"), os.Getenv("MINIO_PORT"))

	/* 	minioClient, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), ""),
		Secure: false,
	}) */

	minioClient, err := minio.New(service.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(service.MinioAccessKey, service.MinioAccessSecret, ""),
		Secure: false,
	})

	if err != nil {
		panic(err)
	}

	name := fmt.Sprintf("golangFile%s.txt", str)
	rd := strings.NewReader(str)

	//return minioClient.PutObject(context.Background(), os.Getenv("MINIO_BUCKET_NAME"), name, rd, rd.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	return minioClient.PutObject(context.Background(), service.MinioBucket, name, rd, rd.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})

}
