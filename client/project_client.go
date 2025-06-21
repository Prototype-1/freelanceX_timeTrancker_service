package client

import (
	projectpb "github.com/Prototype-1/freelanceX_timeTrancker_service/proto/crm_service"
	"google.golang.org/grpc"
	 "google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

var ProjectClient projectpb.ProjectServiceClient

func InitProjectServiceClient() {
    addr := os.Getenv("PROJECT_SERVICE_GRPC_ADDR")
    conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Failed to connect to ProjectService at %s: %v", addr, err)
    }
    
    ProjectClient = projectpb.NewProjectServiceClient(conn)
}
