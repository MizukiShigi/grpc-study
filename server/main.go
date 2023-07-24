package main

import (
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

const storagePath = "/Users/shifi/dev/grpc-lesson/storage"

type server struct {
	pb.UnimplementedFileServiceServer
}

func (*server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	fmt.Println("ListFiles was invoked")

	paths, err := ioutil.ReadDir(storagePath)
	if err != nil {
		return nil, err
	}

	var filenames []string

	for _, path := range paths {
		if !path.IsDir() {
			filenames = append(filenames, path.Name())
		}
	}

	res := &pb.ListFilesResponse{
		Filenames: filenames,
	}

	return res, nil
}

func (*server) Download(req *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	fmt.Println("Download was invoked")

	filename := req.GetFilename()
	path := storagePath + filename

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		res := &pb.DownloadResponse{Data: buf[:n]}
		sendErr := stream.Send(res)
		if sendErr != nil {
			return sendErr
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &server{})

	fmt.Println("Server is running")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
