package main

import (
	"context"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"

	xgrpc "github.com/moorara/goto/grpc"
	"github.com/moorara/goto/grpc/example/zonePB"
	"github.com/moorara/goto/log"
	"github.com/moorara/goto/metrics"
	"github.com/moorara/goto/trace"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const grpcPort = ":10080"
const httpPort = ":10081"

// ZoneServer is an implementation of zonePB.ZoneManagerServer
type ZoneServer struct{}

// GetContainingZone the zone containing all the given locations
func (s *ZoneServer) GetContainingZone(stream zonePB.ZoneManager_GetContainingZoneServer) error {
	// A random delay between 5ms to 50ms
	d := 5 + rand.Intn(45)
	time.Sleep(time.Duration(d) * time.Millisecond)

	logger, _ := xgrpc.LoggerFromContext(stream.Context())
	logger.Info("message", "GetContainingZone handled!")

	for {
		_, err := stream.Recv()
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF {
			return stream.SendAndClose(&zonePB.Zone{
				Location: &zonePB.Location{
					Latitude:  43.661370,
					Longitude: 79.383096,
				},
				Radius: 1200,
			})
		}
	}
}

// GetPlacesInZone returns all places in a zone
func (s *ZoneServer) GetPlacesInZone(ctx context.Context, zone *zonePB.Zone) (*zonePB.GetPlacesResponse, error) {
	// A random delay between 5ms to 50ms
	d := 5 + rand.Intn(45)
	time.Sleep(time.Duration(d) * time.Millisecond)

	logger, _ := xgrpc.LoggerFromContext(ctx)
	logger.Info("message", "GetPlacesInZone handled!")

	return &zonePB.GetPlacesResponse{
		Zone: zone,
		Places: []*zonePB.Place{
			{
				Id:   "1111-1111-1111-1111",
				Name: "CN Tower",
				Location: &zonePB.Location{
					Latitude:  43.642581,
					Longitude: -79.386907,
				},
			},
			{
				Id:   "2222-2222-2222-2222",
				Name: "Yonge-Dundas Square",
				Location: &zonePB.Location{
					Latitude:  43.656095,
					Longitude: -79.380171,
				},
			},
		},
	}, nil
}

// GetUsersInZone returns all the users entering a zone
func (s *ZoneServer) GetUsersInZone(zone *zonePB.Zone, stream zonePB.ZoneManager_GetUsersInZoneServer) error {
	// A random delay between 5ms to 50ms
	d := 5 + rand.Intn(45)
	time.Sleep(time.Duration(d) * time.Millisecond)

	logger, _ := xgrpc.LoggerFromContext(stream.Context())
	logger.Info("message", "GetUsersInZone handled!")

	users := []*zonePB.UserInZone{
		{
			Location: &zonePB.Location{
				Latitude:  43.645710,
				Longitude: -79.376115,
			},
			User: &zonePB.User{
				Id:   "aaaa-aaaa-aaaa-aaaa",
				Name: "Milad",
			},
		},
		{
			Location: &zonePB.Location{
				Latitude:  43.646075,
				Longitude: -79.376294,
			},
			User: &zonePB.User{
				Id:   "bbbb-bbbb-bbbb-bbbb",
				Name: "Mona",
			},
		},
	}

	for _, uz := range users {
		err := stream.Send(uz)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetUsersInZones returns all the users entering any of the given zones
func (s *ZoneServer) GetUsersInZones(stream zonePB.ZoneManager_GetUsersInZonesServer) error {
	// A random delay between 5ms to 50ms
	d := 5 + rand.Intn(45)
	time.Sleep(time.Duration(d) * time.Millisecond)

	logger, _ := xgrpc.LoggerFromContext(stream.Context())
	logger.Info("message", "GetUsersInZones handled!")

	users := []*zonePB.UserInZone{
		{
			Location: &zonePB.Location{
				Latitude:  43.645710,
				Longitude: -79.376115,
			},
			User: &zonePB.User{
				Id:   "aaaa-aaaa-aaaa-aaaa",
				Name: "Milad",
			},
		},
		{
			Location: &zonePB.Location{
				Latitude:  43.646075,
				Longitude: -79.376294,
			},
			User: &zonePB.User{
				Id:   "bbbb-bbbb-bbbb-bbbb",
				Name: "Mona",
			},
		},
	}

	for {
		_, err := stream.Recv()
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF {
			return nil
		}

		i := rand.Intn(2)
		err = stream.Send(users[i])
		if err != nil {
			return err
		}
	}
}

func main() {
	// Create a logger
	logger := log.NewLogger(log.Options{
		Name:        "server",
		Environment: "dev",
		Region:      "us-east-1",
		Component:   "zone-server",
	})

	// Create a metrics factory
	mf := metrics.NewFactory(metrics.FactoryOptions{})

	// Create a tracer
	tracer, closer, _ := trace.NewTracer(trace.Options{Name: "server"})
	defer closer.Close()

	// Create a gRPC interceptor
	i := xgrpc.NewServerInterceptor(logger, mf, tracer)

	optUnaryInterceptor := grpc.UnaryInterceptor(i.UnaryInterceptor)
	optStreamInterceptor := grpc.StreamInterceptor(i.StreamInterceptor)
	server := grpc.NewServer(optUnaryInterceptor, optStreamInterceptor)
	zonePB.RegisterZoneManagerServer(server, &ZoneServer{})

	// Start HTTP server for exposing metrics
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Info("message", "starting http server ...", "port", httpPort)
		panic(http.ListenAndServe(httpPort, nil))
	}()

	conn, err := net.Listen("tcp", grpcPort)
	if err != nil {
		panic(err)
	}

	logger.Info("message", "starting grpc server ...", "port", grpcPort)
	panic(server.Serve(conn))
}
