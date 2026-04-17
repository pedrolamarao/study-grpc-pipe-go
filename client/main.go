// Copyright (c) 2025 Pedro Lamarão. All rights reserved.

package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/Microsoft/go-winio"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"pedrolamarao.dev.br/study/protocol"
)

const (
	path = `\\.\pipe\pedrolamarao.dev.br\study`
)

func closeOrPanic(closeable *grpc.ClientConn) {
	err := closeable.Close()
	if err != nil {
		panic(err)
	}
}

func main() {
	connection, err := grpc.NewClient(
		"passthrough:///"+path,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return winio.DialPipeContext(ctx, addr)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer closeOrPanic(connection)

	requestor := protocol.NewProtocolClient(connection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request := protocol.Request_builder{}
	response, err := requestor.Operation(ctx, request.Build())
	if err != nil {
		log.Fatal(err)
	}
	log.Print(response.GetValue())
}
