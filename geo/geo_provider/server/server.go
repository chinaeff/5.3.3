package main

import (
	"context"
	"fmt"
	"geotask_pprof/geo"
	"google.golang.org/grpc"
	"log"
	"net"
)

// @title GeoService API
// @version 1.0
// @description GeoService provides geo-related functionality.
// @host localhost:50051
// @BasePath /api
type geoServer struct{}

// AddressSearch is a function to search for addresses.
// @Summary Search for addresses.
// @Description Search for addresses based on the provided request.
// @Accept json
// @Produce json
// @Param request body geo.SearchRequest true "Search request parameters"
// @Success 200 {object} geo.SearchResponse
// @Router /address/search [post]
func (s *geoServer) AddressSearch(ctx context.Context, req *geo.SearchRequest) (*geo.SearchResponse, error) {
	return &geo.SearchResponse{
		Addresses: []*geo.Address{},
	}, nil
}

// GeoCode is a function to perform geocoding.
// @Summary Perform geocoding.
// @Description Perform geocoding based on the provided coordinates.
// @Accept json
// @Produce json
// @Param request body geo.GeoCodeRequest true "Geocoding request parameters"
// @Success 200 {object} geo.GeoCodeResponse
// @Router /geo/code [post]
func (s *geoServer) GeoCode(ctx context.Context, req *geo.GeoCodeRequest) (*geo.GeoCodeResponse, error) {
	return &geo.GeoCodeResponse{
		Addresses: []*geo.Address{},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	geoServer := geo.GeoServiceServer(nil)
	grpcServer := grpc.NewServer()
	geo.RegisterGeoServiceServer(grpcServer, geoServer)
	fmt.Println("gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
