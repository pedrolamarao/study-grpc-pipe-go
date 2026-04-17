package main

import (
	"context"
	"log"
	"log/slog"
	"net"

	winio "github.com/Microsoft/go-winio"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"pedrolamarao.dev.br/study/protocol"
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
	path = `\\.\pipe\pedrolamarao.dev.br\study`
)

func closeOrPanic(closeable net.Listener) {
	err := closeable.Close()
	if err != nil {
		panic(err)
	}
}

func main() {
	pipe, err := winio.ListenPipe(path, &winio.PipeConfig{})
	if err != nil {
		log.Fatal(err)
	}
	defer closeOrPanic(pipe)

	service := newService()

	server := grpc.NewServer()
	protocol.RegisterProtocolServer(server, service)
	err = server.Serve(pipe)
	if err != nil {
		log.Fatal(err)
	}
}
