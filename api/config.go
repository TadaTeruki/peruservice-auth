package api

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type envConfig struct {
	Mode             string
	AuthPort         string
	AuthAllowOrigins []string
	PrivateKeyFile   string
	PublicKeyFile    string
	DBHost           string
	DBUser           string
	DBPassWord       string
	DBName           string
}

type jsonConfig struct {
	RefreshTokenExpDurationHour int `json:"refresh_token_exp_duration_hour"`
	AccessTokenExpDurationMin   int `json:"access_token_exp_duration_min"`
}

type ServerConfig struct {
	envConf  envConfig
	jsonConf jsonConfig
}

func getEnvVar(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", errors.New("environment variable not found: " + key)
	}
	return value, nil
}

func QueryServerConfig() (*ServerConfig, error) {
	authPort, err := getEnvVar("AUTH_PORT")
	if err != nil {
		return nil, err
	}

	mode, err := getEnvVar("MODE")
	if err != nil {
		return nil, err
	}

	authAllowOrigins, err := getEnvVar("AUTH_ALLOW_ORIGINS")
	if err != nil {
		return nil, err
	}

	privateKeyFile, err := getEnvVar("PRIVATE_KEY_FILE")
	if err != nil {
		return nil, err
	}

	publicKeyFile, err := getEnvVar("PUBLIC_KEY_FILE")
	if err != nil {
		return nil, err
	}

	dbHost, err := getEnvVar("DB_HOST")
	if err != nil {
		return nil, err
	}

	dbUser, err := getEnvVar("DB_USER")
	if err != nil {
		return nil, err
	}

	dbPassword, err := getEnvVar("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbName, err := getEnvVar("DB_NAME")
	if err != nil {
		return nil, err
	}

	jsonPath, err := getEnvVar("CONFIG_JSON_FILE")
	if err != nil {
		return nil, err
	}

	// Open the JSON file
	file, err := os.Open(jsonPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a JSON decoder and decode the file contents into jsonConf
	var jsonConf jsonConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&jsonConf); err != nil {
		return nil, err
	}

	return &ServerConfig{
		envConf: envConfig{
			AuthPort:         authPort,
			Mode:             mode,
			AuthAllowOrigins: strings.Split(authAllowOrigins, ","),
			PrivateKeyFile:   privateKeyFile,
			PublicKeyFile:    publicKeyFile,
			DBHost:           dbHost,
			DBUser:           dbUser,
			DBPassWord:       dbPassword,
			DBName:           dbName,
		},
		jsonConf: jsonConf,
	}, nil
}
