package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

type GetOsVariables struct {
	EndPoint          string
	MinioAccessKey    string
	MinioAccessSecret string
	MinioBucket       string
	SecondServiceUrl  string
	MinioSSL          bool
}

var once sync.Once

var instance *GetOsVariables

func NewVariableService() *GetOsVariables {
	once.Do(func() {
		var isSSL bool
		isSSL, err := strconv.ParseBool(os.Getenv("MINIO_SSL"))
		if err != nil {
			isSSL = false
		}
		instance = &GetOsVariables{
			EndPoint:          fmt.Sprintf("%s:%s", os.Getenv("MINIO_ENDPOINT"), os.Getenv("MINIO_PORT")),
			MinioAccessKey:    os.Getenv("MINIO_ACCESS_KEY"),
			MinioAccessSecret: os.Getenv("MINIO_SECRET_KEY"),
			MinioBucket:       os.Getenv("MINIO_BUCKET_NAME"),
			MinioSSL:          isSSL,
			SecondServiceUrl:  os.Getenv("SECOND_SERVICE"),
		}
	})

	return instance
}
