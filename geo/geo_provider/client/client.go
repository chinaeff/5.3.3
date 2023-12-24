package main

import (
	"context"
	"fmt"
	"geotask_pprof/geo"
	"log"

	"google.golang.org/grpc"
)

// @Summary Connect to GeoService server and perform AddressSearch
// @Description Connects to the GeoService gRPC server and performs an AddressSearch using the provided input.
// @Accept json
// @Produce json
// @Success 200 {object} geo.SearchResponse
// @Router /address/search [get]
func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := geo.NewGeoServiceClient(conn)

	searchResponse, err := client.AddressSearch(context.Background(), &geo.SearchRequest{Input: "Search Input"})
	if err != nil {
		log.Fatalf("AddressSearch failed: %v", err)
	}
	fmt.Println("AddressSearch Response:", searchResponse)

	geoCodeResponse, err := client.GeoCode(context.Background(), &geo.GeoCodeRequest{Lat: "Latitude", Lng: "Longitude"})
	if err != nil {
		log.Fatalf("GeoCode failed: %v", err)
	}
	fmt.Println("GeoCode Response:", geoCodeResponse)
}
