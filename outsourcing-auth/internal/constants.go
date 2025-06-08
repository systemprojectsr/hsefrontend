package internal

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

var PostgresUser string
var PostgresPassword string
var PostgresDB string
var PostgresHost string
var PostgresPort string
var KeyJWT string

var LifeTimeJWT int

func InitEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PW")
	PostgresDB = os.Getenv("POSTGRES_DB")
	PostgresHost = os.Getenv("POSTGRES_HOST")
	PostgresPort = os.Getenv("POSTGRES_PORT")
	KeyJWT = os.Getenv("KEY_JWT")
	lifeTime, err := strconv.ParseInt(os.Getenv("LIFE_TIME_JWT"), 10, 64)
	if err != nil {
		return err
	}
	LifeTimeJWT = int(lifeTime)
	return nil
}

type StatusResponse struct {
	Status string `json:"status"`
}

type InfoResponse struct {
	*StatusResponse
	Description string `json:"description"`
}
