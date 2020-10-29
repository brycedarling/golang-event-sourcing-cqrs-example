package config

import (
	"fmt"
	"os"
)

// Env ...
type Env struct {
	Env                        string
	Port                       string
	EventStoreConnectionString string
	QueryConnectionString      string
}

// NewEnv ...
func NewEnv() (*Env, error) {
	appEnv, err := getEnv("APP_ENV")
	if err != nil {
		return nil, err
	}
	port, err := getEnv("PORT")
	if err != nil {
		return nil, err
	}
	eventStoreConnString, err := getEnv("EVENT_STORE_CONNECTION_STRING")
	if err != nil {
		return nil, err
	}
	queryConnString, err := getEnv("QUERY_CONNECTION_STRING")
	if err != nil {
		return nil, err
	}
	return &Env{
		Env:                        appEnv,
		Port:                       port,
		EventStoreConnectionString: eventStoreConnString,
		QueryConnectionString:      queryConnString,
	}, nil
}

func getEnv(key string) (string, error) {
	envVar, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("Missing required %s environment variable", key)
	}
	return envVar, nil
}
