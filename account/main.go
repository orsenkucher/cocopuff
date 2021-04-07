//go:generate mkdir ./pb -p
//go:generate protoc ../api/proto/account.proto -I ../api/proto --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/orsenkucher/cocopuff/account/pb"
	"github.com/orsenkucher/cocopuff/account/service"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go handleSignals(cancel)

	if err := startServer(ctx); err != nil {
		log.Fatalln(err)
	}
}

func handleSignals(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	for {
		sig := <-sigCh
		switch sig {
		case os.Interrupt:
			cancel()
			return
		}
	}
}

func startServer(ctx context.Context) error {
	port := 8081
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterAccountingServer(grpcServer, &service.AccountingServer{})
	go grpcServer.Serve(lis)

	laddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		return err
	}

	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}

	defer ln.Close()

	for {
		select {
		case <-ctx.Done():
			log.Println("Server stopped")
			return nil
		default:
			if err := ln.SetDeadline(time.Now().Add(time.Second)); err != nil {
				return err
			}

			_, err := ln.Accept()
			if err != nil {
				if os.IsTimeout(err) {
					continue
				}

				return err
			}

			log.Println("New client connected")
		}
	}
}
