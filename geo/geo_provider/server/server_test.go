package main

import (
	"context"
	"geotask_pprof/geo"
	"testing"
)

func TestAddressSearch(t *testing.T) {
	server := &geoServer{}

	ctx := context.Background()

	request := &geo.SearchRequest{
		Input: "Search Input",
	}

	response, err := server.AddressSearch(ctx, request)
	if err != nil {
		t.Fatalf("AddressSearch failed: %v", err)
	}

	if response == nil {
		t.Fatal("Expected non-nil response")
	}

	if len(response.Addresses) != 0 {
		t.Fatalf("Expected empty Addresses array, got: %v", response.Addresses)
	}
}

func TestGeoCode(t *testing.T) {
	server := &geoServer{}

	ctx := context.Background()

	request := &geo.GeoCodeRequest{
		Lat: "Latitude",
		Lng: "Longitude",
	}

	response, err := server.GeoCode(ctx, request)
	if err != nil {
		t.Fatalf("GeoCode failed: %v", err)
	}

	if response == nil {
		t.Fatal("Expected non-nil response")
	}

	if len(response.Addresses) != 0 {
		t.Fatalf("Expected empty Addresses array, got: %v", response.Addresses)
	}
}
