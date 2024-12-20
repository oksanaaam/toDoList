package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBType              string
	DBConnectionString  string
	MongoURI            string
	MongoDBName         string
	MongoCollectionName string
	ServerAddress       string
}

func LoadConfig() *Config {
	err := godotenv.Load() // for local running
	// err := godotenv.Load("/app/.env") // for running in docker
	if err != nil {
		log.Fatal("Cannot find file .env, err: ", err)
	}

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		log.Fatal("DB_TYPE not found, specify 'postgres' or 'mongo'")
	}

	var dbConnectionString, mongoURI, mongoDBName, mongoCollectionName string
	if dbType == "postgres" {
		dbConnectionString = os.Getenv("DB_CONNECTION_STRING")
		if dbConnectionString == "" {
			log.Fatal("DB_CONNECTION_STRING not found")
		}
	} else if dbType == "mongo" {
		mongoURI = os.Getenv("MONGO_URI")
		if mongoURI == "" {
			log.Fatal("MONGO_URI not found")
		}
		mongoDBName = os.Getenv("MONGO_DB_NAME")
		if mongoDBName == "" {
			log.Fatal("MONGO_DB_NAME not found")
		}
		mongoCollectionName = os.Getenv("MONGO_COLLECTION_NAME")
		if mongoCollectionName == "" {
			log.Fatal("MONGO_COLLECTION_NAME not found")
		}
	} else {
		log.Printf("Unsupported DB type: %v", dbType)
	}

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "localhost:8080"
	}

	return &Config{
		DBType:              dbType,
		DBConnectionString:  dbConnectionString,
		MongoURI:            mongoURI,
		MongoDBName:         mongoDBName,
		MongoCollectionName: mongoCollectionName,
		ServerAddress:       serverAddress,
	}
}
