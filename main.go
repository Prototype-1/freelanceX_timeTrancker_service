package main

import (
	"log"
	"net"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	time_logpb "github.com/Prototype-1/freelanceX_timeTrancker_service/proto"
	"github.com/Prototype-1/freelanceX_timeTrancker_service/internal/repository"
	"github.com/Prototype-1/freelanceX_timeTrancker_service/internal/service"
	"github.com/Prototype-1/freelanceX_timeTrancker_service/internal/model"
	"github.com/Prototype-1/freelanceX_timeTrancker_service/config" 
)

func main() {
	cfg := config.LoadConfig()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&model.TimeLog{}); err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	repo := repository.NewTimeLogRepository(db)
	svc := service.NewTimeLogService(repo)

	listener, err := net.Listen("tcp", cfg.ServerPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.ServerPort, err)
	}
	grpcServer := grpc.NewServer()
	time_logpb.RegisterTimeLogServiceServer(grpcServer, svc)


	log.Printf("Starting server on %s...\n", cfg.ServerPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
