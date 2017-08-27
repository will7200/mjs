package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/metadata"

	"github.com/will7200/mjs/apischeduler"

	"github.com/will7200/mjs/apischeduler/grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAPISchedulerClient(conn)
	//rr, errr := c.Start(context.Background(), &pb.StartRequest{})
	//fmt.Printf("Error %+v\n", errr)
	// Contact the server and print out its response.
	ctx := context.Background()
	ctx = metadata.NewContext(ctx,
		metadata.Pairs(apischeduler.JobUniqueness, "UNIQUE"))
	fmt.Println(ctx.Value(apischeduler.JobUniqueness))
	r, err := c.Add(ctx, &pb.AddRequest{Reqjob: &pb.Job{
		Name:     "Test",
		Command:  []string{"Here", "we", "go"},
		Schedule: "R/2017-09-01T01:01:01/PT10S",
	}})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %+v", r)
}
