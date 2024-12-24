package main

import (
	//"context"
	"log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/abhishek-kamat-nutanix/orchestrator/move/proto"
	   
)

var addr string = "localhost:50051"
// 10.15.170.150:30051 nke
// 10.46.63.221:30051 ocp
// 10.15.168.215:30051 new ocp 

func main() {

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err!=nil {
		log.Fatalf("Failed to connect %v\n", err)
	}

	defer conn.Close()

	c := pb.NewMoveServiceClient(conn)

	doMigrateApp(c)

}