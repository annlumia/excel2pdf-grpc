package main

import (
	"context"

	"github.com/annlumia/excel2pdf-grpc/service"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := service.NewClientService("localhost:8355", "/home/hotman/Downloads/M4M Modbus map - v.1.3N.xlsx", 1024)
	fileName, err := client.Convert(ctx)
	if err != nil {
		panic(err)
	}

	println(fileName)
}
