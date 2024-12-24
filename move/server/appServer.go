package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	//"time"

	pb "github.com/abhishek-kamat-nutanix/orchestrator/move/proto"
	//types "k8s.io/apimachinery/pkg/types"

	//v2 "github.com/kubernetes-csi/external-snapshotter/client/v8/apis/volumesnapshot/v1"
	//"github.com/kubernetes-csi/external-snapshotter/client/v8/clientset/versioned"
	//batchv1 "k8s.io/api/batch/v1"
	//v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	pr "github.com/abhishek-kamat-nutanix/go-client/reader/proto"
	"k8s.io/client-go/tools/clientcmd"
)

func (s *Server) MigrateApp(ctx context.Context, in *pb.AppRequest) (*pb.AppResponse, error) {
    
    fmt.Printf("MigrateApp was invoked\n")

	kubeconfig := flag.String("kubeconfig", "/home/nutanix/kubeconfig", "location to your kubeconfig file")
    flag.Parse() // Ensure flags are parsed before use

	  // Build the Kubernetes config
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        return nil, fmt.Errorf("error getting kubeconfig: %v", err)
    }


	//kubeconfig := in.Kubeconfig 

    // config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
    // if err != nil {
    //     return nil, fmt.Errorf("error getting kubeconfig: %v", err)
    // }

    // Create Kubernetes clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, fmt.Errorf("error creating Kubernetes client: %v", err)
    }

    namespace := in.Namespace
    labl := in.Labels
    serverip := in.ServerAddr
    rsc := in.Resources
	raddr := in.ReaderAddr

    // List PersistentVolumeClaims in the namespace
    pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labl})
    if err != nil {
        return nil, fmt.Errorf("error listing PVCs in namespace %s: %v", namespace, err)
    }

	conn, err := grpc.NewClient(raddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err!=nil {
		log.Fatalf("Failed to connect %v\n", err)
	}

	defer conn.Close()

	c := pr.NewReaderServiceClient(conn)

    // migrate each volume in list
    for _,item := range pvc.Items {
    fmt.Printf("found PVC in namespace %s: %s\n", namespace, item.Name)
	

    res, err := c.MigrateVolume(context.Background(), &pr.VolumeRequest{
        				ServerAddr: serverip, // writer server address on target cluster
        				BackupName: item.Name, // volume on source cluster
        				Namespace: namespace,
        			})
        if err != nil {
            log.Printf("error migrating pv %v error: %v\n",item.Name,err)
        }

		log.Printf("Message recieved from readers server: %v\n", res.Message)

    }

	res, err := c.MigrateConfig(context.Background(), &pr.ConfigRequest{
		Namespace: namespace,
        Resources: rsc,
        Labels: labl,
        ServerAddr: serverip,
	})

	if err != nil {
		log.Fatalf("could not MigrateConfig: %v\n", err)
	}

	log.Printf("Message recieved from readers server: %v\n", res.Message)

    return &pb.AppResponse{
        Message: "Migrate App has Completed",
    }, nil
}
