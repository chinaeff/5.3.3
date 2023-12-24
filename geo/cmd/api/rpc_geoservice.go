package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// @title GeoService
// @version 1.0
// @description GeoService provides geo-related functionality.
// @host localhost:50051
// @BasePath /api
//
//	@contact {
//	  url: "https://github.com/chinaeff/5.1.2,
//	}
//
//	@license {
//	  name: "MIT",
//	  url: "https://opensource.org/licenses/MIT",
//	}

type Address struct {
	Addr string
}

type GeoProvider interface {
	AddressSearch(input string) ([]*Address, error)
	GeoCode(lat, lng string) ([]*Address, error)
}

type SimpleGeoProvider struct{}

func (sgp *SimpleGeoProvider) AddressSearch(input string) ([]*Address, error) {
	return []*Address{{Addr: input}}, nil
}

func (sgp *SimpleGeoProvider) GeoCode(lat, lng string) ([]*Address, error) {
	addr := lat + lng
	return []*Address{{Addr: addr}}, nil
}

func NewGeoProvider() *SimpleGeoProvider {
	return &SimpleGeoProvider{}
}

type GeoService struct {
	Provider GeoProvider
}

func (g *GeoService) AddressSearch(req string, res *[]*Address) error {
	*res, _ = g.Provider.AddressSearch(req)
	return nil
}

func (g *GeoService) GeoCode(req [2]string, res *[]*Address) error {
	*res, _ = g.Provider.GeoCode(req[0], req[1])
	return nil
}

func startJSONRPCServer() {
	server, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println("Failed to start JSON-RPC server:", err)
		return
	}
	fmt.Println("JSON-RPC server started on :50051")

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accept:", err)
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
}

func startRPCServer() {
	server, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println("Failed to start RPC server:", err)
		return
	}
	fmt.Println("RPC server started on :50051")

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accept:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
