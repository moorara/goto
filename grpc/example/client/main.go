package main

import (
	"context"
	"io"
	"net/http"

	xgrpc "github.com/moorara/goto/grpc"
	"github.com/moorara/goto/grpc/example/zonePB"
	"github.com/moorara/goto/log"
	"github.com/moorara/goto/metrics"
	"github.com/moorara/goto/trace"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const address = "localhost:8080"

func getContainingZone(client zonePB.ZoneManagerClient) {
	ctx := context.Background()

	stream, err := client.GetContainingZone(ctx)
	if err != nil {
		panic(err)
	}

	locations := []*zonePB.Location{
		{
			Latitude:  43.662892,
			Longitude: -79.395684,
		},
		{
			Latitude:  43.658776,
			Longitude: -79.379327,
		},
	}

	for _, loc := range locations {
		err := stream.Send(loc)
		if err != nil {
			panic(err)
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		panic(err)
	}
}

func getPlacesInZone(client zonePB.ZoneManagerClient) {
	ctx := context.Background()
	zone := &zonePB.Zone{
		Location: &zonePB.Location{
			Latitude:  43.645844,
			Longitude: -79.379742,
		},
		Radius: 1200,
	}

	_, err := client.GetPlacesInZone(ctx, zone)
	if err != nil {
		panic(err)
	}
}

func getUsersInZone(client zonePB.ZoneManagerClient) {
	ctx := context.Background()
	zone := &zonePB.Zone{
		Location: &zonePB.Location{
			Latitude:  43.645844,
			Longitude: -79.379742,
		},
		Radius: 1200,
	}

	stream, err := client.GetUsersInZone(ctx, zone)
	if err != nil {
		panic(err)
	}

	for {
		_, err := stream.Recv()
		if err != nil && err != io.EOF {
			panic(err)
		}

		if err == io.EOF {
			return
		}
	}
}

func getUsersInZones(client zonePB.ZoneManagerClient) {
	ctx := context.Background()
	zones := []*zonePB.Zone{
		{
			Location: &zonePB.Location{
				Latitude:  45.424688,
				Longitude: -75.699565,
			},
			Radius: 1500,
		},
		{
			Location: &zonePB.Location{
				Latitude:  43.472920,
				Longitude: -80.542378,
			},
			Radius: 1000,
		},
	}

	stream, err := client.GetUsersInZones(ctx)
	if err != nil {
		panic(err)
	}

	waitc := make(chan struct{})

	// Receiving
	go func() {
		for {
			_, err := stream.Recv()
			if err != nil && err != io.EOF {
				panic(err)
			}

			if err == io.EOF {
				close(waitc)
				return
			}
		}
	}()

	// Sending
	for _, zone := range zones {
		err := stream.Send(zone)
		if err != nil {
			panic(err)
		}
	}

	err = stream.CloseSend()
	if err != nil {
		panic(err)
	}

	<-waitc
}

func main() {
	// Create a logger
	logger := log.NewLogger(log.Options{
		Name:        "client",
		Environment: "dev",
		Region:      "us-east-1",
		Component:   "zone-client",
	})

	// Create a metrics factory
	mf := metrics.NewFactory(metrics.FactoryOptions{})

	// Create a tracer
	tracer, closer, _ := trace.NewTracer(trace.Options{Name: "zone-client"})
	defer closer.Close()

	// Create a gRPC interceptor for observability
	i := xgrpc.NewClientObservabilityInterceptor(logger, mf, tracer)

	optInsecure := grpc.WithInsecure()
	optUnaryInterceptor := grpc.WithUnaryInterceptor(i.UnaryInterceptor)
	optStreamInterceptor := grpc.WithStreamInterceptor(i.StreamInterceptor)
	conn, err := grpc.Dial(address, optInsecure, optUnaryInterceptor, optStreamInterceptor)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := zonePB.NewZoneManagerClient(conn)
	logger.Info("message", "client connected to server", "address", address)

	getContainingZone(client)
	getPlacesInZone(client)
	getUsersInZone(client)
	getUsersInZones(client)

	http.Handle("/metrics", promhttp.Handler())
	logger.Info("message", "starting http server on :8082")
	panic(http.ListenAndServe(":8082", nil))
}
