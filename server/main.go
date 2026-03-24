package main

import (
	"context"
	"log"
	"log/slog"

	winio "github.com/Microsoft/go-winio"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"purpura.dev.br/study/grpc/pipe/protocol"
)

type service struct {
	protocol.UnimplementedProtocolServer
}

func newService() *service {
	return &service{}
}

func (s *service) Operation(_ context.Context, _ *protocol.Request) (*protocol.Response, error) {
	response := &protocol.Response_builder{
		Value: proto.String("secret"),
	}
	slog.Info("operation")
	return response.Build(), nil
}

const (
	path = `\\.\pipe\purpura.dev.br\study\grpc\pipe`
)

func main() {
	pipe, err := winio.ListenPipe(path, &winio.PipeConfig{})
	if err != nil {
		log.Fatal(err)
	}
	defer pipe.Close()

	service := newService()

	server := grpc.NewServer()
	protocol.RegisterProtocolServer(server, service)
	err = server.Serve(pipe)
	if err != nil {
		slog.Error("", err)
	}
}
