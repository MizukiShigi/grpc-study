package main

import (
	"io"
	"context"
	"fmt"
	"grpc-lesson/pb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	// 普通は第2引数にSSL通信のオプションを設定すべき
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	// callListFiles(client)
	callDownload(client)
}

func callListFiles(client pb.FileServiceClient) {
	res, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalf("Failed to request %v", err)
	}
	fmt.Println(res.GetFilenames())
}

func callDownload(client pb.FileServiceClient) {
	req := &pb.DownloadRequest{Filename: "/name.txt"}
	stream, err := client.Download(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to request %v", err)
	}

	// サーバーストリーミングRPCなのでレスポンスをループ
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Response from donwload %v", res.GetData())
		log.Printf("Response from donwload %v", string(res.GetData()))
	}
}
