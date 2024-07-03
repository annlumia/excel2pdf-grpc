package main

import (
	"fmt"
	"log"
	"net"

	o2pdf "github.com/annlumia/excel2pdf-grpc/proto"
	"github.com/annlumia/excel2pdf-grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run() *grpc.Server {
	svr := grpc.NewServer()
	svc := service.NewServerService()

	o2pdf.RegisterConverterServiceServer(svr, svc)
	reflection.Register(svr)

	listener, err := net.Listen("tcp", ":8354")
	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		fmt.Println("Server is running on port :8354")
		if err := svr.Serve(listener); err != nil {
			log.Fatal(err.Error())
		}
	}()

	return svr
}
