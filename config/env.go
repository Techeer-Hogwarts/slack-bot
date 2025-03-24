package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv 함수는 .env (filename) 파일을 읽어서 환경변수를 설정합니다.
func LoadEnvFile(filename string) {
	err := godotenv.Load(filename)
	if err != nil {
		log.Println(".env 파일을 읽는데 실패했습니다.")
	}
}

// getEnvVarAsString 함수는 key에 해당하는 환경변수를 읽어서 string으로 반환합니다. 만약 key에 해당하는 환경변수가 없으면 defaultString을 반환합니다.
func GetEnvVarAsString(key string, defaultString string) string {
	value, found := os.LookupEnv(key)
	if !found {
		log.Printf("%s키가 존재하지 않습니다.", key)
		return defaultString
	}
	log.Printf("키 %s가 존재합니다.", key)
	return value
}

// getEnvVarAsInt 함수는 key에 해당하는 환경변수를 읽어서 int로 변환한 후 반환합니다. 만약 key에 해당하는 환경변수가 없거나 변환에 실패하면 defaultInt를 반환합니다.
func GetEnvVarAsInt(key string, defaultInt int) int {
	value, found := os.LookupEnv(key)
	if !found {
		log.Println("키가 존재하지 않습니다.")
		return defaultInt
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Println("키를 int로 변환할 수 없습니다.")
		return defaultInt
	}
	log.Printf("키 %s가 존재합니다.", key)
	return intValue
}

// getEnvVarAsBool 함수는 key에 해당하는 환경변수를 읽어서 bool로 변환한 후 반환합니다. 만약 key에 해당하는 환경변수가 없거나 변환에 실패하면 defaultBool을 반환합니다.
func GetEnvVarAsBool(key string, defaultBool bool) bool {
	value, found := os.LookupEnv(key)
	if !found {
		log.Println("키가 존재하지 않습니다.")
		return defaultBool
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Println("키를 bool로 변환할 수 없습니다.")
		return defaultBool
	}
	log.Printf("키 %s가 존재합니다.", key)
	return boolValue
}
