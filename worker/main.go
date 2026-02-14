package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "hivemind/proto" 

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Setup connection to the Coordinator
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewWorkerServiceClient(conn)

	// 2. Prepare registration data
	workerID := "worker-01" // In a real app, maybe use a UUID
	hostname, _ := os.Hostname()

	// 3. Register with the Coordinator
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := client.RegisterWorker(ctx, &pb.RegisterRequest{
		WorkerId:       workerID,
		WorkerHostname: hostname,
		Message:        "Hello from the worker!",
	})

	if err != nil {
		log.Fatalf("Could not register: %v", err)
	}

	fmt.Printf("Registered! Heartbeat interval: %d seconds\n", res.HeartbeatInterval)

	// 4. Start Heartbeat loop
	interval := time.Duration(res.HeartbeatInterval) * time.Second
	for {
		time.Sleep(interval)
		
		hCtx, hCancel := context.WithTimeout(context.Background(), time.Second*2)
		_, err := client.Heartbeat(hCtx, &pb.HeartbeatRequest{
			WorkerId: workerID,
		})
		
		if err != nil {
			fmt.Printf("Heartbeat failed: %v\n", err)
		} else {
			fmt.Println("Heartbeat sent successfully")
		}
		hCancel()
	}
}