package service

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/annlumia/excel2pdf-grpc/office2pdf"
	o2pdf "github.com/annlumia/excel2pdf-grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type o2pdfServer struct {
	o2pdf.UnimplementedConverterServiceServer
}

func (o *o2pdfServer) convert(filename string) (string, error) {
	excel2pdf := office2pdf.Excel{}
	pdf, err := excel2pdf.Export(filename)
	if err != nil {
		return "", err
	}

	return pdf, nil
}

func (o *o2pdfServer) Convert(stream o2pdf.ConverterService_ConvertServer) error {
	file := NewFile()
	var fileSize uint32

	fileSize = 0
	defer func() {
		if err := file.OutputFile.Close(); err != nil {
			log.Printf("error %s", err.Error())
		}
	}()

	randomFileName := shortRandom()

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Y", err.Error())
			return status.Error(codes.Internal, err.Error())
		}

		if file.FilePath == "" {
			randomFileName += "_" + path.Base(req.GetFileName())
			file.SetFile(randomFileName, "temp")
		}

		chunk := req.GetChunk()
		fileSize += uint32(len(chunk))
		log.Printf("received a chunk with size: %d\n", fileSize)
		if err := file.Write(chunk); err != nil {
			fmt.Println(err.Error())
			return status.Error(codes.Internal, err.Error())
		}

	}

	cwd, _ := os.Getwd()

	excelFileName := filepath.Join(cwd, "temp", randomFileName)
	result, err := o.convert(excelFileName)
	if err != nil {
		return err
	}

	go func() {
		time.Sleep(time.Second * 2)
		os.Remove(excelFileName)
	}()

	pdfFile, err := os.Open(result)
	if err != nil {
		return err
	}
	defer func() {
		pdfFile.Close()
		os.Remove(result)
	}()

	buf := make([]byte, 1024)
	batchNumber := 1

	for {
		n, err := pdfFile.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		err = stream.Send(&o2pdf.ConvertResponse{
			FileName: result,
			Chunk:    buf[:n],
		})
		if err != nil {
			return err
		}

		log.Printf("sent batch number %d\n", batchNumber)
		batchNumber++
	}

	fmt.Println("Everithing is done")
	return nil
}

func NewServerService() *o2pdfServer {
	return &o2pdfServer{}
}
