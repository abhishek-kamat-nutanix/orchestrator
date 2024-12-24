package main

import (
	"context"
	"flag"
	"log"
	"os"

	pb "github.com/abhishek-kamat-nutanix/orchestrator/move/proto"

)

func doMigrateApp(c pb.MoveServiceClient){
	log.Println("doAppMigrate was invoked")

	kubeconfigPath := flag.String("kubeconfig", "/home/nutanix/kubeconfig", "location to your kubeconfig file")
    flag.Parse() // Ensure flags are parsed before use

	// Read the Kubernetes config
    config, err := os.ReadFile(*kubeconfigPath)
	if err != nil {
        log.Fatalf("error reading kubeconfig file: %v", err)
    }
	
	res, err := c.MigrateApp(context.Background(), &pb.AppRequest{
		Namespace: "wordpress", 
		Resources: "deployments,secrets,svc",
		Labels: "app=wordpress",
		ServerAddr: "10.15.170.48:50051",
		Kubeconfig: string(config),
		ReaderAddr: "10.15.168.215:30051", //10.15.168.215:30051 ocp 10.15.170.150:30051 nke 10.15.174.11:30051 ocp2

	})

	if err != nil {
		log.Fatalf("could not MigrateApp: %v\n", err)
	}

	log.Printf("Message recieved from Move server: %v\n", res.Message)
}