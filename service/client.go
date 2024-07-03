package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	o2pdf "github.com/annlumia/excel2pdf-grpc/proto"
	"google.golang.org/grpc"
)

type o2pdfClient struct {
	client    o2pdf.ConverterServiceClient
	addr      string
	filePath  string
	batchSize int
}

func NewClientService(addr string, filePath string, batchSize int) *o2pdfClient {
	return &o2pdfClient{
		addr:      addr,
		filePath:  filePath,
		batchSize: batchSize,
	}
}

func (o *o2pdfClient) Convert(ctx context.Context) (file_name string, err error) {
	log.Println(o.addr, o.filePath)
	conn, err := grpc.Dial(o.addr, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	o.client = o2pdf.NewConverterServiceClient(conn)

	pdfFilename, err := o.upload(ctx)
	if err != nil {
		return "", err
	}

	return pdfFilename, nil
}

func (o *o2pdfClient) upload(ctx context.Context) (string, error) {
	stream, err := o.client.Convert(ctx)
	if err != nil {
		return "", err
	}

	waitc := make(chan struct{})

	pdfFile := NewFile()
	fileSize := uint32(0)
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				log.Fatalf("Failed to receive a converted file: %v", err)
			}

			if pdfFile.FilePath == "" {
				baseFile := path.Base(strings.ReplaceAll(in.FileName, "\\", "/"))
				pdfFile.SetFile(baseFile, ".temp")
			}

			chunk := in.GetChunk()
			fileSize += uint32(len(chunk))
			log.Printf("received a chunk with size: %d\n", fileSize)
			if err := pdfFile.Write(chunk); err != nil {
				fmt.Println(err.Error())
				return
			}

		}
	}()

	file, err := os.Open(o.filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, o.batchSize)
	batchNumber := 1

	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		if err := stream.Send(&o2pdf.ConvertRequest{
			FileName: file.Name(),
			Chunk:    buf[:n],
		}); err != nil {
			return "", err
		}

		log.Printf("sent batch number %d\n", batchNumber)
		batchNumber++
	}

	err = stream.CloseSend()
	if err != nil {
		return "", err
	}

	<-waitc

	return pdfFile.FilePath, nil
}
