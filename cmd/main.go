package main

import (
	"github.com/SaidovZohid/minio-medias/api"
	"github.com/SaidovZohid/minio-medias/config"
	"github.com/SaidovZohid/minio-medias/pkg/logging"
	"github.com/SaidovZohid/minio-medias/storage/minio"
	_ "github.com/lib/pq"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Println("logger initialized")

	logger.Println("config initializing")
	cfg := config.GetConfig(".")

	minioStorage, err := minio.NewStorage(&cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	apiServer := api.New(&api.RouterOptions{
		Cfg:          &cfg,
		MinioStorage: minioStorage,
		Logger:       logger,
	})

	if err := apiServer.Run(cfg.HttpPort); err != nil {
		logger.Printf("server not started in %s port\n", cfg.HttpPort)
	}

	// psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	cfg.Postgres.Host,
	// 	cfg.Postgres.Port,
	// 	cfg.Postgres.User,
	// 	cfg.Postgres.Password,
	// 	cfg.Postgres.Database,
	// )

	// psqlConn, err := sqlx.Connect("postgres", psqlUrl)
	// if err != nil {
	// 	log.Fatalf("failed to connect to database: %v", err)
	// }

	// _ = psqlConn
}
