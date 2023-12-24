package main

import "testing"

func TestAddressSearch(t *testing.T) {
	provider := NewGeoProvider()

	input := "TestAddress"
	result, err := provider.AddressSearch(input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 1 || result[0].Addr != input {
		t.Errorf("Unexpected result. Expected [%s], got %v", input, result)
	}
}

func TestGeoCode(t *testing.T) {
	provider := NewGeoProvider()

	lat, lng := "12.345", "67.890"
	expectedResult := lat + lng

	result, err := provider.GeoCode(lat, lng)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 1 || result[0].Addr != expectedResult {
		t.Errorf("Unexpected result. Expected [%s], got %v", expectedResult, result)
	}
}

func TestGeoServiceAddressSearch(t *testing.T) {
	provider := NewGeoProvider()
	geoService := &GeoService{Provider: provider}

	req := "TestAddress"
	var res []*Address
	err := geoService.AddressSearch(req, &res)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(res) != 1 || res[0].Addr != req {
		t.Errorf("Unexpected result. Expected [%s], got %v", req, res)
	}
}

func TestGeoServiceGeoCode(t *testing.T) {
	provider := NewGeoProvider()
	geoService := &GeoService{Provider: provider}

	lat, lng := "12.345", "67.890"
	expectedResult := lat + lng

	req := [2]string{lat, lng}
	var res []*Address
	err := geoService.GeoCode(req, &res)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(res) != 1 || res[0].Addr != expectedResult {
		t.Errorf("Unexpected result. Expected [%s], got %v", expectedResult, res)
	}
}
