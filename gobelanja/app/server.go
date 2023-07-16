package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
}

type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBDriver   string
}

func (server *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("Welcome to " + appConfig.AppName)

	server.InitializeDB(dbConfig)
	server.Router = mux.NewRouter()
	server.InitializeRoutes()

}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func (server *Server) InitializeDB(dbConfig DBConfig) {
	var err error
	var dsn string

	if dbConfig.DBDriver == "mysql" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
	}

	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect database")
	}

	for _, model := range RegisterModels() {
		err = server.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal(err)
		}

	}
	fmt.Println("Database Migrated sukses.")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Run() {
	var server = Server{}
	var appConfig = AppConfig{}
	var dbConfig = DBConfig{}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error on loading .env file")
	}

	appConfig.AppName = getEnv("APP_NAME", "GoBelanjaApp")
	appConfig.AppEnv = getEnv("APP_ENV", "development")
	appConfig.AppPort = getEnv("APP_PORT", "9000")

	dbConfig.DBHost = getEnv("DB_HOST", "localhost")
	dbConfig.DBUser = getEnv("DB_USER", "postgres")
	dbConfig.DBPassword = getEnv("DB_PASSWORD", "1mAu200Ribu")
	dbConfig.DBName = getEnv("DB_NAME", "balqis_gobelanjadb")
	dbConfig.DBPort = getEnv("DB_PORT", "5432")

	server.Initialize(appConfig, dbConfig)
	server.Run(":" + appConfig.AppPort)

}
