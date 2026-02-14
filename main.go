package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "hivemind/proto"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
	timeoutSeconds = 15
	cleanupInterval = 3
)

type WorkerInfo struct{
	ID	string
	LastSeen time.Time
	WorkerHostname string
}

type Coordinator struct{
	pb.UnimplementedWorkerServiceServer
	mu 		sync.Mutex
	workers map[string]*WorkerInfo
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		workers: make(map[string]*WorkerInfo),
	}
}

//Register

func(c *Coordinator) RegisterWorker(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error){
	c.mu.Lock()
	defer c.mu.Unlock()

	c.workers[req.WorkerId] = &WorkerInfo{
		ID:		req.WorkerId,
		LastSeen: time.Now(),
		WorkerHostname: req.WorkerHostname,
	}

	fmt.Printf("Worker joined: %s (%s) \n", req.WorkerId, req.WorkerHostname)

	return &pb.RegisterResponse{
		HeartbeatInterval: 5,
	}, nil

}

//Heartbeat
func(c *Coordinator) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error){
	c.mu.Lock()
	defer c.mu.Unlock()

	worker, exists := c.workers[req.WorkerId]
	if !exists {
		return &pb.HeartbeatResponse{Ok: false}, nil
	}

	worker.LastSeen = time.Now()
	fmt.Printf("Heartbeat from: %s\n", req.WorkerId)

	return &pb.HeartbeatResponse{Ok:true}, nil
}

//Timeout cleaner

func (c *Coordinator) cleanupDeadWorkers(){
	for{
		time.Sleep(cleanupInterval * time.Second)

		c.mu.Lock()
		now := time.Now()

		for id, worker := range c.workers {
			//basically now - lastSeen > timeout seconds (raw 15 * type = 15 seconds)
			if now.Sub(worker.LastSeen) > timeoutSeconds*time.Second {
				fmt.Printf("Worker time out: %s \n", id)
				delete(c.workers, id)
			}
			
		}

		c.mu.Unlock()
	}
}

//Main driver
func main(){
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	server := grpc.NewServer()
	coord := NewCoordinator()

	pb.RegisterWorkerServiceServer(server, coord)

	go coord.cleanupDeadWorkers()

	fmt.Println("Coordinator running on the port", port)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}


}